package rpc

import (
	"errors"
	"fmt"
	"io"
	"math/rand"
	"sort"
	"strconv"
	"sync"
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
	NodeStatusClosed
)

type Nodes []Node

type Node struct {
	Id              string
	Port            int
	Status          NodeStatus
	CreateTime      time.Time
	UpdateTime      time.Time
	UploadBytes     int64
	UploadDur       time.Duration
	ManualSetWeight int
	SuccessedTimes  int
	FailedTimes     int
	Client          *shell.Shell
}

func (ns Nodes) Len() int {
	return len(ns)
}

func (ns Nodes) Less(i, j int) bool {
	return ns[i].Id < ns[j].Id
}

func (ns Nodes) Swap(i, j int) {
	ns[i], ns[j] = ns[j], ns[i]
}

func (ns Nodes) Sort() {
	sort.Sort(ns)
}

func (c Client) Add(r io.Reader) (id string, err error) {
	n, err := c.NodeByManulWeight()
	if err != nil {
		return
	}
	defer func() {
		fmt.Printf("upload node id: %s, conn port: %d, block hash: %s, err: %+v\n", n.Id, n.Port, id, err)
		c.P2PClose(n.Port, n.Id)
	}()

	nr := reader.NewReader(r)
	start := time.Now()
	id, err = n.Client.Add(nr)
	if err != nil {
		c.NodesMux[n.Id].Lock()
		n.FailedTimes++
		c.NodesMux[n.Id].Unlock()
		return
	}

	c.NodesMux[n.Id].Lock()
	n.UploadDur += time.Now().Sub(start)
	n.UploadBytes += nr.N()
	n.SuccessedTimes++
	n.UpdateTime = time.Now()
	c.NodesMux[n.Id].Unlock()

	return
}

func (c *Client) RandomNode() (n *Node, err error) {
	ns, err := c.GetAvailableNodes()
	if err != nil {
		return
	}

	var speedSum int
	for i := range ns {
		if ns[i].UploadBytes == 0 {
			continue
		}

		speedSum += int(float64(ns[i].UploadBytes) / ns[i].UploadDur.Seconds())
	}

	baseWeight := speedSum / len(ns)
	rand.Seed(time.Now().UnixNano())

	randLimit := speedSum * 2
	if randLimit == 0 {
		randLimit = 1
	}
	r := rand.Intn(randLimit)
	for i := range ns {
		speed := 0
		if ns[i].UploadBytes != 0 {
			speed = int(float64(ns[i].UploadBytes) / ns[i].UploadDur.Seconds())
		}
		r -= baseWeight
		r -= speed
		if r <= 0 {
			n = ns[i]
			if n.Status == NodeStatusClosed {
				c.NodesMux[n.Id].Lock()
				n, err = c.NewNode(n.Id)
				c.NodesMux[n.Id].Unlock()
			}
			return
		}
	}

	return
}

func (c *Client) NodeByManulWeight() (n *Node, err error) {
	ns, err := c.GetAvailableOrClosedNodes()
	if err != nil {
		return
	}

	var weightSum int
	for i := range ns {
		weightSum += ns[i].ManualSetWeight
	}

	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(weightSum)
	for i := range ns {
		r -= ns[i].ManualSetWeight
		if r <= 0 {
			n, err = c.NewNode(ns[i].Id)
			return
		}
	}

	return
}

func (c *Client) GetNode(nid string) (n *Node, err error) {
	ns, err := c.GetAvailableNodes()
	if err != nil {
		return
	}

	for i := range ns {
		if ns[i].Id == nid {
			n = ns[i]
			if n.Status == NodeStatusClosed {
				c.NodesMux[n.Id].Lock()
				n, err = c.NewNode(n.Id)
				c.NodesMux[n.Id].Unlock()
			}
			return
		}
	}

	err = ErrNodeNotFound
	return
}

func (c *Client) GetAvailableNodes() (ns []*Node, err error) {
	return c.GetNodes(NodeStatusAvailable)
}

func (c *Client) GetAvailableOrClosedNodes() (ns []*Node, err error) {
	cNodes, err := c.GetNodes(NodeStatusClosed)
	if err != nil {
		return
	}

	aNodes, err := c.GetNodes(NodeStatusAvailable)
	if err != nil {
		return
	}

	ns = append(aNodes, cNodes...)
	return
}

func (c *Client) GetNodes(status NodeStatus) (ns []*Node, err error) {
	if len(c.Nodes) == 0 {
		err = c.refreshNodes()
	}

	if len(c.Nodes) == 0 {
		err = ErrNodeNotFound
	}

	for _, n := range c.Nodes {
		if n.Status == status {
			ns = append(ns, n)
		}
	}

	return
}

func (c *Client) NewNode(peerId string) (n *Node, err error) {
	n = &Node{
		Id:         peerId,
		Status:     NodeStatusUnavailable,
		CreateTime: time.Now(),
		UpdateTime: time.Now(),
	}

	port, err := netools.GetFreePort()
	if err != nil {
		return
	}
	n.Port = port

	err = c.P2PForward(port, peerId)
	if err != nil {
		return
	}

	url := fmt.Sprintf("127.0.0.1:%d", port)
	cli := shell.NewShell(url)
	cli.SetTimeout(c.NodeRequestTimeout)
	n.Client = cli
	n.UpdateTime = time.Now()

	return
}

func (c *Client) NewNodeAndVerify(peerId string) (n *Node, err error) {
	n, err = c.NewNode(peerId)
	if err != nil {
		return
	}

	info, err := n.Client.ID()
	if err != nil {
		c.P2PClose(n.Port, peerId)
		fmt.Println("bad peer: ", peerId, " err: ", err)
		return
	}
	n.Status = NodeStatusAvailable

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

func (c *Client) closeNodesTick() {
	fmt.Println("closeNodes tick.")
	for {
		select {
		case <-time.Tick(c.NodeCloseDuration):
			for _, n := range c.Nodes {
				c.NodesMux[n.Id].Lock()
				dur := time.Now().Sub(n.UpdateTime.Add(c.NodeCloseDuration))
				if n.Status == NodeStatusAvailable && dur.Seconds() > 0 {
					n.Status = NodeStatusClosed
					c.P2PClose(0, n.Id)
				}
				c.NodesMux[n.Id].Unlock()
			}
		}
	}
}

func (c *Client) refreshNodes() error {
	c.NodeRefreshTime = time.Now()
	sema := make(chan int, c.NodeRefreshWorkers)
	fmt.Println("nodes refreshing time: ", time.Now())

	nodeAccMux := sync.RWMutex{}
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

			n, err := c.NewNodeAndVerify(id)
			nodeAccMux.Lock()
			c.Nodes[id] = n
			c.NodesMux[id] = &sync.RWMutex{}
			nodeAccMux.Unlock()
			if err != nil {
				c.P2PClose(0, id)
			}

			if w, ok := c.NodesWeightInfo[id]; ok {
				n.ManualSetWeight += w + 1
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
