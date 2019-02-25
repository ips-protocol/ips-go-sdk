package p2p

import (
	"context"
	"crypto/rand"
	"fmt"
	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p-crypto"
	pstore "github.com/libp2p/go-libp2p-peerstore"
	"github.com/multiformats/go-multiaddr"
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

func FindPeers(c Config) (peers chan pstore.PeerInfo, err error) {

	r := rand.Reader
	prvKey, _, err := crypto.GenerateKeyPairWithReader(crypto.RSA, 2048, r)
	if err != nil {
		return
	}
	ctx := context.Background()
	sourceMultiAddr, _ := multiaddr.NewMultiaddr(fmt.Sprintf("/ip4/%s/tcp/%d", c.ListenHost, c.ListenPort))
	host, err := libp2p.New(
		ctx,
		libp2p.ListenAddrs(sourceMultiAddr),
		libp2p.Identity(prvKey),
	)
	if err != nil {
		return
	}

	return initMDNS(ctx, host, c.RendezvousString), nil
}
