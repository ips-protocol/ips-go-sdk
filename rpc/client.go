package rpc

import (
	"context"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"strconv"
	"time"

	"github.com/ipfs/go-ipfs-api"
	"github.com/ipfs/go-ipfs/core"
	"github.com/ipweb-group/go-sdk/conf"
	"github.com/ipweb-group/go-sdk/contracts/storage"
	"github.com/ipweb-group/go-sdk/p2p"
	"github.com/ipweb-group/go-sdk/utils/netools"
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

func getRandonNode(nodes []NodeClient) NodeClient {
	rand.Seed(time.Now().UnixNano())
	i := rand.Intn(len(nodes))
	return nodes[i]
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
