package main

import (
	"github.com/ipweb-group/go-sdk/conf"
	"github.com/ipweb-group/go-sdk/putPolicy"
	"github.com/ipweb-group/go-sdk/putPolicy/persistent"
	"github.com/ipweb-group/go-sdk/rpc"
	"github.com/ipweb-group/go-sdk/websvr/uploadController"
	"github.com/kataras/iris"
)

// 最大允许上传的文件大小：500MB
const MaxFileSize int64 = 500 << 20

func main() {
	putPolicy.LoadAppClients("./websvr/app-clients.json")
	conf.LoadConfig("./websvr/conf.json")

	cfg := conf.GetConfig()

	// 初始化 RPC 客户端
	// FIXME 将 rpcClient 改为单例，并通过工厂方法获取
	rpcClient, err := rpc.NewClient(cfg.NodeConf)
	if err != nil {
		panic(err)
	}

	// 启动转换器线程
	go persistent.ConvertMediaJob()

	// 初始化 Web 服务器
	app := iris.Default()

	// 构建路由
	routers(app, rpcClient)

	err = app.Run(iris.Addr(cfg.ServerHost))
	if err != nil {
		panic(err)
	}
}

// 构建路由
func routers(app *iris.Application, rpcClient *rpc.Client) {
	// Version 1
	v1 := app.Party("/v1")
	{
		controller := uploadController.UploadController{Node: rpcClient}
		v1.Post("/upload", iris.LimitRequestBodySize(MaxFileSize), controller.Upload)
	}
}
