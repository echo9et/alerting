package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"

	"log/slog"
)

type Config struct {
	AddrServer    string `json:"address,omitempty"`
	AddrDatabase  string `json:"database_dsn,omitempty"`
	LogLevel      string `json:"log_level,omitempty"`
	StoreInterval uint64 `json:"store_interval,omitempty"`
	FilenameSave  string `json:"file_storage_path,omitempty"`
	RestoreData   bool   `json:"restore,omitempty"`
	SecretKey     string `json:"key,omitempty"`
	CryptoKey     string `json:"crypto_key,omitempty"`
}

func ParseFlags() (*Config, error) {
	cfg := &Config{}

	var configFilePath string
	flag.StringVar(&configFilePath, "c", "", "путь к конфигурационному файлу")
	flag.StringVar(&configFilePath, "config", "", "путь к конфигурационному файлу")

	// Флаги
	flag.StringVar(&cfg.AddrServer, "a", "localhost:8080", "server and port to run server")
	flag.StringVar(&cfg.AddrDatabase, "d", "", "address to postgres base")
	flag.StringVar(&cfg.LogLevel, "l", "info", "log level")
	flag.Uint64Var(&cfg.StoreInterval, "i", 300, "save to file interval")
	flag.StringVar(&cfg.FilenameSave, "f", "data.json", "filename for save and restore data")
	flag.BoolVar(&cfg.RestoreData, "r", true, "is restor data from file")
	flag.StringVar(&cfg.SecretKey, "k", "", "secret key for encryption")
	flag.StringVar(&cfg.CryptoKey, "crypto-key", "", "privat key")

	// Переменные окружения
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
		uValue, err := strconv.ParseUint(envStoreInterval, 10, 64)
		if err == nil {
			cfg.StoreInterval = uValue
		}
	}

	if envFileNameSave := os.Getenv("FILE_STORAGE_PATH"); envFileNameSave != "" {
		cfg.FilenameSave = envFileNameSave
	}

	if envRestoreData := os.Getenv("RESTORE"); envRestoreData != "" {
		cfg.RestoreData = envRestoreData == "true"
	}

	if envSecretKey := os.Getenv("KEY"); envSecretKey != "" {
		cfg.SecretKey = envSecretKey
	}

	if envCryptoKey := os.Getenv("CRYPTO-KEY"); envCryptoKey != "" {
		cfg.CryptoKey = envCryptoKey
	}

	flag.Parse()

	// Чтение JSON-конфига (если указан)
	if configFilePath == "" {
		configFilePath = os.Getenv("CONFIG")
	}

	if configFilePath != "" {
		fileData, err := os.ReadFile(configFilePath)
		if err != nil {
			slog.Error("Не удалось открыть конфигурационный файл", "error", err)
			return nil, err
		}

		tmpCfg := &Config{}
		err = json.Unmarshal(fileData, tmpCfg)
		if err != nil {
			slog.Error("Ошибка при разборе JSON конфига", "error", err)
			return nil, err
		}

		// JSON
		if flag.Lookup("a").Value.String() == "localhost:8080" && tmpCfg.AddrServer != "" {
			cfg.AddrServer = tmpCfg.AddrServer
		}
		if flag.Lookup("d").Value.String() == "" && tmpCfg.AddrDatabase != "" {
			cfg.AddrDatabase = tmpCfg.AddrDatabase
		}
		if flag.Lookup("l").Value.String() == "info" && tmpCfg.LogLevel != "" {
			cfg.LogLevel = tmpCfg.LogLevel
		}
		if flag.Lookup("i").Value.String() == "300" && tmpCfg.StoreInterval > 0 {
			cfg.StoreInterval = tmpCfg.StoreInterval
		}
		if flag.Lookup("f").Value.String() == "data.json" && tmpCfg.FilenameSave != "" {
			cfg.FilenameSave = tmpCfg.FilenameSave
		}
		if flag.Lookup("r").Value.String() == "true" {
			cfg.RestoreData = tmpCfg.RestoreData
		}
		if flag.Lookup("k").Value.String() == "" && tmpCfg.SecretKey != "" {
			cfg.SecretKey = tmpCfg.SecretKey
		}
		if flag.Lookup("crypto-key").Value.String() == "" && tmpCfg.CryptoKey != "" {
			cfg.CryptoKey = tmpCfg.CryptoKey
		}
	}

	// Валидация
	_, _, err := net.SplitHostPort(cfg.AddrServer)
	if err != nil {
		return nil, fmt.Errorf("неверный адрес сервера: %w", err)
	}

	return cfg, nil
}
