package websvr

import (
	"github.com/ipweb-group/go-sdk/conf"
	"github.com/ipweb-group/go-sdk/websvr/uploadController"
	"github.com/kataras/iris"
)

// 最大允许上传的文件大小：500MB
const MaxFileSize int64 = 500 << 20

func Run() {
	cfg := conf.GetConfig()

	service, err := NewService(cfg.NodeConf)
	if err != nil {
		panic(err)
	}

	app := iris.Default()
	app.Post("/file/upload", service.FileUpload)
	app.Get("/file/{cid: string}", service.FileDownload)
	app.Get("/file/stream/{cid: string}", service.FileStreamRead)
	app.Delete("/file/{cid: string}", service.FileDelete)
	app.Get("/nodes", service.NodesList)

	/**
	 * Version 1
	 */
	v1 := app.Party("/v1")
	{
		controller := uploadController.New()
		v1.Post("/upload", iris.LimitRequestBodySize(MaxFileSize), controller.Upload)
	}

	err = app.Run(iris.Addr(cfg.ServerHost))
	if err != nil {
		panic(err)
	}

}
