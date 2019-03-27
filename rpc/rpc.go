package rpc

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"os"
	"path"
	"strconv"
	"time"

	contract "../contract"
	file "../file"

	p2p "../p2p"

	shell "github.com/ipfs/go-ipfs-api"
	"github.com/ipfs/go-ipfs/core"
	"github.com/ironsmile/nedomi/utils"
)

var ErrNodeNotFound = errors.New("node not found")

const P2pProtocl = "/sys/http"

type Client struct {
	IpfsClients          map[string]*shell.Shell
	NodesRefreshTime     time.Time
	NodesRefreshInterval time.Duration
	contract.Client
	contract.StorageAccount
	*core.IpfsNode
}

func NewClient(node *core.IpfsNode) (cli *Client, err error) {
	cli = &Client{IpfsNode: node}
	cli.IpfsClients = make(map[string]*shell.Shell)
	cli.NodesRefreshInterval = time.Second
	return
}

func (c *Client) Upload(fpath string) (cid string, err error) {

	fname := path.Base(fpath)
	fh, err := os.Open(fpath)
	if err != nil {
		return
	}
	fi, err := fh.Stat()
	if err != nil {
		return
	}

	cid, err = file.GetCID(nil)
	if err != nil {
		return
	}

	metaEx := file.MetaEx{
		FName: fname,
		FSize: fi.Size(),
		FHash: cid,
	}

	metaExB, err := json.Marshal(metaEx)
	if err != nil {
		return
	}
	totalMetaLength := len(metaExB) + file.MetaBytes
	dataShards, parShards := file.BlockCount(totalMetaLength, fi.Size(), 1/3)
	totalDataLength := int64(dataShards*totalMetaLength) + fi.Size()
	getMeta := func(i int) []byte {
		meta := file.Meta{
			DataShards:   uint32(dataShards),
			ParShards:    uint32(parShards),
			BlockIdx:     uint32(i),
			MetaExLength: uint32(len(metaExB)),
		}
		metaBytes := file.EncodeMeta(&meta)
		metaBytes = append(metaBytes, metaExB...)
		return metaBytes
	}

	cfg := file.Config{DataShards: dataShards, ParShards: parShards}
	mgr, err := file.NewBlockMgr(cfg)
	if err != nil {
		return
	}

	shardsRdr, err := mgr.ECShards(fh, getMeta, totalDataLength)
	if err != nil {
		return
	}

	nodes, err := c.GetIpfsClients()
	if err != nil {
		return
	}

	for i := range shardsRdr {
		nodeIdx := i % len(nodes)
		_, err = nodes[nodeIdx].Add(shardsRdr[i])
		if err != nil {
			return
		}
	}

	_, err = c.NewUploadJob(cid, uint64(totalDataLength), uint64(dataShards+parShards), 0)
	return
}

func (c *Client) Download(hash string) (rc io.ReadCloser, metaAll file.MetaAll, err error) {
	blocksHash, peersInfo, err := c.GetAllBlocksInfo()
	if err != nil {
		return
	}

	node, err := c.getIpfsClientOrRandom(peersInfo[0])
	if err != nil {
		return
	}

	metaAll, err = GetMeta(node, blocksHash[0])
	if err != nil {
		return
	}

	metaAllLength := file.MetaBytes + metaAll.MetaExLength
	dataShards := int(metaAll.DataShards)
	rcs := make([]io.ReadCloser, dataShards)
	for i := 0; i < dataShards; i++ {
		node, err = c.getIpfsClientOrRandom(peersInfo[i])
		if err != nil {
			return
		}
		var rc3 io.ReadCloser
		rc3, err = ReadAt(node, blocksHash[i], int64(metaAllLength), 0)
		if err != nil {
			return
		}
		rcs[i] = rc3
	}
	rc = utils.MultiReadCloser(rcs...)
	return
}

func GetMeta(node *shell.Shell, blockHash string) (metaAll file.MetaAll, err error) {

	rc1, err := ReadAt(node, blockHash, 0, file.MetaBytes)
	if err != nil {
		return
	}
	defer rc1.Close()
	metaB, err := ioutil.ReadAll(rc1)
	if err != nil {
		return
	}
	meta, err := file.DecodeMeta(metaB)
	if err != nil {
		return
	}

	rc2, err := ReadAt(node, blockHash, file.MetaBytes, int64(meta.MetaExLength))
	if err != nil {
		return
	}
	defer rc2.Close()
	metaExB, err := ioutil.ReadAll(rc2)
	if err != nil {
		return
	}
	metaEx := file.MetaEx{}
	err = json.Unmarshal(metaExB, &metaEx)
	if err != nil {
		return
	}
	metaAll = file.MetaAll{Meta: *meta, MetaEx: metaEx}
	return
}

func ReadAt(node *shell.Shell, fp string, offset, count int64) (rc io.ReadCloser, err error) {
	req := node.Request("files/read", fp)
	if offset != 0 {
		req.Option("offset", offset)
	}
	if count != 0 {
		req.Option("count", count)
	}
	resp, err := req.Send(context.Background())
	if err != nil {
		return
	}
	rc = resp.Output
	return
}

func (c *Client) getIpfsClientOrRandom(pId string) (node *shell.Shell, err error) {
	node, exist := c.GetIpfsClient(pId)
	if exist {
		return
	}

	node, err = c.GetRandomIpfsClient()
	return
}

func (c *Client) GetRandomIpfsClient() (node *shell.Shell, err error) {
	nodes, err := c.GetIpfsClients()
	if err != nil {
		return
	}

	rand.Seed(time.Now().Unix())
	node = nodes[rand.Intn(len(nodes))]
	return
}

func (c *Client) GetIpfsClient(pId string) (node *shell.Shell, exist bool) {
	if c.needRefresh() {
		c.refreshIpfsClients()
	}

	node, exist = c.IpfsClients[pId]
	return
}

func (c *Client) GetIpfsClients() (ss []*shell.Shell, err error) {

	if !c.needRefresh() {
		ss = c.getIpfsClients()
		return
	}

	err = c.refreshIpfsClients()
	if err == ErrNodeNotFound {
		time.Sleep(time.Second * 3)
		err = c.refreshIpfsClients()
	}
	if err != nil {
		return
	}

	ss = c.getIpfsClients()
	return
}

func (c *Client) getIpfsClients() (ss []*shell.Shell) {
	for _, ic := range c.IpfsClients {
		ss = append(ss, ic)
	}

	return
}

func (c *Client) needRefresh() bool {
	timeOut := c.NodesRefreshTime.Add(c.NodesRefreshInterval).Before(time.Now())
	if timeOut || len(c.IpfsClients) == 0 {
		return true
	}
	return false
}

func (c *Client) refreshIpfsClients() error {

	clients := make(map[string]*shell.Shell)
	ps := c.IpfsNode.Peerstore.Peers()
	for _, p := range ps {
		fmt.Println("peer: ---->", p.Pretty())

		id := p.Pretty()
		if !isP2PNode(id) {
			continue
		}

		cli, ok := c.IpfsClients[id]
		if ok {
			clients[id] = cli
			continue
		}

		port, err := GetFreePort()
		if err != nil {
			return err
		}

		err = c.P2PForward(port, id)
		if err != nil {
			log.Println(err)
			continue
		}

		clients[id] = shell.NewShell(localAddr(port))
	}

	for id := range c.IpfsClients {
		if _, ok := clients[id]; ok {
			continue
		}

		c.P2PClose(0, id)
	}

	if len(clients) == 0 {
		return ErrNodeNotFound
	}

	c.IpfsClients = clients
	return nil
}

func localAddr(port int) string {
	return "127.0.0.1:" + strconv.Itoa(port)
}

func (c *Client) P2PForward(port int, peerId string) error {
	listenOpt := "/ip4/127.0.0.1/tcp/" + strconv.Itoa(port)
	targetOpt := "/ipfs/" + peerId
	return p2p.Forward(c.IpfsNode.P2P, P2pProtocl, listenOpt, targetOpt)
}

func (c *Client) P2PClose(port int, peerId string) error {
	listenOpt := "/ip4/127.0.0.1/tcp/" + strconv.Itoa(port)
	targetOpt := "/ipfs/" + peerId
	return p2p.Close(c.IpfsNode.P2P, false, "", listenOpt, targetOpt)
}

func (c *Client) P2PCloseAll() error {
	return p2p.Close(c.IpfsNode.P2P, true, "", "", "")
}

func isP2PNode(id string) bool {
	nodes := map[string]bool{
		"Qmain1GGsLNtPmDPJsmWGYv7QxyFnbTjvFceH16yC2PCRd": true,
		"Qma2z1RDNNTH2NVpQbzaoZZs8CxHU9N881ZGKB4oZuabXw": true,
	}
	_, ok := nodes[id]
	return ok
}

func GetFreePort() (port int, err error) {
	addr, err := net.ResolveTCPAddr("tcp", "localhost:0")
	if err != nil {
		return
	}

	l, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return
	}
	defer l.Close()

	port = l.Addr().(*net.TCPAddr).Port
	return
}
