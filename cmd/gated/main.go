package main

import (
	"github.com/gopherd/doge/cmd/gated/server"
	"github.com/gopherd/doge/service"
)

func main() {
	service.Run(server.New())
}
