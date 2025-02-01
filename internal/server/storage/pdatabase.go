package storage

import (
	"context"
	"database/sql"
	"fmt"
	"net"
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
	ip, port, err := net.SplitHostPort(b.addr)
	if err != nil {
		fmt.Println("error parse addr")
		return false
	}
	cmd := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		ip, port, `postgres`, `postgres`, `postgres`)
	db, err := sql.Open("pgx", cmd)
	if err != nil {
		fmt.Println("dont open db", cmd)
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
