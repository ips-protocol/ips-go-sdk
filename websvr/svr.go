package websvr

import (
	"github.com/ipweb-group/go-sdk/conf"
	"github.com/kataras/iris"
)

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

	err = app.Run(iris.Addr(cfg.ServerHost))
	if err != nil {
		panic(err)
	}

}
