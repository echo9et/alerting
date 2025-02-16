package main

import (
	"time"

	"github.com/echo9et/alerting/internal/agent/client"
)

func main() {
	config, status := GetConfig()
	if !status {
		panic("Не верно проинцелизирован конфиг файл")
	}
	a := client.NewAgent(config.AddrServer)
	r := time.Duration(config.ReportTimeout) * time.Second
	p := time.Duration(config.PollTimeout) * time.Second
	a.UpdateMetrics(r, p, config.SecretKey, config.RateLimit)
}
