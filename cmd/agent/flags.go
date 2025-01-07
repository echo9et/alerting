package main

import (
	"flag"
)

var addrServer *string = flag.String("a", "localhost:8080", "address and port to run server")
var pollTimeout *int = flag.Int("p", 2, "pool interval")
var reportTimeout *int = flag.Int("r", 10, "report interval")

func parseFlags() {
	flag.Parse()
}
