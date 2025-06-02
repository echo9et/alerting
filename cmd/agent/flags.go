package main

import (
	"encoding/json"
	"flag"
	"net"
	"os"
	"strconv"

	"log/slog"
)

type Config struct {
	AddrServer    string `json:"address,omitempty"`
	SelfIP        string `json:"selfIP,omitempty"`
	PollTimeout   int64  `json:"poll_interval,omitempty"`
	ReportTimeout int64  `json:"report_interval,omitempty"`
	SecretKey     string `json:"key,omitempty"`
	RateLimit     int64  `json:"rate_limit,omitempty"`
	CryptoKey     string `json:"crypto_key,omitempty"`
	UseGRPC       bool   `json:"use_grpc,omitempty"`
}

func (cfg Config) isValid() bool {
	_, _, err := net.SplitHostPort(cfg.AddrServer)
	if err != nil {
		slog.Error("Ошибка в передаче параметра сервера")
		return false
	}
	if cfg.PollTimeout < 1 {
		slog.Error("Частота отправки данных на сервер должна быть больше 0")
		return false
	}

	if cfg.ReportTimeout < 1 {
		slog.Error("Частота снятия данных должна быть больше 0")
		return false
	}

	if cfg.RateLimit < 1 {
		slog.Error("Количество одновременно исходящих запросов должно быть больше 0")
		return false
	}
	return true
}

func GetConfig() (*Config, bool) {
	cfg := &Config{}

	// Флаги командной строки
	var configFilePath string
	flag.StringVar(&configFilePath, "c", "", "путь к конфигурационному файлу")
	flag.StringVar(&configFilePath, "config", "", "путь к конфигурационному файлу")

	flag.StringVar(&cfg.AddrServer, "a", "localhost:8080", "server and port to run server")
	flag.Int64Var(&cfg.PollTimeout, "p", 2, "pool interval")
	flag.Int64Var(&cfg.ReportTimeout, "r", 10, "report interval")
	flag.StringVar(&cfg.SecretKey, "k", "", "secret key for encryption")
	flag.Int64Var(&cfg.RateLimit, "l", 2, "rate limit")
	flag.StringVar(&cfg.CryptoKey, "crypto-key", "", "public key")
	flag.StringVar(&cfg.SelfIP, "self-ip", "127.0.0.1", "your ip address")
	flag.BoolVar(&cfg.UseGRPC, "g", false, "use grpc")

	// Читаем переменные окружения
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		cfg.AddrServer = envRunAddr
	}

	// Читаем переменные окружения
	if envSelfIP := os.Getenv("SELF_IP"); envSelfIP != "" {
		cfg.SelfIP = envSelfIP
	}

	if envReportInterval := os.Getenv("REPORT_INTERVAL"); envReportInterval != "" {
		value, ok := strconv.ParseInt(envReportInterval, 10, 0)
		if ok == nil {
			cfg.ReportTimeout = value
		}
	}

	if envPoolInterval := os.Getenv("POLL_INTERVAL"); envPoolInterval != "" {
		value, ok := strconv.ParseInt(envPoolInterval, 10, 0)
		if ok == nil {
			cfg.PollTimeout = value
		}
	}

	if envSecretKey := os.Getenv("KEY"); envSecretKey != "" {
		cfg.SecretKey = envSecretKey
	}

	if envRateLimit := os.Getenv("RATE_LIMIT"); envRateLimit != "" {
		value, ok := strconv.ParseInt(envRateLimit, 10, 0)
		if ok == nil {
			cfg.RateLimit = value
		}
	}

	if envCryptoKey := os.Getenv("CRYPTO_KEY"); envCryptoKey != "" {
		cfg.CryptoKey = envCryptoKey
	}

	if envUseGRPC := os.Getenv("USE_GRPC"); envUseGRPC != "" {
		cfg.UseGRPC = envUseGRPC == "true"
	}

	flag.Parse()

	if configFilePath != "" {
		fileData, err := os.ReadFile(configFilePath)
		if err != nil {
			slog.Error("Не удалось открыть конфигурационный файл", "error", err)
			return nil, false
		}
		tmpCfg := &Config{}
		err = json.Unmarshal(fileData, tmpCfg)
		if err != nil {
			slog.Error("Не удалось распарсить конфигурационный файл", "error", err)
			return nil, false
		}

		if cfg.AddrServer == "localhost:8080" && tmpCfg.AddrServer != "" {
			cfg.AddrServer = tmpCfg.AddrServer
		}
		if cfg.SelfIP == "127.0.0.1" && tmpCfg.AddrServer != "" {
			cfg.SelfIP = tmpCfg.SelfIP
		}
		if flag.Lookup("p").Value.String() == "2" {
			cfg.PollTimeout = tmpCfg.PollTimeout
		}
		if flag.Lookup("r").Value.String() == "10" {
			cfg.ReportTimeout = tmpCfg.ReportTimeout
		}
		if flag.Lookup("k").Value.String() == "" && tmpCfg.SecretKey != "" {
			cfg.SecretKey = tmpCfg.SecretKey
		}
		if flag.Lookup("l").Value.String() == "2" {
			cfg.RateLimit = tmpCfg.RateLimit
		}
		if flag.Lookup("crypto-key").Value.String() == "" && tmpCfg.CryptoKey != "" {
			cfg.CryptoKey = tmpCfg.CryptoKey
		}
	}

	return cfg, cfg.isValid()
}
