package main

import (
	"github.com/ipweb-group/go-sdk/conf"
	"github.com/ipweb-group/go-sdk/putPolicy/persistent"
)

func main() {
	conf.LoadConfig("./websvr/conf.json")

	persistent.ConvertMediaJob()
}
