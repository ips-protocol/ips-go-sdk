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
	NodeStatusUnavailable = iota // unavailable
	NodeStatusOpen               // available
	NodeStatusClosed             // available
)

type Nodes []*Node

type Node struct {
	Id              string
	Port            int
	Status          NodeStatus
	UpdateTime      time.Time
	AccessTime      time.Time
	UploadBytes     int64
	UploadDur       time.Duration
	ManualSetWeight int
	SuccessedTimes  int
	FailedTimes     int
	ConnQuota       int
	Client          *shell.Shell
}

func (ns Nodes) Len() int {
	return len(ns)
}

func (ns Nodes) Less(i, j int) bool {
	return ns[i].ConnQuota > ns[j].ConnQuota
}

func (ns Nodes) Swap(i, j int) {
	ns[i], ns[j] = ns[j], ns[i]
}

func (ns Nodes) Sort() {
	sort.Sort(ns)
}

func (c Client) Add(r io.Reader) (id string, err error) {
	n, err := c.NodeByManulWeight()
	defer func() {
		if n == nil {
			fmt.Printf("get upload node failed: %+v\n", err)
			return
		}

		if err != nil {
			c.NodesMux[n.Id].Lock()
			n.FailedTimes++
			failedRate := float64(n.FailedTimes) / float64(n.FailedTimes+n.SuccessedTimes)
			if failedRate >= 0.2 {
				n.Status = NodeStatusUnavailable
			}
			c.NodesMux[n.Id].Unlock()

			fmt.Printf("upload failed id: %s, port: %d, failed times: %d, success times: %d, err: %+v\n",
				n.Id, n.Port, n.FailedTimes, n.SuccessedTimes, err)
		}

		c.NodesAllocCond.L.Lock()
		c.Nodes[n.Id].ConnQuota += 1
		c.NodesAllocCond.L.Unlock()
		c.NodesAllocCond.Broadcast()

		fmt.Printf("upload node id: %s, port: %d, block hash: %s, error: %+v\n", n.Id, n.Port, id, err)
	}()
	if err != nil {
		return
	}
	fmt.Printf("get available node id: %s, port: %d, remain connects quota: %d\n", n.Id, n.Port, n.ConnQuota)

	nr := reader.NewReader(r)
	start := time.Now()
	id, err = n.Client.Add(nr)
	if err != nil {
		return
	}

	c.NodesMux[n.Id].Lock()
	n.UploadDur += time.Now().Sub(start)
	n.UploadBytes += nr.N()
	n.SuccessedTimes++
	c.NodesMux[n.Id].Unlock()

	return
}

func (c *Client) RandomNode() (n *Node, err error) {
	ns, err := c.GetOpenNodes()
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
	c.NodesAllocCond.L.Lock()
	defer func() {
		if err == nil {
			n.ConnQuota -= 1
		}
		c.NodesAllocCond.L.Unlock()
	}()

reTry:
	ns, err := c.GetAvailableNodes()
	if err != nil {
		return
	}
	if len(ns) == 0 {
		err = ErrNodeNotFound
		return
	}

	Nodes(ns).Sort()
	if ns[0].ConnQuota <= 0 {
		c.NodesAllocCond.Wait()
		goto reTry
	}

	var availNodes []*Node
	connQuotaFlag := ns[0].ConnQuota
	for i := range ns {
		if ns[i].ConnQuota < connQuotaFlag {
			break
		}

		availNodes = append(availNodes, ns[i])
	}

	var weightSum int
	for i := range availNodes {
		weightSum += availNodes[i].ManualSetWeight
	}
	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(weightSum)
	for i := range availNodes {
		r -= availNodes[i].ManualSetWeight
		if r <= 0 {
			n = availNodes[i]
			if n.Status == NodeStatusClosed {
				n, err = c.NewNode(n.Id)
			}
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

func (c *Client) GetOpenNodes() (ns []*Node, err error) {
	return c.GetNodes(NodeStatusOpen)
}

func (c *Client) GetAvailableNodes() (ns []*Node, err error) {
	cNodes, err := c.GetNodes(NodeStatusClosed)
	if err != nil {
		return
	}

	aNodes, err := c.GetNodes(NodeStatusOpen)
	if err != nil {
		return
	}

	ns = append(aNodes, cNodes...)
	return
}

func (c *Client) GetNodes(status NodeStatus) (ns []*Node, err error) {
	if len(c.Nodes) == 0 {
		err = c.refreshNodes(1)
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
		Id:              peerId,
		Status:          NodeStatusUnavailable,
		ManualSetWeight: 1,
		AccessTime:      time.Now(),
		UpdateTime:      time.Now(),
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
	_, err = cli.ID()
	if err != nil {
		c.P2PClose(port, peerId)
		fmt.Println("bad peer: ", peerId, "port: ", port, " err: ", err)
		return
	}

	n.Port = port
	n.Status = NodeStatusOpen
	n.Client = cli
	fmt.Println("p2p peer: ", peerId, " port: ", port)
	return
}

func (c *Client) CloseNode(id string) error {
	c.NodesMux[id].Lock()
	defer c.NodesMux[id].Unlock()

	n := c.Nodes[id]
	if n.Status != NodeStatusOpen || n.ConnQuota != c.ConnQuotaPerNode {
		return nil
	}

	err := c.P2PClose(n.Port, n.Id)
	if err != nil {
		return err
	}
	n.Status = NodeStatusClosed
	return nil
}

func (c *Client) OpenNode(id string) error {
	c.NodesMux[id].Lock()
	defer c.NodesMux[id].Unlock()

	n := c.Nodes[id]
	if n.Status != NodeStatusClosed {
		return nil
	}

	n.ConnQuota = c.ConnQuotaPerNode
	err := c.P2PForward(n.Port, n.Id)
	if err != nil {
		return err
	}
	n.Status = NodeStatusOpen
	return nil
}

func (c *Client) refreshNodesTick() {
	err := c.refreshNodes(c.NodeRefreshWorkers)
	if err != nil {
		fmt.Println("refreshNodes err: ", err)
	}
	for {
		select {
		case <-time.Tick(c.NodeRefreshDuration):
			err = c.refreshNodes(1)
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
				dur := time.Now().Sub(n.AccessTime.Add(c.NodeCloseDuration))
				if n.Status == NodeStatusOpen && dur.Seconds() > 0 {
					fmt.Printf("close node id: %s, access time: %s", n.Id, n.AccessTime)
					n.Status = NodeStatusClosed
					c.P2PClose(n.Port, n.Id)
				}
				c.NodesMux[n.Id].Unlock()
			}
		}
	}
}

func (c *Client) refreshNodes(workers int) error {
	c.NodeRefreshTime = time.Now()
	sema := make(chan int, workers)
	fmt.Println("nodes refreshing time: ", time.Now())

	ps := c.IpfsNode.Peerstore.Peers()
	for _, p := range ps {
		sema <- 1
		go func(peerId peer.ID) {
			defer func() {
				<-sema
			}()

			id := peerId.Pretty()
			c.NodesRwMux.RLock()
			n, ok := c.Nodes[id]
			c.NodesRwMux.RUnlock()
			if ok && n.Status != NodeStatusUnavailable {
				return
			}

			n, _ = c.NewNode(id)
			n.ConnQuota = c.ConnQuotaPerNode
			c.NodesRwMux.Lock()
			c.NodesMux[id] = &sync.RWMutex{}
			c.Nodes[id] = n
			c.NodesRwMux.Unlock()

			if w, ok := c.NodesWeightInfo[id]; ok {
				if w < 0 {
					w = 0
				}
				n.ManualSetWeight = w + 1
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
	return p2p.Close(c.IpfsNode.P2P, false, P2pProtocl, listenOpt, targetOpt)
}

func (c *Client) P2PCloseAll() error {
	return p2p.Close(c.IpfsNode.P2P, true, P2pProtocl, "", "")
}
