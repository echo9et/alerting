package main

import (
	"time"
	"github.com/echo9et/alerting/internal/agent/client"
)

func main() {
	// time.Sleep(3 * time.Second)
	config, err := GetConfig()
	if err != nil {
		panic(err)
	}
	a := client.NewAgent(config.AddrServer)
	r := time.Duration(config.ReportTimeout) * time.Second
	p := time.Duration(config.PollTimeout) * time.Second
	a.UpdateMetrics(r, p)
}
