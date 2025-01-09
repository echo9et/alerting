package main

import (
	"github.com/echo9et/alerting/cmd/server/coreserver"
	"github.com/echo9et/alerting/cmd/server/storage"
)

func main() {
	ParseFlags()
	storage := storage.NewMemStorage()
	if err := coreserver.Run(*addrServer, storage); err != nil {
		panic(err)
	}
}
