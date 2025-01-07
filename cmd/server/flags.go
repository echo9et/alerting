package main

import (
	"flag"
)

var addrServer *string = flag.String("a", "localhost:8080", "server and port to run server")

func PaeseFlags() {
	flag.Parse()
}
