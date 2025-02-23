package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
)

type Config struct {
	AddrServer    string
	PollTimeout   int64
	ReportTimeout int64
	SecretKey     string
	RateLimit     int64
}

func (cfg Config) isValid() bool {
	_, _, err := net.SplitHostPort(cfg.AddrServer)
	if err != nil {
		fmt.Println("Ошибка в передачи пармерта сервера")
		return false
	}
	if cfg.PollTimeout < 1 {
		fmt.Println("Частота отправки данных на сервер должна быть больше 0")
		return false
	}

	if cfg.ReportTimeout < 1 {
		fmt.Println("Частота снятия данных должна быть больше 0")
		return false
	}

	if cfg.RateLimit < 1 {
		fmt.Println("Количество одновременно исходящих запросов на сервер должна быть больше 0")
		return false
	}
	return true
}

func GetConfig() (*Config, bool) {
	cfg := &Config{}
	flag.StringVar(&cfg.AddrServer, "a", "localhost:8080", "server and port to run server")
	flag.Int64Var(&cfg.PollTimeout, "p", 2, "pool interval")
	flag.Int64Var(&cfg.ReportTimeout, "r", 10, "report interval")
	flag.StringVar(&cfg.SecretKey, "k", "", "secret key for encryption")
	flag.Int64Var(&cfg.RateLimit, "l", 2, "rate limit")

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		cfg.AddrServer = envRunAddr
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

	flag.Parse()
	return cfg, cfg.isValid()
}
