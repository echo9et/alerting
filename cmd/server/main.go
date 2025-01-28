package main

import (
	"time"

	"github.com/echo9et/alerting/internal/entities"
	"github.com/echo9et/alerting/internal/logger"
	"github.com/echo9et/alerting/internal/server/coreserver"
	"github.com/echo9et/alerting/internal/server/storage"
)

func main() {

	ParseFlags()

	settingsServer := entities.Settings{
		Address:         *flagAddrServer,
		LogLevel:        *flagLogLevel,
		StoreInterval:   time.Second * time.Duration(*flagStoreInterval),
		FilenameStorage: *flagFilenameSave,
		IsRestore:       *flagRestoreData,
	}

	// storage := storage.NewMemStorage()
	storage, err := storage.NewSaver(storage.NewMemStore(), settingsServer.FilenameStorage, settingsServer.IsRestore, settingsServer.StoreInterval)
	if err != nil {
		panic(err)
	}
	logger.Initilization(settingsServer.LogLevel)

	if err := coreserver.Run(settingsServer, storage); err != nil {
		panic(err)
	}

}
