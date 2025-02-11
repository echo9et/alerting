package main

import (
	"errors"
	"flag"
	"os"
	"strconv"
	"net"
)

type Config struct {
	AddrServer    string
	PollTimeout   int64
	ReportTimeout int64
}


func (cfg Config)isValid() error{
	_, _, err := net.SplitHostPort(cfg.AddrServer)
	if err != nil {
		return errors.New("Ошибка в передачи пармерта сервера")
	}
	if cfg.PollTimeout < 1 {
		return errors.New("Частота отправки данных на сервер должна быть больше 0")
	}

	if cfg.ReportTimeout < 1 {
		return errors.New("Частота снятия данных должна быть больше 0")
	}
	return nil
}

func GetConfig() (*Config, error){
	cfg := &Config{}
	flag.StringVar(&cfg.AddrServer, "a", "localhost:8080", "server and port to run server")
	flag.Int64Var(&cfg.PollTimeout, "p", 2, "pool interval")
	flag.Int64Var(&cfg.ReportTimeout, "r", 10, "report interval")


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
	flag.Parse()
	return cfg, cfg.isValid()
}
