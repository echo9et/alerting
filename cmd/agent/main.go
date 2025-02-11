package main

import (
	"time"

	"github.com/echo9et/alerting/internal/agent/client"
)

func main() {
	// time.Sleep(3 * time.Second)
	initAgent()
	parseFlags()

	a := client.NewAgent(*addrServer)
	reportTime := time.Duration(*reportTimeout) * time.Second
	poolTime := time.Duration(*pollTimeout) * time.Second
	a.UpdateMetrics(reportTime, poolTime)
}
