package rpc

import (
	"bytes"
	"context"
	"crypto/sha256"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
	"path"
	"strconv"
	"sync"
	"time"

	"github.com/ipfs/go-ipfs-api"
	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/metafile"
	"github.com/ipweb-group/go-sdk/conf"
	"github.com/ipweb-group/go-sdk/contracts/storage"
	"github.com/ipweb-group/go-sdk/file"
	"github.com/ipweb-group/go-sdk/p2p"
	"github.com/ipweb-group/go-sdk/utils/netools"
	"github.com/ironsmile/nedomi/utils"
)

var ErrNodeNotFound = errors.New("node not found")
var ErrContractNotFound = errors.New("no contract code at given address")

const P2pProtocl = "/sys/http"

type Client struct {
	IpfsClients              map[string]*shell.Shell
	IpfsUnabailableClients   map[string]*shell.Shell
	NodesRefreshTime         time.Time
	NodesRefreshDuration     time.Duration
	NodeRequestTimeout       time.Duration
	BlockUpWorkerCount       int
	BlockDownloadWorkerCount int
	WalletPubKey             string
	*storage.Client
	*core.IpfsNode
}

func NewClient(cfg conf.Config) (cli *Client, err error) {
	ctx := context.Background()
	n, err := p2p.NewNode(ctx, cfg)
	if err != nil {
		return
	}

	cli = &Client{IpfsNode: n}
	cli.IpfsClients = make(map[string]*shell.Shell)
	cli.IpfsUnabailableClients = make(map[string]*shell.Shell)
	if cfg.NodesRefreshIntervalInSecond == 0 {
		cfg.NodesRefreshIntervalInSecond = 600
	}
	cli.NodesRefreshDuration = time.Second * time.Duration(cfg.NodesRefreshIntervalInSecond)

	if cfg.NodeRequestTimeoutInSecond == 0 {
		cfg.NodeRequestTimeoutInSecond = 60
	}
	cli.NodeRequestTimeout = time.Second * time.Duration(cfg.NodeRequestTimeoutInSecond)

	if cfg.BlockUpWorkerCount == 0 {
		cfg.BlockUpWorkerCount = 5
	}
	cli.BlockUpWorkerCount = cfg.BlockUpWorkerCount

	if cfg.BlockDownloadWorkerCount == 0 {
		cfg.BlockDownloadWorkerCount = 5
	}
	cli.BlockDownloadWorkerCount = cfg.BlockDownloadWorkerCount

	pubKey, err := conf.GetWalletPubKey()
	if err != nil {
		return
	}
	cli.WalletPubKey = pubKey

	c, err := storage.NewClient(cfg.ContractConf)
	if err != nil {
		return
	}
	cli.Client = c
	go cli.refreshNodePeers()

	return
}

func (c *Client) Upload(rdr io.Reader, fname string, fsize int64) (cid string, err error) {
	dataShards, parShards, shardSize := file.BlockCount(fsize)
	mgr, err := file.NewBlockMgr(dataShards, parShards)
	if err != nil {
		return
	}

	h := sha256.New()
	_, err = h.Write([]byte(c.WalletPubKey))
	if err != nil {
		return
	}
	r := io.TeeReader(rdr, h)

	//dataFhs := make([]*os.File, dataShards)
	dataFhs, err := mgr.Split(r, fsize)
	if err != nil {
		return
	}
	cid, err = file.GetCidV0(h)
	if err != nil {
		return
	}

	parFhs, err := file.CreateTmpFiles(parShards)
	if err != nil {
		return
	}

	fmt.Println("===> shards:", dataShards, "parshards:", parShards, "shardsize:", shardSize)
	dataFhRdrs := make([]io.Reader, dataShards)
	for i := range dataFhRdrs {
		dataFhRdrs[i] = dataFhs[i]
	}
	parFhWtrs := make([]io.Writer, parShards)
	for i := range parFhWtrs {
		parFhWtrs[i] = parFhs[i]
	}
	err = mgr.Encode(dataFhRdrs, parFhWtrs)
	if err != nil {
		return
	}
	fhs := append(dataFhs, parFhs...)
	defer func() {
		err := file.DeleteTempFiles(fhs)
		if err != nil {
			log.Println("delete files failed:", err)
		}
	}()

	meta := metafile.NewMeta(fname, cid, fsize, uint32(dataShards), uint32(parShards))
	meta.WalletPubKey = c.WalletPubKey
	shardSize += int64(len(meta.Encode(0)))
	shards := dataShards + parShards
	_, err = c.NewUploadJob(cid, fsize, shards, shardSize)
	if err != nil {
		return
	}

	for i := range fhs {
		fhs[i].Seek(0, 0)
	}
	err = c.upload(fhs, meta)

	return
}

func (c *Client) upload(fhs []*os.File, meta metafile.Meta) error {
	shards := len(fhs)
	nodes, err := c.GetNodeClients("")
	if err != nil {
		return err
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

			mr := bytes.NewBuffer(meta.Encode(id))
			var r io.Reader
			if retry == 0 {
				r = io.MultiReader(mr, fhs[id])
			} else {
				fhs[id].Seek(0, 0)
				r = io.MultiReader(mr, fhs[id])
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
			log.Printf("block hash: %s, node id: %s, err: %#v", blkHash, node.Id, err)
		}
	}
	for i := 0; i < c.BlockUpWorkerCount; i++ {
		go worker()
	}
	wg.Wait()
	close(errCh)
	for err1 := range errCh {
		return err1
	}
	return nil
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

func (c *Client) Remove(fHash string) error {
	blocksInfo, err := c.GetBlocksInfo(fHash)
	if err != nil {
		return err
	}

	for _, bi := range blocksInfo {
		node, err := c.GetNodeClient(bi.PeerId)
		if err != nil {
			return err
		}

		err = node.Client.Unpin(bi.BlockHash)
		if err != nil {
			return err
		}
	}

	err = c.DeleteFile(fHash)
	if err != nil && err.Error() == ErrContractNotFound.Error() {
		err = ErrContractNotFound
	}
	return err
}

func (c *Client) Download(fileHash string, w io.Writer) (metaAll metafile.Meta, err error) {
	blocksInfo, err := c.GetBlocksInfo(fileHash)
	if err != nil {
		return
	}

	meta, err := c.GetMeta(blocksInfo)
	if err != nil {
		return
	}

	dataShards := int(meta.MetaHeader.DataShards)
	parShards := int(meta.MetaHeader.ParShards)
	shards := dataShards + parShards
	dataFhs, hasBroken, err := c.download(blocksInfo[:dataShards], len(meta.Encode(0)))
	if err != nil {
		return
	}
	if hasBroken {
		log.Println("have broken shards, start download parity shards.")
		parFhs, _, e := c.download(blocksInfo[dataShards:], len(meta.Encode(0)))
		if e != nil {
			err = e
			return
		}
		log.Println("parity shards download success.")
		mgr, e := file.NewBlockMgr(dataShards, parShards)
		if err != nil {
			err = e
			return
		}

		fhs := append(dataFhs, parFhs...)
		rdrs := make([]io.Reader, shards)
		wtrs := make([]io.Writer, shards)
		for i := range fhs {
			if fhs[i] == nil {
				rdrs[i] = nil
				wtrs[i], err = file.CreteTmpFile()
				if err != nil {
					return
				}
			} else {
				rdrs[i] = fhs[i]
			}
		}

		log.Println("reconstruct start.")
		err = mgr.Reconstruct(rdrs, wtrs)
		if err != nil {
			return
		}
		log.Println("reconstruct end.")

		file.SeekToStart(fhs)
		for i := range fhs {
			if fhs[i] == nil {
				fh := wtrs[i].(*os.File)
				fh.Seek(0, 0)
				fhs[i] = fh
			}
		}
		dataFhs = fhs[:dataShards]
	}

	log.Println("download over.")
	drs := make([]io.Reader, dataShards)
	for i := range dataFhs {
		drs[i] = dataFhs[i]
	}
	fr := io.LimitReader(io.MultiReader(drs...), meta.FSize)
	_, err = io.Copy(w, fr)
	return
}

func (c *Client) download(blocksInfo []storage.BlockInfo, metaLen int) (fhs []*os.File, hasBroken bool, err error) {
	blockNum := len(blocksInfo)
	fhs, err = file.CreateTmpFiles(blockNum)
	if err != nil {
		return
	}

	sem := make(chan bool, c.BlockDownloadWorkerCount)
	for i := 0; i < blockNum; i++ {
		sem <- true
		go func(idx int) (err error) {
			log.Println("block idx:", idx)
			defer func() {
				<-sem
				if err != nil {
					log.Printf("block download failed idx: %d, error: %#v\n", idx, err)
					fhs[idx].Close()
					os.Remove(fhs[idx].Name())
					fhs[idx] = nil
					hasBroken = true
				} else {
					fhs[idx].Seek(0, 0)
				}
			}()

			blockInfo := blocksInfo[idx]
			node, err := c.GetNodeClient(blockInfo.PeerId)
			if err != nil {
				return
			}

			dialRetryTimes := 0
			downloadRetryTimes := 0
		lazyTry:
			dialStart := time.Now()
			rc1, err := ReadAt(node.Client, blockInfo.BlockHash, int64(metaLen), 0)
			log.Printf("dial to node finished, idx: %d, node id: %s, block hash: %s, time elapsed: %s,  error: %#v \n", idx, node.Id, blockInfo.BlockHash, time.Now().Sub(dialStart), err)
			if err != nil {
				if dialRetryTimes < 2 {
					log.Printf("retrying, idx: %d, block hash: %s, new node id: %s, diaRetryTimes: %d\n", idx, blockInfo.BlockHash, node.Id, dialRetryTimes)
					dialRetryTimes++
					goto lazyTry
				}

				log.Printf("retry failed, return, idx: %d, diaRetryTimes: %d \n", idx, dialRetryTimes)
				return
			}
			if rc1 == nil {
				log.Printf("get nil reader from node, node down, idx: %d, node id: %s, block hash: %s\n", idx, node.Id, blockInfo.BlockHash)
				err = errors.New("ipfs node down")
				return
			}

			copyStart := time.Now()
			_, err = io.Copy(fhs[idx], rc1)
			log.Printf("block download finished, idx: %d, node id: %s, block hash: %s, time elapsed: %s, error: %#v \n", idx, node.Id, blockInfo.BlockHash, time.Now().Sub(copyStart), err)
			if err != nil {
				if downloadRetryTimes < 3 {
					log.Printf("retrying, idx: %d, downloadRetryTimes: %d \n", idx, downloadRetryTimes)
					downloadRetryTimes++
					fhs[idx].Seek(0, 0)
					goto lazyTry
				}
			}
			return
		}(i)
	}
	//wait
	for i := 0; i < c.BlockDownloadWorkerCount; i++ {
		sem <- true
	}
	return
}

func getRandonNode(nodes []NodeClient) NodeClient {
	rand.Seed(time.Now().UnixNano())
	i := rand.Intn(len(nodes))
	return nodes[i]
}

func (c *Client) StreamRead(fileHash string) (rc io.ReadCloser, metaAll metafile.Meta, err error) {
	blocksInfo, err := c.GetBlocksInfo(fileHash)
	if err != nil {
		return
	}

	meta, err := c.GetMeta(blocksInfo)
	if err != nil {
		return
	}

	dataShards := int(meta.MetaHeader.DataShards)
	_, _, shardSize := file.BlockCount(meta.FSize)
	lastShardSize := meta.FSize - int64(dataShards-1)*shardSize
	rcs := make([]io.ReadCloser, dataShards)
	retryTimes := 0
	var partSize int64 = 0
	for i := 0; i < dataShards; i++ {
		blockInfo := blocksInfo[i]
		nodes, e := c.GetNodeClients(blockInfo.PeerId)
		if e != nil {
			err = e
			log.Printf("GetNodeClient error: %s, node id: %s, block hash: %s \n", blockInfo.PeerId, blockInfo.BlockHash, err)
			return
		}

		if i == dataShards-1 {
			partSize = lastShardSize
		}

		node := nodes[0]
	lazyTry:
		rc1, e := ReadAt(node.Client, blockInfo.BlockHash, int64(meta.Len()), partSize)
		if e != nil {
			err = e
			if retryTimes < 3 {
				log.Printf("read failed, node id: %s, block hash: %s, error: %#v , retryTimes: %d, retrying\n", blockInfo.PeerId, blockInfo.BlockHash, err, retryTimes)
				node = getRandonNode(nodes)
				retryTimes++
				goto lazyTry
			}
			return
		}

		rcs[i] = rc1
		log.Printf("read block, node id: %s, Block Hash: %s, error: %s \n", node.Id, blockInfo.BlockHash, err)
	}

	rc = utils.MultiReadCloser(rcs...)
	return
}

func (c *Client) GetMeta(bis []storage.BlockInfo) (meta *metafile.Meta, err error) {
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

func getMeta(node *shell.Shell, blkHash string) (m *metafile.Meta, err error) {
	rc1, err := ReadAt(node, blkHash, 0, metafile.MetaHeaderLength)
	if err != nil {
		return
	}
	defer rc1.Close()
	metaHeaderB, err := ioutil.ReadAll(rc1)
	if err != nil {
		return
	}
	metaHeader, err := metafile.DecodeMetaHeader(metaHeaderB)
	if err != nil {
		return
	}

	metaLen := metafile.MetaHeaderLength + metaHeader.MetaBodyLength
	rc2, err := ReadAt(node, blkHash, 0, int64(metaLen))
	if err != nil {
		return
	}
	defer rc2.Close()
	metaData, err := ioutil.ReadAll(rc2)
	if err != nil {
		return
	}

	return metafile.DecodeMeta(metaData)
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

// Remove the given path
func FilesRm(s *shell.Shell, path string, recursive, force bool) error {
	return s.Request("files/rm", path).
		Option("recursive", recursive).
		Option("force", force).
		Exec(context.Background(), nil)
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

func (c *Client) GetNodeClient(nid string) (cli NodeClient, err error) {
	ns, err := c.GetNodeClients(nid)
	if err != nil {
		return
	}

	if len(ns) == 0 || ns[0].Id != nid {
		err = ErrNodeNotFound
	}

	cli = ns[0]

	return
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
	port := 4001
	available, _ := netools.IsLocalPortAvailable(port)
	if !available {
		port, err = netools.GetFreePort()
		if err != nil {
			return
		}
	}

	err = c.P2PForward(port, peerId)
	if err != nil {
		return
	}
	url := fmt.Sprintf("127.0.0.1:%d", port)

	cli = shell.NewShell(url)
	cli.SetTimeout(c.NodeRequestTimeout)
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
