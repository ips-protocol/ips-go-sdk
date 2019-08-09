package websvr

import (
	"github.com/ipweb-group/go-sdk/conf"
	"github.com/kataras/iris"
)

func Run() {
	nodeConf := conf.Config{
		NodeRefreshIntervalInSecond: 1,
		NodeRefreshWorkers:          0,
		NodeRequestTimeoutInSecond:  300,
		NodeCloseIntervalInSecond:   0,
		ConnQuotaPerNode:            0,
		BlockUploadWorkers:          8,
		BlockDownloadWorkers:        0,
		ContractConf: conf.ContractConfig{
			ClientKeyHex:     "B2FE66D78810869A64CAAE7B1F2C60CCA3AC2F2261DA2F1DE7040DE3F1FEDA9C",
			ContractNodeAddr: "https://mainnet.ipweb.top",
		},
		ECConfig: conf.ECConfig{},
	}

	service, err := NewService(nodeConf)
	if err != nil {
		panic(err)
	}

	app := iris.Default()
	app.Post("/file/upload", service.FileUpload)
	app.Get("/file/{cid: string}", service.FileDownload)
	app.Get("/file/stream/{cid: string}", service.FileStreamRead)
	app.Delete("/file/{cid: string}", service.FileDelete)
	app.Get("/nodes", service.NodesList)

	err = app.Run(iris.Addr(":9090"))
	if err != nil {
		panic(err)
	}

}
