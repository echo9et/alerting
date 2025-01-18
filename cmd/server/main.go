package main

import (
	"github.com/echo9et/alerting/internal/logger"
	"github.com/echo9et/alerting/internal/server/coreserver"
	"github.com/echo9et/alerting/internal/server/storage"
)

func main() {
	ParseFlags()
	storage := storage.NewMemStorage()
	logger.Initilization(*flagLogLevel)
	if err := coreserver.Run(*addrServer, storage); err != nil {
		panic(err)
	}
}
