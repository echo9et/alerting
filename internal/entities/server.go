package entities

import "time"

type Settings struct {
	Address         string
	LogLevel        string
	StoreInterval   time.Duration
	FilenameStorage string
	IsRestore       bool
}
