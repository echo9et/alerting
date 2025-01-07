package main

import (
	"github.com/echo9et/alerting/cmd/server/coreserver"
)

func main() {
	PaeseFlags()
	if err := coreserver.Run(*addrServer); err != nil {
		panic(err)
	}
}
