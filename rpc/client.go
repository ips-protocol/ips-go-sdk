package rpc

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"github.com/ipfs/go-ipfs-api"
	"github.com/ipfs/go-ipfs/core"
	"github.com/ipweb-group/go-sdk/conf"
	"github.com/ipweb-group/go-sdk/contracts/storage"
	"github.com/ipweb-group/go-sdk/p2p"
	"github.com/ipweb-group/go-sdk/utils/netools"
	"github.com/libp2p/go-libp2p-peer"
)

var ErrNodeNotFound = errors.New("node not found")
var ErrContractNotFound = errors.New("no contract code at given address")

const P2pProtocl = "/sys/http"

type Client struct {
	IpfsClients               map[string]*shell.Shell
	IpfsClientsMux            sync.RWMutex
	IpfsUnavailableClients    map[string]*shell.Shell
	IpfsUnavailableClientsMux sync.RWMutex
	NodeRefreshTime           time.Time
	NodeRefreshDuration       time.Duration
	NodeRequestTimeout        time.Duration
	NodeRefreshWorkers        int
	BlockUploadWorkers        int
	BlockDownloadWorkers      int
	WalletPubKey              string
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
	cli.IpfsClientsMux = sync.RWMutex{}
	cli.IpfsUnavailableClients = make(map[string]*shell.Shell)
	cli.IpfsUnavailableClientsMux = sync.RWMutex{}
	if cfg.NodeRefreshIntervalInSecond == 0 {
		cfg.NodeRefreshIntervalInSecond = 600
	}
	cli.NodeRefreshDuration = time.Second * time.Duration(cfg.NodeRefreshIntervalInSecond)

	if cfg.NodeRefreshWorkers == 0 {
		cfg.NodeRefreshWorkers = 10
	}
	cli.NodeRefreshWorkers = cfg.NodeRefreshWorkers

	if cfg.NodeRequestTimeoutInSecond == 0 {
		cfg.NodeRequestTimeoutInSecond = 60
	}
	cli.NodeRequestTimeout = time.Second * time.Duration(cfg.NodeRequestTimeoutInSecond)

	if cfg.BlockUploadWorkers == 0 {
		cfg.BlockUploadWorkers = 5
	}
	cli.BlockUploadWorkers = cfg.BlockUploadWorkers

	if cfg.BlockDownloadWorkers == 0 {
		cfg.BlockDownloadWorkers = 5
	}
	cli.BlockDownloadWorkers = cfg.BlockDownloadWorkers

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
	go cli.refreshNodesTick()

	return
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
	if len(c.IpfsClients) != 0 {
		ns = getNodes()
		return
	}

	err = c.refreshNodes()
	if err == ErrNodeNotFound {
		err = c.refreshNodes()
	}
	if err != nil {
		return
	}

	ns = getNodes()
	return
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

func (c *Client) refreshNodesTick() {
	err := c.refreshNodes()
	if err != nil {
		fmt.Println("refreshNodes err: ", err)
	}
	for {
		select {
		case <-time.Tick(c.NodeRefreshDuration):
			err = c.refreshNodes()
			if err != nil {
				fmt.Println("refreshNodes err: ", err)
			}
		}
	}
}

func (c *Client) refreshNodes() error {
	c.NodeRefreshTime = time.Now()
	sema := make(chan int, c.NodeRefreshWorkers)

	ps := c.IpfsNode.Peerstore.Peers()
	for _, p := range ps {
		sema <- 1
		go func(peerId peer.ID) {
			defer func() {
				<-sema
			}()

			id := peerId.Pretty()

			c.IpfsClientsMux.RLock()
			cli, ok := c.IpfsClients[id]
			c.IpfsClientsMux.RUnlock()
			if ok {
				return
			}

			c.IpfsUnavailableClientsMux.RLock()
			_, ok = c.IpfsUnavailableClients[id]
			c.IpfsUnavailableClientsMux.RUnlock()
			if ok {
				return
			}

			cli, err := c.NewIpfsClient(id)
			if err != nil {
				fmt.Println("bad peer: ", p.Pretty(), err)

				c.IpfsUnavailableClientsMux.Lock()
				c.IpfsUnavailableClients[id] = cli
				c.IpfsUnavailableClientsMux.Unlock()

				c.P2PClose(0, id)
				return
			}

			fmt.Println("p2p peer: ", p.Pretty())
			c.IpfsClientsMux.Lock()
			c.IpfsClients[id] = cli
			c.IpfsClientsMux.Unlock()
		}(p)
	}
	for i := 0; i < c.NodeRefreshWorkers; i++ {
		sema <- 1
	}

	if len(c.IpfsClients) == 0 {
		return ErrNodeNotFound
	}

	return nil
}

func (c *Client) addClient() {

}

func getRandonNode(nodes []NodeClient) NodeClient {
	rand.Seed(time.Now().UnixNano())
	i := rand.Intn(len(nodes))
	return nodes[i]
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
