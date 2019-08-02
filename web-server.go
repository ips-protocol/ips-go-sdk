package main

import (
	"context"
	"fmt"
	"github.com/ipweb-group/go-sdk/conf"
	"github.com/ipweb-group/go-sdk/putPolicy"
	"github.com/ipweb-group/go-sdk/putPolicy/persistent"
	"github.com/ipweb-group/go-sdk/rpc"
	"github.com/ipweb-group/go-sdk/utils/redis"
	"github.com/ipweb-group/go-sdk/websvr/controllers"
	"github.com/kataras/iris"
	irisContext "github.com/kataras/iris/context"
	"time"
)

// 最大允许上传的文件大小：500MB
const MaxFileSize int64 = 500 << 20

func main() {
	putPolicy.LoadAppClients("./websvr/app-clients.json")
	conf.LoadConfig("./websvr/conf.json")

	cfg := conf.GetConfig()

	// 初始化 RPC 客户端
	rpcClient, err := rpc.GetClientInstance()
	if err != nil {
		panic(err)
	}

	// 启动转换器线程
	go persistent.ConvertMediaJob()

	// 初始化 Web 服务器
	app := iris.Default()

	// 404 错误输出
	app.OnErrorCode(iris.StatusNotFound, func(ctx irisContext.Context) {
		_, _ = ctx.JSON(iris.Map{"error": "Document not found"})
	})

	// 构建路由
	routers(app)

	// 平滑关闭服务
	iris.RegisterOnInterrupt(func() {
		timeout := 15 * time.Second
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()

		// 关闭 RPC 客户端
		err := rpcClient.Close()
		if err != nil {
			fmt.Printf("[WARN] Close RPC Client failed, %v \n", err)
		}

		// 关闭 Redis 连接
		err = redis.GetClient().Close()
		if err != nil {
			fmt.Printf("[WARN] Close redis connection failed, %v \n", err)
		}

		_ = app.Shutdown(ctx)
	})

	_ = app.Run(iris.Addr(cfg.ServerHost), iris.WithoutInterruptHandler)

	// app.Run(iris.AutoTLS(":443", "example.com", "admin@example.com")) 可以自动配置 Lets Encrypt 证书
}

// 构建路由
func routers(app *iris.Application) {
	// Version 1
	v1 := app.Party("/v1")
	{
		uploadController := controllers.UploadController{}
		v1.Post("/upload", iris.LimitRequestBodySize(MaxFileSize), uploadController.Upload)

		downloadController := controllers.DownloadController{}
		v1.Get("/file/{cid:string}", downloadController.StreamedDownload)
	}
}
