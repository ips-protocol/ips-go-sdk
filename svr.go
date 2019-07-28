package main

import (
	"github.com/ipweb-group/go-sdk/conf"
	"github.com/ipweb-group/go-sdk/putPolicy"
	"github.com/ipweb-group/go-sdk/websvr"
)

func main() {
	putPolicy.LoadAppClients("./websvr/app-clients.json")
	conf.LoadConfig("./websvr/conf.json")

	websvr.Run()
}
