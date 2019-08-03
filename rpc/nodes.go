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
	"github.com/libp2p/go-libp2p-peer"
)

var ErrNodeNotFound = errors.New("node not found")

func (c Client) Add(r io.Reader) (id string, err error) {
	return
}

func (c Client) GetAvaiab() {

}

type Node struct {
	Id     string
	Client *shell.Shell
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
	if len(c.IpfsClients) == 0 {
		err = c.refreshNodes()
	}

	for id, n := range c.IpfsClients {
		ns = append(ns, Node{id, n})
	}

	if len(ns) == 0 {
		err = ErrNodeNotFound
	}

	return
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
	cli.SetTimeout(c.NodeRequestTimeout)
	info, err := cli.ID()
	if err != nil {
		c.P2PClose(0, peerId)
		fmt.Println("bad peer: ", peerId, " err: ", err)
		return
	}

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
				c.IpfsUnavailableClientsMux.Lock()
				c.IpfsUnavailableClients[id] = cli
				c.IpfsUnavailableClientsMux.Unlock()

				c.P2PClose(0, id)
				return
			}

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
