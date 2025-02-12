package main

import (
	"flag"
	"os"
	"strconv"
)

var addrServer *string = flag.String("a", "localhost:8080", "address and port to run server")
var pollTimeout *int = flag.Int("p", 2, "pool interval")
var reportTimeout *int = flag.Int("r", 10, "report interval")

func parseFlags() {
	flag.Parse()
}

func initAgent() {
	parseFlags()

	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		*addrServer = envRunAddr
	}

	if envReportInterval := os.Getenv("REPORT_INTERVAL"); envReportInterval != "" {
		value, ok := strconv.Atoi(envReportInterval)
		if ok == nil {
			*reportTimeout = value
		}
	}

	if envPoolInterval := os.Getenv("POLL_INTERVAL"); envPoolInterval != "" {
		value, ok := strconv.Atoi(envPoolInterval)
		if ok == nil {
			*pollTimeout = value
		}
	}

}
