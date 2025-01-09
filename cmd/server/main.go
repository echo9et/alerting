package main

import (
	"flag"
	"os"

	"github.com/echo9et/alerting/cmd/server/coreserver"
	"github.com/echo9et/alerting/cmd/server/storage"
)

var addrServer *string = flag.String("a", "localhost:8080", "server and port to run server")

func ParseFlags() {
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		*addrServer = envRunAddr
	}

	flag.Parse()
}

func main() {
	ParseFlags()
	storage := storage.NewMemStorage()
	if err := coreserver.Run(*addrServer, storage); err != nil {
		panic(err)
	}
}
