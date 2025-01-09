package main

import (
	"flag"
	"os"
)

var addrServer *string = flag.String("a", "localhost:8080", "server and port to run server")

func ParseFlags() {
	if envRunAddr := os.Getenv("ADDRESS"); envRunAddr != "" {
		*addrServer = envRunAddr
	}

	flag.Parse()
}
