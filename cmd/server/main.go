package main

import (
	"fmt"
	"time"

	"github.com/echo9et/alerting/internal/entities"
	"github.com/echo9et/alerting/internal/logger"
	"github.com/echo9et/alerting/internal/server/coreserver"
	"github.com/echo9et/alerting/internal/server/storage"
)

func main() {

	cfg, err := ParseFlags()
	if err != nil {
		panic(err)
	}

	var store entities.Storage
	if cfg.AddrDatabase != "" {
		fmt.Println("start with postgres")
		store, err = storage.NewPDatabase(cfg.AddrDatabase)
		if err != nil {
			panic(err)
		}
	} else {
		fmt.Println("start with mem storage")
		store, err = storage.NewSaver(storage.NewMemStore(), cfg.FilenameSave, cfg.RestoreData, time.Duration(cfg.StoreInterval)*time.Second)
		if err != nil {
			panic(err)
		}
	}

	logger.Initilization(cfg.LogLevel)

	if err := coreserver.Run(cfg.AddrServer, cfg.AddrDatabase, store, cfg.SecretKey); err != nil {
		panic(err)
	}

}
