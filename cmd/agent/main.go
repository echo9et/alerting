package main

import (
	"time"

	"github.com/echo9et/alerting/cmd/agent/client"
)

func main() {
	initAgent()
	parseFlags()

	a := client.NewAgent(*addrServer)
	reportTime := time.Duration(*reportTimeout) * time.Second
	poolTime := time.Duration(*pollTimeout) * time.Second
	a.UpdateMetrics(reportTime, poolTime)
}
