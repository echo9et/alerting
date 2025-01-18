package main

import (
	"flag"
	"os"
)

var addrServer *string = flag.String("a", "localhost:8080", "server and port to run server")
var flagLogLevel *string = flag.String("l", "info", "log level")

func ParseFlags() {
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		*addrServer = envRunAddr
	}

	if envRunLogLVL := os.Getenv("LOG_LVL"); envRunLogLVL != "" {
		*flagLogLevel = envRunLogLVL
	}

	flag.Parse()
}
