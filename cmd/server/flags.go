package main

import (
	"flag"
	"net"
	"os"
	"strconv"
)

type Config struct {
	AddrServer    string
	AddrDatabase  string
	LogLevel      string
	StoreInterval uint64
	FilenameSave  string
	RestoreData   bool
	SecretKey     string
}

func ParseFlags() (*Config, error) {
	cfg := &Config{}
	flag.StringVar(&cfg.AddrServer, "a", "localhost:8080", "server and port to run server")
	flag.StringVar(&cfg.AddrDatabase, "d", "", "address to postgres base")
	// flag.StringVar(&cfg.AddrDatabase, "d", "host=localhost user=username password=123321 dbname=echo9et sslmode=disable", "address to postgres base")
	flag.StringVar(&cfg.LogLevel, "l", "info", "log level")
	flag.Uint64Var(&cfg.StoreInterval, "i", 300, "save to file interval")
	flag.StringVar(&cfg.FilenameSave, "f", "data.json", "filename for save and restore data")
	flag.BoolVar(&cfg.RestoreData, "r", true, "is restor data from file")
	flag.StringVar(&cfg.SecretKey, "k", "", "secret key for encryption")

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		cfg.AddrServer = envRunAddr
	}

	if envDatabaseAddr := os.Getenv("DATABASE_DSN"); envDatabaseAddr != "" {
		cfg.AddrDatabase = envDatabaseAddr
	}

	if envRunLogLVL := os.Getenv("LOG_LVL"); envRunLogLVL != "" {
		cfg.LogLevel = envRunLogLVL
	}

	if envStoreInterval := os.Getenv("STORE_INTERVAL"); envStoreInterval != "" {
		uValue, err := strconv.ParseUint(string(envStoreInterval), 10, 64)
		if err != nil {
			return nil, err
		}
		cfg.StoreInterval = uValue
	}

	if envFileNameSave := os.Getenv("FILE_STORAGE_PATH"); envFileNameSave != "" {
		cfg.FilenameSave = envFileNameSave
	}

	if envRestoreData := os.Getenv("RESTORE"); envRestoreData != "" {
		cfg.RestoreData = envRestoreData == "true"
	}

	if envSecretKey := os.Getenv("KEY"); envSecretKey != "" {
		cfg.AddrServer = envSecretKey
	}

	flag.Parse()

	_, _, err := net.SplitHostPort(cfg.AddrServer)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
