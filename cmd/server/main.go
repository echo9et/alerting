package main

import (
	"github.com/echo9et/alerting/cmd/server/server"
)

func main() {
	if err := server.Run(); err != nil {
		panic(err)
	}
}
