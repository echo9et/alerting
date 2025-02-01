package storage

import (
	"context"
	"database/sql"
	"fmt"
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
	fmt.Println(b.addr)
	db, err := sql.Open("pgx", b.addr)
	if err != nil {
		fmt.Println("dont open db")
		return false
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err = db.PingContext(ctx); err != nil {
		fmt.Println("timeout")
		return false
	}
	fmt.Println("database connect")

	return true
}
