package rpc

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strconv"
	"sync"
	"time"

	"github.com/ipfs/go-ipfs-api"
	"github.com/ipfs/go-ipfs/core"
	"github.com/ipweb-group/go-sdk/conf"
	"github.com/ipweb-group/go-sdk/contracts/storage"
	"github.com/ipweb-group/go-sdk/file"
	"github.com/ipweb-group/go-sdk/p2p"
	"github.com/ipweb-group/go-sdk/utils/netools"
	"github.com/ironsmile/nedomi/utils"
)

var ErrNodeNotFound = errors.New("node not found")

const P2pProtocl = "/sys/http"

type Client struct {
	IpfsClients            map[string]*shell.Shell
	IpfsUnabailableClients map[string]*shell.Shell
	NodesRefreshTime       time.Time
	NodesRefreshDuration   time.Duration
	BlockUpWorkerCount     int
	*storage.Client
	*core.IpfsNode
}

func NewClient(cfg conf.Config) (cli *Client, err error) {
	ctx := context.Background()
	n, err := p2p.NewNode(ctx)
	if err != nil {
		return
	}

	cli = &Client{IpfsNode: n}
	cli.IpfsClients = make(map[string]*shell.Shell)
	cli.IpfsUnabailableClients = make(map[string]*shell.Shell)
	if cfg.NodesRefreshIntervalInSecond == 0 {
		cfg.NodesRefreshIntervalInSecond = 300
	}
	cli.NodesRefreshDuration = time.Minute * time.Duration(cfg.NodesRefreshIntervalInSecond)

	if cfg.BlockUpWorkerCount == 0 {
		cfg.BlockUpWorkerCount = 2
	}
	cli.BlockUpWorkerCount = cfg.BlockUpWorkerCount

	c, err := storage.NewClient(cfg.ContractConf)
	if err != nil {
		return
	}
	cli.Client = c
	go cli.refreshNodePeers()

	return
}

func (c *Client) Upload(rdr io.Reader, fname string, fsize int64) (cid string, err error) {
	br := &bytes.Buffer{}
	cidRdr := io.TeeReader(rdr, br)
	cid, err = file.GetCID(cidRdr)
	if err != nil {
		return
	}

	dataShards, parShards, shardSize := file.BlockCount(fsize)
	mgr, err := file.NewBlockMgr(dataShards, parShards)
	if err != nil {
		return
	}
	meta := file.NewMeta(fname, cid, fsize, uint32(dataShards), uint32(parShards))

	shardsRdr, err := mgr.ECShards(br, fsize)
	if err != nil {
		return
	}

	nodes, err := c.GetNodeClients("")
	if err != nil {
		return
	}

	shards := dataShards + parShards
	shardSize += int64(len(meta.Encode(0)))
	_, err = c.NewUploadJob(cid, fsize, shards, shardSize)
	if err != nil {
		return
	}

	shardIdCh := make(chan int, shards)
	for i := 0; i < shards; i++ {
		shardIdCh <- i
	}
	close(shardIdCh)
	errCh := make(chan error, shards)
	wg := sync.WaitGroup{}
	wg.Add(c.BlockUpWorkerCount)
	worker := func() {
		defer wg.Done()
		for id := range shardIdCh {
			retry := 0
		lazyTry:
			nodeIdx := (id + retry) % len(nodes)
			node := nodes[nodeIdx]

			blkBuf := bytes.Buffer{}
			blkRdr := io.TeeReader(shardsRdr[id], &blkBuf)
			mr := bytes.NewBuffer(meta.Encode(id))
			var r io.Reader
			if retry == 0 {
				r = io.MultiReader(mr, blkRdr)
			} else {
				bufRdr := bytes.NewReader(blkBuf.Bytes())
				r = io.MultiReader(mr, bufRdr)
			}

			blkHash, err := node.Client.Add(r)
			if err != nil {
				if retry < len(nodes) {
					retry++
					log.Println("err:", err, "retrying... retry times:", retry)
					goto lazyTry
				}
				errCh <- err
			}
			log.Println("Block Hash:", blkHash, err)
		}
	}
	for i := 0; i < c.BlockUpWorkerCount; i++ {
		go worker()
	}
	wg.Wait()
	close(errCh)
	for err1 := range errCh {
		cid = ""
		err = err1
	}

	return
}

func (c *Client) UploadWithPath(fpath string) (cid string, err error) {
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

	cid, err = c.Upload(fh, fname, fi.Size())
	return
}

func (c *Client) Download(fileHash string) (rc io.ReadCloser, metaAll file.Meta, err error) {
	blocksInfo, err := c.GetBlocksInfo(fileHash)
	if err != nil {
		return
	}

	meta, err := c.GetMeta(blocksInfo)
	if err != nil {
		return
	}

	dataShards := int(meta.MetaHeader.DataShards)
	rcs := make([]io.ReadCloser, dataShards)
	for i := 0; i < dataShards; i++ {
		blockInfo := blocksInfo[i]
		ns, err1 := c.GetNodeClients(blockInfo.PeerId)
		if err1 != nil {
			err = err1
			return
		}

		var rc1 io.ReadCloser
		for _, n := range ns {
			rc1, err1 = ReadAt(n.Client, blockInfo.BlockHash, int64(meta.Len()), 0)
			if err1 == nil {
				break
			}
		}

		rcs[i] = rc1
	}
	rc = utils.MultiReadCloser(rcs...)
	return
}

func (c *Client) GetMeta(bis []storage.BlockInfo) (meta *file.Meta, err error) {
	for _, bi := range bis {
		clients, err := c.GetNodeClients(bi.PeerId)
		if err != nil {
			return nil, err
		}

		for _, n := range clients {
			meta, err := getMeta(n.Client, bi.BlockHash)
			if err == nil {
				return meta, err
			}
		}
	}
	return
}

func getMeta(node *shell.Shell, blkHash string) (m *file.Meta, err error) {
	rc1, err := ReadAt(node, blkHash, 0, file.MetaHeaderLength)
	if err != nil {
		return
	}
	defer rc1.Close()
	metaHeaderB, err := ioutil.ReadAll(rc1)
	if err != nil {
		return
	}
	metaHeader, err := file.DecodeMetaHeader(metaHeaderB)
	if err != nil {
		return
	}

	metaLen := file.MetaHeaderLength + metaHeader.MetaBodyLength
	rc2, err := ReadAt(node, blkHash, 0, int64(metaLen))
	if err != nil {
		return
	}
	defer rc2.Close()
	metaData, err := ioutil.ReadAll(rc2)
	if err != nil {
		return
	}

	return file.DecodeMeta(metaData)
}

func ReadAt(node *shell.Shell, fp string, offset, length int64) (rc io.ReadCloser, err error) {
	req := node.Request("cat", fp).Option("offset", offset)
	if length != 0 {
		req.Option("length", length)
	}
	resp, err := req.Send(context.Background())
	if err != nil {
		return
	}
	rc = resp.Output
	return
}

func (c *Client) GetClientByPeerId(pId string) (node *shell.Shell, exist bool) {
	if c.needRefresh() {
		c.refreshNodePeers()
	}

	node, exist = c.IpfsClients[pId]
	return
}

type NodeClient struct {
	Id     string
	Client *shell.Shell
}

func (c *Client) GetNodeClients(nodeIdMoveToFirstElement string) (ns []NodeClient, err error) {
	getNodes := func() []NodeClient {
		var ns1, ns2 []NodeClient
		for id, n := range c.IpfsClients {
			if id == nodeIdMoveToFirstElement {
				ns1 = append(ns1, NodeClient{Id: id, Client: n})
				continue
			}
			ns2 = append(ns2, NodeClient{id, n})
		}
		return append(ns1, ns2...)
	}
	if !c.needRefresh() {
		ns = getNodes()
		return
	}

	err = c.refreshNodePeers()
	if err == ErrNodeNotFound {
		err = c.refreshNodePeers()
	}
	if err != nil {
		return
	}

	ns = getNodes()
	return
}

func (c *Client) needRefresh() bool {
	timeOut := c.NodesRefreshTime.Add(c.NodesRefreshDuration).Before(time.Now())

	if timeOut || len(c.IpfsClients) == 0 {
		return true
	}
	return false
}

func (c *Client) refreshNodePeers() error {
	c.NodesRefreshTime = time.Now()
	clients := make(map[string]*shell.Shell)
	ps := c.IpfsNode.Peerstore.Peers()
	for _, p := range ps {
		id := p.Pretty()
		cli, ok := c.IpfsClients[id]
		if ok {
			clients[id] = cli
			continue
		}

		_, ok = c.IpfsUnabailableClients[id]
		if ok {
			continue
		}

		cli, err := c.NewIpfsClient(id)
		if err != nil {
			c.IpfsUnabailableClients[id] = cli
			c.P2PClose(0, id)
			log.Println("bad peer: ", p.Pretty(), err)
			continue
		}

		log.Println("p2p peer: ", p.Pretty())
		clients[id] = cli
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

func (c *Client) NewIpfsClient(peerId string) (cli *shell.Shell, err error) {
	port, err := netools.GetFreePort()
	if err != nil {
		return
	}

	err = c.P2PForward(port, peerId)
	if err != nil {
		return
	}
	url := fmt.Sprintf("127.0.0.1:%d", port)

	cli = shell.NewShell(url)
	_, err = cli.ID()
	if err != nil {
		c.P2PClose(0, peerId)
	}

	return
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
