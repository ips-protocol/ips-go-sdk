package p2p

import (
	"context"
	"time"

	"github.com/libp2p/go-libp2p-host"
	"github.com/libp2p/go-libp2p-peer"
	pstore "github.com/libp2p/go-libp2p-peerstore"
	"github.com/libp2p/go-libp2p/p2p/discovery"
)

var activePeers map[peer.ID]pstore.PeerInfo

func init() {
	activePeers = make(map[peer.ID]pstore.PeerInfo)
}

type discoveryNotifee struct {
	PeerChan chan pstore.PeerInfo
}

//interface to be called when new  peer is found
func (n *discoveryNotifee) HandlePeerFound(pi pstore.PeerInfo) {
	n.PeerChan <- pi
	if _, ok := activePeers[pi.ID]; !ok {
		activePeers[pi.ID] = pi
	}
}

//Initialize the MDNS service
func initMDNS(ctx context.Context, peerhost host.Host, rendezvous string) chan pstore.PeerInfo {
	ser, err := discovery.NewMdnsService(ctx, peerhost, time.Second, rendezvous)
	if err != nil {
		panic(err)
	}

	n := &discoveryNotifee{}
	n.PeerChan = make(chan pstore.PeerInfo)

	ser.RegisterNotifee(n)
	return n.PeerChan
}
