package storage

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
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

	conn, err := pgx.Connect(context.Background(), b.addr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		return false
	}
	defer conn.Close(context.Background())

	return true
}
