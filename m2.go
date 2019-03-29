package main

import (
	"context"
	"fmt"
	"github.com/ipfs/go-ipfs/core"
	"go-sdk/conf"
	"go-sdk/p2p"
	"go-sdk/rpc"
	"gx/ipfs/QmUadX5EcvrBmxAV9sE7wUWtWSqxns5K84qKJBixmcT1w9/go-datastore"
)

func main() {

	ctx := context.Background()

	ds := datastore.NewNullDatastore()
	repo, err := p2p.DefaultRepo(ds)
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
	}

	n, err := core.NewNode(ctx, ncfg)
	if err != nil {
		panic(err)
	}

	ccfg := conf.ContractConfig{
		ClientKeyHex:       "92D38B6F671F575EC9E47102364F53CA7F75B706A43606AA570E53917CBE2F9C",
		StorageKeyHex:      "CED8FF231B09B14F09D8FF977C5C6C079EF4B485FC2A0D3B2955182B77310A04",
		ContractNodeAddr:   "http://127.0.0.1:8545",
		TransactorGasLimit: 8000,
		TransactorGasPrice: 1,
		TransactorValue:    1e6,
	}
	cfg := conf.Config{ContractConfig: ccfg}

	cli, err := rpc.NewClient(cfg, n)
	if err != nil {
		panic(err)
	}

	// cid, err := cli.Upload(f)
	cid, err := cli.Upload("/tmp/test.txt")
	if err != nil {
		panic(err)
	}
	fmt.Println("cid:", cid)

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

	// list nodes ---------------------------------------------------------------------------------------------------
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
	//n, err := core.NewNode(ctx, ncfg)
	//if err != nil {
	//	panic(err)
	//}
	//ctx1 := p2p.Ctx(n, "")
	//api, err := ctx1.GetAPI()
	//if err != nil {
	//	panic(err)
	//}
	//
	//fn, err := api.Unixfs().Add()
	//if err != nil {
	//	panic(err)
	//}
}
