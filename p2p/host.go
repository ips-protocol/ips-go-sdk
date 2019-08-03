package p2p

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"strconv"

	"github.com/ipfs/go-datastore"
	ci "github.com/libp2p/go-libp2p-crypto"
	peer "github.com/libp2p/go-libp2p-peer"

	"github.com/ipfs/go-ipfs/commands"
	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/repo"
	"github.com/ipweb-group/go-ipws-config"
	"github.com/ipweb-group/go-sdk/conf"
	"github.com/ipweb-group/go-sdk/utils/netools"
)

func NewNode(ctx context.Context, cfg conf.Config) (n *core.IpfsNode, err error) {
	ds := datastore.NewNullDatastore()
	repo, err := DefaultRepo(ds, cfg)
	if err != nil {
		return
	}

	ncfg := &core.BuildCfg{
		Repo:   repo, //opt
		Online: true, //required, true
	}

	n, err = core.NewNode(ctx, ncfg)
	return
}

func DefaultRepo(dstore repo.Datastore, cfg conf.Config) (repo.Repo, error) {
	c := config.Config{}
	priv, pub, err := ci.GenerateKeyPairWithReader(ci.RSA, 1024, rand.Reader)
	if err != nil {
		return nil, err
	}

	pid, err := peer.IDFromPublicKey(pub)
	if err != nil {
		return nil, err
	}

	privkeyb, err := priv.Bytes()
	if err != nil {
		return nil, err
	}

	bootstrapPeers, err := config.DefaultBootstrapPeers()
	if err != nil {
		return nil, err
	}
	c.Bootstrap = config.BootstrapPeerStrings(bootstrapPeers)

	port := 4001
	available, _ := netools.IsLocalPortAvailable(port)
	if !available {
		port, err = netools.GetFreePort()
		if err != nil {
			return nil, err
		}
	}
	c.Addresses.Swarm = []string{"/ip4/0.0.0.0/tcp/" + strconv.Itoa(port), "/ip6/::/tcp/" + strconv.Itoa(port)}
	c.Identity.PeerID = pid.Pretty()
	c.Identity.PrivKey = base64.StdEncoding.EncodeToString(privkeyb)
	// c.Discovery.MDNS.Enabled = true
	c.Discovery.MDNS.Interval = 1
	c.Routing.Type = "dht"

	c.Chain.URL = cfg.ContractConf.ContractNodeAddr
	// c.Chain.WalletPriKey = cfg.ContractConf.ClientKeyHex

	mockRepo := &repo.Mock{
		D: dstore,
		C: c,
	}

	return Repo{mockRepo}, nil
}

func Ctx(node *core.IpfsNode, repoPath string) commands.Context {
	return commands.Context{
		// Online:     true,
		ConfigRoot: repoPath,
		LoadConfig: func(path string) (*config.Config, error) {
			return node.Repo.Config()
		},
		ConstructNode: func() (*core.IpfsNode, error) {
			return node, nil
		},
	}
}
