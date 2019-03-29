package rpc

import (
	"context"
	"crypto/md5"
	"encoding/json"
	"errors"
	"fmt"
	"go-sdk/conf"
	"go-sdk/file"
	"go-sdk/p2p"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/ipfs/go-ipfs-api"
	"github.com/ipfs/go-ipfs/core"
	"github.com/ironsmile/nedomi/utils"
	"go-sdk/contracts/storage"
)

var ErrNodeNotFound = errors.New("node not found")

const P2pProtocl = "/sys/http"

type Client struct {
	IpfsClients              map[string]*shell.Shell
	NodesRefreshTime         time.Time
	NodesRefreshInterval     time.Duration
	DurationToDiscoveryNodes time.Duration
	*storage.Client
	*core.IpfsNode
}

func NewClient(cfg conf.Config, node *core.IpfsNode) (cli *Client, err error) {
	cli = &Client{IpfsNode: node}
	cli.IpfsClients = make(map[string]*shell.Shell)
	cli.NodesRefreshInterval = time.Second
	cli.DurationToDiscoveryNodes = time.Second * 3

	c, err := storage.NewClient(cfg.ContractConfig)
	if err != nil {
		return
	}
	cli.Client = c

	return
}

func (c *Client) Upload(fpath string) (cid string, err error) {

	fname := path.Base(fpath)
	fh, err := os.Open(fpath)
	if err != nil {
		return
	}
	defer fh.Close()

	fi, err := fh.Stat()
	if err != nil {
		return
	}

	hs := md5.New()
	freader := io.TeeReader(fh, hs)

	//cid, err = file.GetCID(nil)
	//if err != nil {
	//	return
	//}

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
	dataShards, parShards := file.BlockCount(totalMetaLength, fi.Size(), 0.3)
	totalDataLength := int64(dataShards*totalMetaLength) + fi.Size()

	fmt.Println("dshards, parshards", dataShards, parShards)
	mgr, err := file.NewBlockMgr(dataShards, parShards)
	if err != nil {
		return
	}

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
	shardsRdr, err := mgr.ECShards(freader, getMeta, totalDataLength)
	if err != nil {
		return
	}

	nodes, err := c.GetIpfsClientsWithId()
	if err != nil {
		return
	}

	//cid
	cid = fmt.Sprintf("%x", hs.Sum(nil))
	_, err = c.NewUploadJob(cid, fi.Size(), dataShards+parShards)
	if err != nil {
		return
	}

	//upload block
	for i := range shardsRdr {

		//body, err := ioutil.ReadAll(shardsRdr[i])
		//fmt.Println("index:", i, "body:", string(body), "err:", err)

		nodeIdx := i % len(nodes)
		node := nodes[nodeIdx]
		blockHash, err := node.c.Add(shardsRdr[i])
		if err != nil {
			return "", err
		}

		fmt.Println("blockHash: ", blockHash)
		//err = c.CommitBlock(job, i, blockHash, node.id)
		if err != nil {
			return "", err
		}
	}

	return
}

func (c *Client) Download(fileHash string) (rc io.ReadCloser, metaAll file.MetaAll, err error) {
	blocksInfo, err := c.GetBlocksInfo(fileHash)
	log.Println("GetBlocksInfo fileHash:", fileHash, "\tblocks info:", blocksInfo, "\terr:", err)
	if err != nil {
		return
	}

	firstBlock := blocksInfo[0]
	node, err := c.getIpfsClientOrRandom(string(firstBlock.PeerId))
	if err != nil {
		return
	}

	metaAll, err = GetMeta(node, string(firstBlock.BlockHash))
	if err != nil {
		return
	}

	metaAllLength := file.MetaBytes + metaAll.MetaExLength
	dataShards := int(metaAll.DataShards)
	rcs := make([]io.ReadCloser, dataShards)
	for i := 0; i < dataShards; i++ {
		blockInfo := blocksInfo[i]
		node, err = c.getIpfsClientOrRandom(string(blockInfo.PeerId))
		if err != nil {
			return
		}
		rc1, err1 := ReadAt(node, string(blockInfo.BlockHash), int64(metaAllLength), 0)
		if err1 != nil {
			err = err1
			return
		}
		rcs[i] = rc1
	}
	rc = utils.MultiReadCloser(rcs...)
	err = c.DownloadSuccess(fileHash)
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

type IpfsClientWithId struct {
	id string
	c  *shell.Shell
}

func (c *Client) GetIpfsClientsWithId() (cs []IpfsClientWithId, err error) {

	getClientWithId := func() []IpfsClientWithId {
		cs := []IpfsClientWithId{}
		for id, c := range c.IpfsClients {
			cs = append(cs, IpfsClientWithId{id, c})
		}
		return cs
	}
	if !c.needRefresh() {
		cs = getClientWithId()
		return
	}

	err = c.refreshIpfsClients()
	if err == ErrNodeNotFound {
		time.Sleep(c.DurationToDiscoveryNodes)
		err = c.refreshIpfsClients()
	}
	if err != nil {
		return
	}

	cs = getClientWithId()
	return
}

func (c *Client) GetIpfsClients() (ss []*shell.Shell, err error) {

	if !c.needRefresh() {
		ss = c.ipfsClients()
		return
	}

	err = c.refreshIpfsClients()
	if err == ErrNodeNotFound {
		time.Sleep(c.DurationToDiscoveryNodes)
		err = c.refreshIpfsClients()
	}
	if err != nil {
		return
	}

	ss = c.ipfsClients()
	return
}

func (c *Client) ipfsClients() (ss []*shell.Shell) {
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
		//"Qmctjf76oDSKjLXS5QzNMbkDT4hVkSENRjpbaYn3dqeZ78": true,
		"QmWb9ra6trs9HpXp4dRH1WPucV7Xin3cG3AD4Dswp4sEmk": true,
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
