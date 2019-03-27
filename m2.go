package main

import (
	"context"
	"fmt"
	"github.com/ipfs/go-ipfs/core"
	"go-sdk/p2p"
	"go-sdk/rpc"
	"gx/ipfs/QmVSbopkxvLSRFuUn1SeHoEcArhCLn2okUbVpLvhQ1pm1X/interface-go-ipfs-core"
	ds "gx/ipfs/Qmf4xQhNomPNhrtZc67qSnfJSjxjXs9LWvknJtSXwimPrM/go-datastore"
	"os"

	rpc "./rpc"

	p2p "./p2p"

	"github.com/ipfs/go-ipfs/core"
)

func main() {

	ctx := context.Background()
	lds := ds.NewMapDatastore()
	repo, err := p2p.DefaultRepo(lds)
	if err != nil {
		panic(err)
	}

	key, err := repo.SwarmKey()
	fmt.Println("============>", string(key), err)
	ncfg := &core.BuildCfg{
		Repo: repo, //opt
		//Permanent:                 true, //opt, true|false
		Online: true, //required, true
		//Online: false, //required, true
		//DisableEncryptedConnections: false, //opt, false
		//ExtraOpts: map[string]bool{
		//	"pubsub": false, //opt, true|false
		//	"ipnsps": false, //opt, false|false
		//	"mplex":  true,  //opt,	true|false
		//}, //opt
		NilRepo: true,
	}

	n, err := core.NewNode(ctx, ncfg)
	if err != nil {
		panic(err)
	}

	cli, err := rpc.NewClient(n)
	if err != nil {
		panic(err)
	}

	f, err := os.Open("/tmp/test.txt")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	// cid, err := cli.Upload(f)
	cid, err := cli.Upload("/tmp/test.txt")
	if err != nil {
		panic(err)
	}

	// upload file ---------------------------------------------------------------------------------------------------
	//n, err := core.NewNode(ctx, ncfg)
	//if err != nil {
	//	panic(err)
	//}
	//
	//cli, err := rpc.NewClient(n)
	//if err != nil {
	//	panic(err)
	//}
	//
	//f, err := os.Open("/tmp/test.txt")
	//if err != nil {
	//	panic(err)
	//}
	//defer f.Close()
	//
	//cid, err := cli.Upload(f)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println("-------->", cid)

	// ---------------------------------------------------------------------------------------------------
	//n, err := core.NewNode(ctx, ncfg)
	//if err != nil {
	//	panic(err)
	//}
	//
	//for {
	//	ps := n.Peerstore.Peers()
	//	for _, p := range ps {
	//		fmt.Println("peer: ---->", p.Pretty())
	//
	//		remotePeer := "Qmain1GGsLNtPmDPJsmWGYv7QxyFnbTjvFceH16yC2PCRd"
	//		if p.Pretty() == remotePeer {
	//			p2p.Forward(n.P2P, "/sys/http", "/ip4/127.0.0.1/tcp/8888", "/ipfs/Qmain1GGsLNtPmDPJsmWGYv7QxyFnbTjvFceH16yC2PCRd")
	//		}
	//	}
	//
	//	cns := n.PeerHost.Network().Conns()
	//	for _, cn := range cns {
	//		streams := cn.GetStreams()
	//		fmt.Println("stream:", streams)
	//	}
	//}

	// cacl cid ---------------------------------------------------------------------------------------------------
	n, err := core.NewNode(ctx, ncfg)
	if err != nil {
		panic(err)
	}
	ctx1 := p2p.Ctx(n, "")
	api, err := ctx1.GetAPI()
	if err != nil {
		panic(err)
	}

	fn, err := api.Unixfs().Add()
	if err != nil {
		panic(err)
	}
}
