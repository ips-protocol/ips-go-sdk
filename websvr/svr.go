package websvr

import (
	"github.com/ipweb-group/go-sdk/conf"
	"github.com/kataras/iris"
)

type Config struct {
	ServerWriteTimeoutInSecond int         `json:"server_write_timeout_in_second"`
	ServerReadTimeoutInSecond  int         `json:"server_read_timeout_in_second"`
	ServerHost                 string      `json:"server_host"`
	NodeConf                   conf.Config `json:"node_conf"`
}

func Run(cfgPath string) {
	cfg := Config{}
	err := conf.LoadConf(&cfg, cfgPath)
	if err != nil {
		panic(err)
	}

	service, err := NewService(cfg.NodeConf)
	if err != nil {
		panic(err)
	}

	app := iris.Default()
	app.Post("/file/upload", service.FileUpload)
	app.Get("/file/{cid: string}", service.FileDownload)
	app.Get("/file/stream/{cid: string}", service.FileStreamRead)
	app.Delete("/file/{cid: string}", service.FileDelete)
	app.Run(iris.Addr(cfg.ServerHost))
}
