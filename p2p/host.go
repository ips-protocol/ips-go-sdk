package p2p

import (
	"context"
	"crypto/rand"
	"encoding/base64"

	"github.com/ipfs/go-ipfs/commands"
	"github.com/ipfs/go-ipfs/config"
	"github.com/ipfs/go-ipfs/core"
	"github.com/ipfs/go-ipfs/repo"
	"github.com/ipweb-group/go-sdk/p2p"
	ci "gx/ipfs/QmTW4SdgBWq9GjsBsHeUx8WuGxzhgzAf88UMH2w62PC8yK/go-libp2p-crypto"
	"gx/ipfs/QmUadX5EcvrBmxAV9sE7wUWtWSqxns5K84qKJBixmcT1w9/go-datastore"
	peer "gx/ipfs/QmYVXrKrKHDC9FobgmcmshCDyWwdrfwfanNQN4oxJ9Fk3h/go-libp2p-peer"
)

func NewNode(ctx context.Context) (n *core.IpfsNode, err error) {
	ds := datastore.NewNullDatastore()
	repo, err := p2p.DefaultRepo(ds)
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

func DefaultRepo(dstore repo.Datastore) (repo.Repo, error) {
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

	c.Bootstrap = config.DefaultBootstrapAddresses
	c.Addresses.Swarm = []string{"/ip4/0.0.0.0/tcp/4001"}
	c.Identity.PeerID = pid.Pretty()
	c.Identity.PrivKey = base64.StdEncoding.EncodeToString(privkeyb)
	c.Discovery.MDNS.Enabled = true
	c.Discovery.MDNS.Interval = 1
	c.Routing.Type = "dhtclient"

	mockRepo := &repo.Mock{
		D: dstore,
		C: c,
	}

	return Repo{mockRepo}, nil
}

func Ctx(node *core.IpfsNode, repoPath string) commands.Context {
	return commands.Context{
		Online:     true,
		ConfigRoot: repoPath,
		LoadConfig: func(path string) (*config.Config, error) {
			return node.Repo.Config()
		},
		ConstructNode: func() (*core.IpfsNode, error) {
			return node, nil
		},
	}
}
