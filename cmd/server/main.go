package main

import (
	"github.com/echo9et/alerting/cmd/server/coreserver"
)

func main() {
	if err := coreserver.Run(); err != nil {
		panic(err)
	}
}
