package main

import (
	"github.com/ipweb-group/go-sdk/conf"
	"github.com/ipweb-group/go-sdk/websvr"
)

func main() {
	conf.LoadConfig("./websvr/conf.json")

	websvr.Run()
}
