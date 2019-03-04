package main

import (
	"context"
	"fmt"
	"github.com/ipfs/go-ipfs/core"
	"go-sdk/p2p"
	"gx/ipfs/QmQmhotPUzVrMEWNK3x1R5jQ5ZHWyL7tVUrmRPjrBrvyCb/go-ipfs-files"
	"gx/ipfs/QmVSbopkxvLSRFuUn1SeHoEcArhCLn2okUbVpLvhQ1pm1X/interface-go-ipfs-core"
	ds "gx/ipfs/Qmf4xQhNomPNhrtZc67qSnfJSjxjXs9LWvknJtSXwimPrM/go-datastore"
	"time"
)

func main() {

	ctx := context.Background()
	lds := ds.NewMapDatastore()
	repo, _ := p2p.DefaultRepo(lds)
	ncfg := &core.BuildCfg{
		Repo: repo, //opt
		//Permanent:                 true, //opt, true|false
		//Repo:   mockRepo, //opt
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

	for {
		//ps := n.Peerstore.Peers()
		//for _, p := range ps {
		//	fmt.Println("---->", p.Pretty())
		//}
		//fmt.Println("count ====>", len(ps))
		//time.Sleep(time.Second)

		cns := n.PeerHost.Network().Conns()
		for _, cn := range cns {
			streams := cn.GetStreams()
			fmt.Println("stream:", streams)
		}
		fmt.Println("---->", len(cns), cns)
		time.Sleep(time.Second)
	}

	ctx1 := p2p.Ctx(n, "")
	api, err := ctx1.GetAPI()
	if err != nil {
		panic(err)
	}

	p, err := iface.ParsePath("QmVwyPQRKvoxvfxnVx2Y4BEitHFYDVtEqdqtspYGrVKAxF")
	if err != nil {
		panic(err)
	}

	fn, err := api.Unixfs().Get(ctx, p)
	if err != nil {
		panic(err)
	}
	fmt.Println(fn)
	file, ok := fn.(files.File)
	if !ok {
		panic("not file")
	}

	body := []byte{}
	size, err := file.Read(body)
	if err != nil {
		panic(err)
	}
	fmt.Println("file size:", size, string(body))
}
