package main

import (
	"fmt"
	"github.com/ipweb-group/go-sdk/conf"
	"github.com/ipweb-group/go-sdk/rpc"
)

func main() {

	ccfg := conf.ContractConfig{
		ClientKeyHex:     "B2FE66D78810869A64CAAE7B1F2C60CCA3AC2F2261DA2F1DE7040DE3F1FEDA9C",
		ContractNodeAddr: "https://mainnet.ipweb.top",
	}
	cfg := conf.Config{ContractConf: ccfg, BlockUploadWorkers: 3}

	cli, err := rpc.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	// upload file ---------------------------------------------------------------------------------------------------
	cid, err := cli.UploadWithPath("/tmp/abc.txt")
	if err != nil {
		panic(err)
	}
	fmt.Println("Upload success cid:", cid)

	// delete file ---------------------------------------------------------------------------------------------------
	//err = cli.Remove("QmaP7refemdjCt55jWrbNdepkoRfvabkvQKK98iUdt7JKb")
	//if err != nil {
	//	panic(err)
	//} else {
	//	fmt.Println("delete file success!")
	//}

	//
	// download file ---------------------------------------------------------------------------------------------------
	//rc, meta, err := cli.Download("QmVQfe4DN8oAFRiTy9PLuCNhRxcbR6cYFc2XhJQJfwMc9d")
	//if err != nil && err != io.EOF {
	//	panic(err)
	//}
	//defer rc.Close()
	//fc, err := ioutil.ReadAll(rc)
	//if err != nil {
	//	panic(err)
	//}
	//fmt.Println("Download file content:", string(fc), "\tmeta data:", meta)

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
