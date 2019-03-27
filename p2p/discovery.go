package p2p

import (
	"context"
	"crypto/rand"
	"fmt"
	pstore "gx/ipfs/QmaCTz9RkrU13bm9kMB54f7atgqM4qkjDZpRwRoJiWXEqs/go-libp2p-peerstore"
	pstoremem "gx/ipfs/QmaCTz9RkrU13bm9kMB54f7atgqM4qkjDZpRwRoJiWXEqs/go-libp2p-peerstore/pstoremem"

	"gx/ipfs/QmRxk6AUaGaKCfzS1xSNRojiAPd7h2ih8GuCdjJBF3Y6GK/go-libp2p"

	circuit "gx/ipfs/QmZBfqr863PYD7BKbmCFSNmzsqYmtr2DKgzubsQaiTQkMc/go-libp2p-circuit"

	crypto "gx/ipfs/QmTW4SdgBWq9GjsBsHeUx8WuGxzhgzAf88UMH2w62PC8yK/go-libp2p-crypto"

	host "gx/ipfs/QmYrWiWM4qtrnCeT3R14jY3ZZyirDNJgwK57q4qFYePgbd/go-libp2p-host"

	"gx/ipfs/QmTZBfrPJmjWsCvHEtX5FE6KimVJhsJg5sBbqEFYf4UZtL/go-multiaddr"
)

type Config struct {
	RendezvousString string
	ProtocolID       string
	ListenHost       string
	ListenPort       int
}

func GetActivePeers() []pstore.PeerInfo {
	ps := []pstore.PeerInfo{}
	for _, p := range activePeers {
		ps = append(ps, p)
	}
	return ps
}

func RmFromActivePeers(p pstore.PeerInfo) {
	delete(activePeers, p.ID)
	return
}

func FindPeers(c Config) (peers chan pstore.PeerInfo, host2 host.Host, err error) {

	r := rand.Reader
	prvKey, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	if err != nil {
		return
	}
	ctx := context.Background()

	libp2pOpts := []libp2p.Option{}
	relayOpts := []circuit.RelayOpt{circuit.OptDiscovery}
	relayOpts = append(relayOpts, circuit.OptHop)
	libp2pOpts = append(libp2pOpts, libp2p.EnableRelay(relayOpts...))
	sourceMultiAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%d", c.ListenHost, c.ListenPort))

	libp2pOpts = append(libp2pOpts, []libp2p.Option{libp2p.ListenAddrs(sourceMultiAddr), libp2p.Identity(prvKey), libp2p.Peerstore(pstoremem.NewPeerstore())}...)

	host, err := libp2p.New(
		ctx,
		libp2pOpts...,
	)
	if err != nil {
		return
	}

	return initMDNS(ctx, host, c.RendezvousString), host, nil
}
