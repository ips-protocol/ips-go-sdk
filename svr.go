package main

import (
	"./putPolicy"
	"./websvr"
)

func main() {
	putPolicy.LoadAppClients("./websvr/app-clients.json")

	websvr.Run("./websvr/conf.json")
}
