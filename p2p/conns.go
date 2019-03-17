package p2p

import (
	"context"
	"errors"
	"time"

	//"gx/ipfs/QmPJxxDsX2UbchSHobbYuvz7qnyJTFKvaKMzE2rZWJ4x5B/go-libp2p-peer"
	"github.com/ipfs/go-ipfs/p2p"
	ma "gx/ipfs/QmNTCey11oxhb1AxDnQBRHtdhap6Ctud872NjAYPYYXPuc/go-multiaddr"
	madns "gx/ipfs/QmQc7jbDUsxUJZyFJzxVrnrWeECCct6fErEpMqtjyWvCX8/go-multiaddr-dns"
	protocol "gx/ipfs/QmZNkThpqfVXs9GNbexPrfBbXSLNYeKrE7jwFM2oqHbyqN/go-libp2p-protocol"
	ipfsaddr "gx/ipfs/Qmba1WKLAf87nRgvKha56DZaAvN2tbTyyDz8XDDyzhvLKM/go-ipfs-addr"
)

var resolveTimeout = 10 * time.Second

func Forward(p *p2p.P2P, protoOpt, listenOpt, targetOpt string) error {

	proto := protocol.ID(protoOpt)

	listen, err := ma.NewMultiaddr(listenOpt)
	if err != nil {
		return err
	}

	targets, err := parseIpfsAddr(targetOpt)
	if err != nil {
		return err
	}

	_, err = p.ForwardLocal(context.Background(), targets[0].ID(), proto, listen)
	return err
}

func Close(p *p2p.P2P, closeAll bool, protoOpt, listenOpt, targetOpt string) error {

	proto := protocol.ID(protoOpt)

	listen, err := ma.NewMultiaddr(listenOpt)
	if err != nil {
		return err
	}

	target, err := ma.NewMultiaddr(targetOpt)
	if err != nil {
		return err
	}

	match := func(listener p2p.Listener) bool {
		if closeAll {
			return true
		}
		if proto != listener.Protocol() {
			return false
		}
		if !listen.Equal(listener.ListenAddress()) {
			return false
		}
		if !target.Equal(listener.TargetAddress()) {
			return false
		}

		return true
	}

	done := p.ListenersLocal.Close(match)
	done += p.ListenersP2P.Close(match)
	if done != 0 {
		err = errors.New("close failed")
	}

	return err
}

func parseIpfsAddr(addr string) ([]ipfsaddr.IPFSAddr, error) {
	mutiladdr, err := ma.NewMultiaddr(addr)
	if err != nil {
		return nil, err
	}
	if _, err := mutiladdr.ValueForProtocol(ma.P_IPFS); err == nil {
		iaddrs := make([]ipfsaddr.IPFSAddr, 1)
		iaddrs[0], err = ipfsaddr.ParseMultiaddr(mutiladdr)
		if err != nil {
			return nil, err
		}
		return iaddrs, nil
	}
	// resolve mutiladdr whose protocol is not ma.P_IPFS
	ctx, cancel := context.WithTimeout(context.Background(), resolveTimeout)
	addrs, err := madns.Resolve(ctx, mutiladdr)
	cancel()
	if len(addrs) == 0 {
		return nil, errors.New("fail to resolve the multiaddr:" + mutiladdr.String())
	}
	iaddrs := make([]ipfsaddr.IPFSAddr, len(addrs))
	for i, addr := range addrs {
		iaddrs[i], err = ipfsaddr.ParseMultiaddr(addr)
		if err != nil {
			return nil, err
		}
	}
	return iaddrs, nil
}
