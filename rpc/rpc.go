package rpc

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"sync"
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
	NodesAccLock         sync.RWMutex
	*core.IpfsNode
}

func NewClient(ctx context.Context, cfg *core.BuildCfg) (cli *Client, err error) {
	n, err := core.NewNode(ctx, cfg)
	if err != nil {
		return
	}
	cli = &Client{IpfsNode: n}
	return
}

func (c *Client) Upload(r io.Reader) (cid string, err error) {
	ss, err := c.GetIpfsClients()
	if err != nil {
		return
	}

	cid, err = ss[0].Add(r)
	return
}

func (c *Client) GetIpfsClients() (ss []*shell.Shell, err error) {
	addrs, err := c.GetLocalAddrs()
	if err == ErrNodeNotFound {
		time.Sleep(time.Second * 3)
		addrs, err = c.GetLocalAddrs()
	}
	if err != nil {
		return
	}

	for i := range addrs {
		s := shell.NewShell(addrs[i])
		ss = append(ss, s)
	}

	return
}

func (c *Client) GetLocalAddrs() (addrs map[string]string, err error) {

	ps := c.IpfsNode.Peerstore.Peers()
	for _, p := range ps {
		fmt.Println("peer: ---->", p.Pretty())
		port, err := GetFreePort()
		if err != nil {
			return nil, err
		}

		id := p.Pretty()
		if isP2PNode(id) {
			err = c.P2PForward(port, id)
			if err != nil {
				log.Println(err)
				continue
			}
			addrs[id] = localAddr(port)
		}
	}

	if len(addrs) == 0 {
		err = ErrNodeNotFound
	}

	return
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
