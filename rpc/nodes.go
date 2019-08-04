package rpc

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"strconv"
	"time"

	"github.com/ipfs/go-ipfs-api"
	"github.com/ipweb-group/go-sdk/p2p"
	"github.com/ipweb-group/go-sdk/utils/netools"
	"github.com/ipweb-group/go-sdk/utils/reader"
	"github.com/libp2p/go-libp2p-peer"
)

var ErrNodeNotFound = errors.New("node not found")

type NodeStatus int

const (
	NodeStatusUnavailable = iota
	NodeStatusAvailable
	NodeStatusUsing
)

type Nodes map[string]Node

type Node struct {
	Id          string
	Status      NodeStatus
	CreateTime  time.Time
	UpdateTime  time.Time
	UploadBytes int64
	UploadDur   time.Duration
	Client      *shell.Shell
}

func (c Client) Add(r io.Reader) (id string, err error) {
	nr := reader.NewReader(r)

	return
}

func (c *Client) GetNode(nid string) (cli Node, err error) {
	ns, err := c.GetNodes()
	if err != nil {
		return
	}

	for i := range ns {
		if ns[i].Id == nid {
			cli = ns[i]
			return
		}
	}

	err = ErrNodeNotFound
	return
}

func (c *Client) GetNodes() (ns []Node, err error) {
	if len(c.Nodes) == 0 {
		err = c.refreshNodes()
	}

	if len(c.Nodes) == 0 {
		err = ErrNodeNotFound
	}

	for _, n := range c.Nodes {
		ns = append(ns, n)
	}

	return
}

func (c *Client) NewNode(peerId string) (n Node, err error) {
	n = Node{
		Id:         peerId,
		Status:     NodeStatusUnavailable,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}

	port, err := netools.GetFreePort()
	if err != nil {
		return
	}

	err = c.P2PForward(port, peerId)
	if err != nil {
		return
	}
	url := fmt.Sprintf("127.0.0.1:%d", port)

	cli := shell.NewShell(url)
	cli.SetTimeout(c.NodeRequestTimeout)
	info, err := cli.ID()
	if err != nil {
		c.P2PClose(0, peerId)
		fmt.Println("bad peer: ", peerId, " err: ", err)
		return
	}

	n.Status = NodeStatusAvailable
	n.Client = cli
	fmt.Println("p2p peer: ", peerId, " addr: ", info.Addresses)
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
	fmt.Println("nodes refreshing time: ", time.Now())

	ps := c.IpfsNode.Peerstore.Peers()
	for _, p := range ps {
		sema <- 1
		go func(peerId peer.ID) {
			defer func() {
				<-sema
			}()

			id := peerId.Pretty()
			if _, ok := c.Nodes[id]; ok {
				return
			}

			n, err := c.NewNode(id)
			c.NodesMux.Lock()
			c.Nodes[id] = n
			c.NodesMux.Unlock()
			if err != nil {
				c.P2PClose(0, id)
			}
			return
		}(p)
	}
	for i := 0; i < c.NodeRefreshWorkers; i++ {
		sema <- 1
	}

	if len(c.Nodes) == 0 {
		return ErrNodeNotFound
	}

	return nil
}

func getRandonNode(nodes []Node) Node {
	rand.Seed(time.Now().UnixNano())
	i := rand.Intn(len(nodes))
	return nodes[i]
}

func (c *Client) P2PForward(port int, peerId string) error {
	listenOpt := "/ip4/127.0.0.1/tcp/" + strconv.Itoa(port)
	targetOpt := "/ipfs/" + peerId
	return p2p.Forward(c.IpfsNode.Context(), c.IpfsNode.P2P, P2pProtocl, listenOpt, targetOpt)
}

func (c *Client) P2PClose(port int, peerId string) error {
	listenOpt := "/ip4/127.0.0.1/tcp/" + strconv.Itoa(port)
	targetOpt := "/ipfs/" + peerId
	return p2p.Close(c.IpfsNode.P2P, false, "", listenOpt, targetOpt)
}

func (c *Client) P2PCloseAll() error {
	return p2p.Close(c.IpfsNode.P2P, true, "", "", "")
}
