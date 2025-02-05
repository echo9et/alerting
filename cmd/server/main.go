package main

import (
	"time"

	"github.com/echo9et/alerting/internal/logger"
	"github.com/echo9et/alerting/internal/server/coreserver"
	"github.com/echo9et/alerting/internal/server/handlers"
	"github.com/echo9et/alerting/internal/server/storage"
)

func main() {

	cfg, err := ParseFlags()
	if err != nil {
		panic(err)
	}

	var store handlers.Storage
	if cfg.AddrDatabase != "" {
		store, err = storage.NewPDatabase(cfg.AddrDatabase)
		if err != nil {
			panic(err)
		}
	} else {
		store, err = storage.NewSaver(storage.NewMemStore(), cfg.FilenameSave, cfg.RestoreData, time.Duration(cfg.StoreInterval)*time.Second)
		if err != nil {
			panic(err)
		}
	}

	logger.Initilization(cfg.LogLevel)

	if err := coreserver.Run(cfg.AddrServer, cfg.AddrDatabase, store); err != nil {
		panic(err)
	}

}
