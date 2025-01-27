package main

import (
	"flag"
	"os"
	"strconv"
)

var flagAddrServer *string = flag.String("a", "localhost:8080", "server and port to run server")
var flagLogLevel *string = flag.String("l", "info", "log level")

var flagStoreInterval *uint64 = flag.Uint64("i", 300, "save to file interval")
var flagFilenameSave *string = flag.String("f", "data.json", "filename for save and restore data")
var flagRestoreData *bool = flag.Bool("r", true, "is restor data from file")

func ParseFlags() {
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		*flagAddrServer = envRunAddr
	}

	if envRunLogLVL := os.Getenv("LOG_LVL"); envRunLogLVL != "" {
		*flagLogLevel = envRunLogLVL
	}

	if envStoreInterval := os.Getenv("STORE_INTERVAL"); envStoreInterval != "" {
		uValue, err := strconv.ParseUint(string(envStoreInterval), 10, 64)
		if err != nil {
			panic(err)
		}
		*flagStoreInterval = uValue
	}

	if envFileNameSave := os.Getenv("FILE_STORAGE_PATH"); envFileNameSave != "" {
		*flagFilenameSave = envFileNameSave
	}

	if envRestoreData := os.Getenv("RESTORE"); envRestoreData != "" {
		*flagRestoreData = envRestoreData == "true"
	}

	flag.Parse()
}
