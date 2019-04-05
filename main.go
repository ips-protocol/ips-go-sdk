package main

import (
	"context"
	"fmt"
	"github.com/ipweb-group/go-sdk/conf"
	"github.com/ipweb-group/go-sdk/p2p"
	"github.com/ipweb-group/go-sdk/rpc"
)

func main() {

	ctx := context.Background()
	n, err := p2p.NewNode(ctx)
	if err != nil {
		panic(err)
	}

	ccfg := conf.ContractConfig{
		ClientKeyHex:     "92D38B6F671F575EC9E47102364F53CA7F75B706A43606AA570E53917CBE2F9C",
		ContractNodeAddr: "http://127.0.0.1:8545",
	}
	cfg := conf.Config{ContractConfig: ccfg, BlockUpWorkerCount: 3}

	cli, err := rpc.NewClient(cfg, n)
	if err != nil {
		panic(err)
	}

	cid, err := cli.Upload("/tmp/5m.txt")
	if err != nil {
		panic(err)
	}
	fmt.Println("Upload success cid:", cid)

	//rc, meta, err := cli.Download("")
	//if err != nil && err != io.EOF {
	//	panic(err)
	//}
	//defer rc.Close()
	//fc, err := ioutil.ReadAll(rc)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println("Download file content:", string(fc), "\tmeta data:", meta)

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
