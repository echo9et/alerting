package storage

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Base struct {
	addr string
	conn *sql.DB
}

func NewPDatabase(a string) (*Base, error) {
	base := &Base{
		addr: a,
		conn: nil,
	}
	if err := base.Open(); err != nil {
		return nil, err
	}
	fmt.Println("--- dataBase is open")
	return base, nil
}

func (b *Base) Open() error {
	// host=localhost user=echo9et password=123321 dbname=echo9et sslmode=disable

	bd, err := sql.Open("pgx", b.addr)
	if err != nil {
		fmt.Println("---", err)
		return err
	}

	b.conn = bd

	if err := b.InitTable(); err != nil {
		fmt.Println("---", err)
		return err
	}
	return nil

}

func (b *Base) InitTable() error {
	if b.conn == nil {
		return fmt.Errorf("b.conn is nil")
	}
	_, err := b.conn.QueryContext(context.Background(),
		`CREATE TABLE IF NOT EXISTS metrics_gauge (id serial PRIMARY KEY, name varchar(255) UNIQUE NOT NULL, value DOUBLE PRECISION NOT NULL);`)
	// err := row.Scan()
	if err != nil {
		print("metric gauge")
		return err
	}
	_, err = b.conn.QueryContext(context.Background(),
		`CREATE TABLE IF NOT EXISTS metrics_counter (id serial PRIMARY KEY, name varchar(255) UNIQUE NOT NULL, value INTEGER NOT NULL);`)
	// err = row.Scan()

	if err != nil {
		return err
	}

	fmt.Println("Create table ok")

	return nil
}

func (b *Base) Ping() bool {

	defer b.conn.Close()
	if b.conn == nil {
		fmt.Println("---", "Ping nil b.conn")
		return false
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	if err := b.conn.PingContext(ctx); err != nil {
		fmt.Println("--- Context False", err)
		return false
	}
	return true
}

func (b *Base) GetCounter(string) (string, bool) {
	return "", false
}

func (b *Base) SetCounter(string, int64) {
}

func (b *Base) GetGauge(string) (string, bool) {
	return "", false
}

func (b *Base) SetGauge(string, float64) {
}

func (b *Base) AllMetrics() map[string]string {
	return map[string]string{}
}
