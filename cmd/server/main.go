package main

import (
	"fmt"
	"time"

	"github.com/echo9et/alerting/internal/entities"
	"github.com/echo9et/alerting/internal/logger"
	"github.com/echo9et/alerting/internal/server/coreserver"
	"github.com/echo9et/alerting/internal/server/storage"

	"log/slog"
)

// Испоользуй флаги сборки
// go build -ldflags "-X main.buildVersion=1.0.0"
var (
	buildVersion string = "N/A"
	buildDate    string = "N/A"
	buildCommit  string = "N/A"
)

func main() {
	fmt.Printf("Build version: %s\n", buildVersion)
	fmt.Printf("Build date: %s\n", buildDate)
	fmt.Printf("Build commit: %s\n", buildCommit)

	cfg, err := ParseFlags()
	if err != nil {
		panic(err)
	}

	var store entities.Storage
	if cfg.AddrDatabase != "" {
		slog.Info("start with postgres")
		store, err = storage.NewPDatabase(cfg.AddrDatabase)
		if err != nil {
			panic(err)
		}
	} else {
		slog.Info("start with mem storage")
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
