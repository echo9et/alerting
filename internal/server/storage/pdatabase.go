package storage

import (
	"context"
	"database/sql"
	"time"
)

type Base struct {
	addr string
}

func NewPDatabase(a string) *Base {
	return &Base{
		addr: a,
	}
}

func (b *Base) Ping() bool {

	db, err := sql.Open("pgx", b.addr)
	if err != nil {
		return false
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err = db.PingContext(ctx); err != nil {
		return false
	}
	return true
}
