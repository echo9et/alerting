package main

import (
	"fmt"
	"time"

	"github.com/echo9et/alerting/internal/agent/client"
	"github.com/echo9et/alerting/internal/entities"
)

// Используй флаги сборки
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

	config, status := GetConfig()
	if !status {
		panic("Не верно проинцелизирован конфиг файл")
	}

	a := client.NewAgent(config.AddrServer, config.SelfIP, config.UseGRPC)
	r := time.Duration(config.ReportTimeout) * time.Second
	p := time.Duration(config.PollTimeout) * time.Second

	if config.CryptoKey != "" {
		pub, err := entities.GetPubKey(config.CryptoKey)
		if err != nil {
			panic(err)
		}
		a.UpdateMetrics(r, p, config.SecretKey, config.RateLimit, pub)
	}

	a.UpdateMetrics(r, p, config.SecretKey, config.RateLimit, nil)
}
