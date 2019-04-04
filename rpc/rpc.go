package rpc

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"math/rand"
	"net"
	"os"
	"path"
	"strconv"
	"sync"
	"time"

	"github.com/ipfs/go-ipfs-api"
	"github.com/ipweb-group/go-ipfs/core"
	"github.com/ipweb-group/go-sdk/conf"
	"github.com/ipweb-group/go-sdk/contracts/storage"
	"github.com/ipweb-group/go-sdk/file"
	"github.com/ipweb-group/go-sdk/p2p"
	"github.com/ironsmile/nedomi/utils"
)

var ErrNodeNotFound = errors.New("node not found")

const P2pProtocl = "/sys/http"

type Client struct {
	IpfsClients              map[string]*shell.Shell
	IpfsUnabailableClients   map[string]*shell.Shell
	NodesRefreshTime         time.Time
	NodesRefreshInterval     time.Duration
	DurationToDiscoveryNodes time.Duration
	WorkerCounts             int
	*storage.Client
	*core.IpfsNode
}

func NewClient(cfg conf.Config, node *core.IpfsNode) (cli *Client, err error) {
	cli = &Client{IpfsNode: node}
	cli.IpfsClients = make(map[string]*shell.Shell)
	cli.IpfsUnabailableClients = make(map[string]*shell.Shell)
	cli.NodesRefreshInterval = time.Second
	cli.DurationToDiscoveryNodes = time.Second * 3
	cli.WorkerCounts = cfg.BlockUpWorkerCount
	if cli.WorkerCounts == 0 {
		cli.WorkerCounts = 2
	}

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

	//cid
	cid, err = file.GetCID(fh)
	if err != nil {
		return
	}
	_, err = fh.Seek(io.SeekStart, io.SeekStart)
	if err != nil {
		return
	}

	dataShards, parShards := file.BlockCount(fi.Size())
	mgr, err := file.NewBlockMgr(dataShards, parShards)
	if err != nil {
		return
	}

	meta := file.NewMeta(fname, cid, fi.Size(), uint32(dataShards), uint32(parShards))

	shardsRdr, err := mgr.ECShards(fh, fi.Size())
	if err != nil {
		return
	}

	nodes, err := c.GetIpfsClientsWithId()
	if err != nil {
		return
	}

	shards := dataShards + parShards
	//_, err = c.NewUploadJob(cid, fi.Size(), shards)
	//if err != nil {
	//	return
	//}

	//upload block
	shardIdCh := make(chan int, shards)
	for i := 0; i < shards; i++ {
		shardIdCh <- i
	}
	close(shardIdCh)
	errCh := make(chan error, shards)
	wg := sync.WaitGroup{}
	wg.Add(c.WorkerCounts)
	worker := func() {
		defer wg.Done()
		for id := range shardIdCh {
			//for retry
			nodeIdx := id % len(nodes)
			node := nodes[nodeIdx]

			mr := bytes.NewBuffer(meta.Encode(id))
			r := io.MultiReader(mr, shardsRdr[id])

			blkHash, err := node.c.Add(r)
			if err != nil {
				errCh <- err
			}
			fmt.Println("Block Hash:", blkHash, err)
		}
	}
	for i := 0; i < c.WorkerCounts; i++ {
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

func (c *Client) Download(fileHash string) (rc io.ReadCloser, metaAll file.Meta, err error) {
	//blocksInfo, err := c.GetBlocksInfo(fileHash)
	//log.Println("GetBlocksInfo fileHash:", fileHash, "\tblocks info:", blocksInfo, "\terr:", err)
	//if err != nil {
	//	return
	//}

	blocksInfo := []storage.BlockInfo{
		{[]byte("Qmd6Nj3GFKxo3LT2gAdoVzSjoBWSfYKGdFRyLjtjKDkvJ1"), []byte("QmWb9ra6trs9HpXp4dRH1WPucV7Xin3cG3AD4Dswp4sEmk")},
		{[]byte("QmVdtWzv1Nem7BiJTMDCcvza38hMKCYARBWb4LZjsvtNYU"), []byte("QmWb9ra6trs9HpXp4dRH1WPucV7Xin3cG3AD4Dswp4sEmk")},
		{[]byte("QmYYwHMtceagVxG79fE74LETmMvyi9A6QZ9TDHeU3K6Dqv"), []byte("QmWb9ra6trs9HpXp4dRH1WPucV7Xin3cG3AD4Dswp4sEmk")},
		{[]byte("QmZM2uYRVFvWryiw6pExBberjQNSASBZsEAWeJDAy7AQDt"), []byte("QmWb9ra6trs9HpXp4dRH1WPucV7Xin3cG3AD4Dswp4sEmk")},
		{[]byte("QmTXpM13pqYToQMRgwXUhwXrmkGfBPiZ68MG9ijvMaf69G"), []byte("QmWb9ra6trs9HpXp4dRH1WPucV7Xin3cG3AD4Dswp4sEmk")},
		{[]byte("QmXcZpPCSvs6ZvRtEWmNNAfYFxEdiCYqU3TwNZ2P1sFbTW"), []byte("QmWb9ra6trs9HpXp4dRH1WPucV7Xin3cG3AD4Dswp4sEmk")},
	}

	firstBlock := blocksInfo[0]
	node, err := c.getIpfsClientOrRandom(string(firstBlock.PeerId))
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
		node, err = c.getIpfsClientOrRandom(string(blockInfo.PeerId))
		if err != nil {
			return
		}
		rc1, err1 := ReadAt(node, string(blockInfo.BlockHash), int64(meta.Len()), 0)
		if err1 != nil {
			err = err1
			return
		}
		rcs[i] = rc1
	}
	rc = utils.MultiReadCloser(rcs...)
	//err = c.DownloadSuccess(fileHash)
	return
}

func (c *Client) GetMeta(bis []storage.BlockInfo) (meta *file.Meta, err error) {
	for _, bi := range bis {
		cli, err1 := c.getIpfsClientOrRandom(string(bi.PeerId))
		if err1 != nil {
			continue
		}
		meta, err = getMeta(cli, string(bi.BlockHash))
		if err == nil {
			return
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
			fmt.Println("bad peer: ", p.Pretty(), err)
			continue
		}

		fmt.Println("p2p peer: ", p.Pretty())
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
	port, err := GetFreePort()
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
