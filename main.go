package main

import (
	"fmt"
	pstore "gx/ipfs/QmaCTz9RkrU13bm9kMB54f7atgqM4qkjDZpRwRoJiWXEqs/go-libp2p-peerstore"
	"time"

	p2p "./p2p"
)

func main1() {
	cfg := p2p.Config{
		RendezvousString: "_ipfs-discovery._udp",
		ProtocolID:       "test/abc",
		ListenHost:       "0.0.0.0",
		ListenPort:       9999,
	}

	ps, host, err := p2p.FindPeers(cfg)
	if err != nil {
		panic(err)
	}

	peers := []pstore.PeerInfo{}
	go func() {
		pm := map[string]pstore.PeerInfo{}
		for p := range ps {
			if _, ok := pm[p.ID.String()]; !ok {
				fmt.Println("find peer:", p, p.ID.Pretty())
				pm[p.ID.String()] = p
				peers = append(peers, p)
			} else {

			}
		}
	}()

	for {
		fmt.Println("==>:", len(peers))
		fmt.Println("==>:", len(p2p.GetActivePeers()))
		fmt.Println("net work peers ==>:", host.Network().Conns())
		time.Sleep(time.Second)
	}
	select {} //wait here
}
