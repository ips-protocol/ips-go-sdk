package rpc

import (
	"encoding/json"
	"errors"
	"fmt"
	"go-sdk/contract"
	"go-sdk/file"
	"io"
	"log"
	"net"
	"os"
	"path"
	"strconv"
	"time"

	"github.com/ipfs/go-ipfs-api"
	"github.com/ipfs/go-ipfs/core"
	"go-sdk/p2p"
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

func (c *Client) Download(hash string) (rc io.ReadCloser, meta file.MetaAll, err error) {
	nodes, err := c.GetIpfsClients()
	if err != nil {
		return
	}

	//cid, err = ss[0].Get()
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
