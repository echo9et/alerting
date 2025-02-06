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
	rows, err := b.conn.QueryContext(context.Background(),
		`CREATE TABLE IF NOT EXISTS metrics_gauge (name varchar(255) PRIMARY KEY UNIQUE NOT NULL, value DOUBLE PRECISION NOT NULL);`)
	if err != nil {
		return err
	}
	if rows.Err() != nil {
		return err
	}
	rows.Close()

	rows, err = b.conn.QueryContext(context.Background(),
		`CREATE TABLE IF NOT EXISTS metrics_counter (name varchar(255) PRIMARY KEY UNIQUE NOT NULL, value INTEGER NOT NULL);`)

	if err != nil {
		return err
	}

	if rows.Err() != nil {
		return err
	}
	rows.Close()

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

func (b *Base) GetCounter(name string) (string, bool) {
	query := `SELECT value FROM metrics_counter WHERE name=$1;`
	var desc sql.NullString
	err := b.conn.QueryRow(query, name).Scan(&desc)

	if err != nil {
		fmt.Println("ERROR GetCounter", err)
		return "", false
	}

	return desc.String, true
}

func (b *Base) SetCounter(name string, iValue int64) {
	_, err := b.conn.Exec(
		`INSERT INTO metrics_counter (name, value) 
		VALUES ($1, $2) ON CONFLICT (name) 
		DO UPDATE SET value = EXCLUDED.value;`, name, iValue)
	if err != nil {
		fmt.Println("---", err)
	}
}

func (b *Base) GetGauge(name string) (string, bool) {
	query := `SELECT value FROM metrics_gauge WHERE name=$1;`
	var desc sql.NullString
	err := b.conn.QueryRow(query, name).Scan(&desc)

	if err != nil {
		fmt.Println("ERROR GetCounter", err)
		return "", false
	}

	return desc.String, true
}

func (b *Base) SetGauge(name string, fValue float64) {
	_, err := b.conn.Exec(
		`INSERT INTO metrics_gauge (name, value) 
		VALUES ($1, $2) ON CONFLICT (name) 
		DO UPDATE SET value = EXCLUDED.value;`, name, fValue)
	if err != nil {
		fmt.Println("---", err)
	}
}

func (b *Base) AllMetrics() map[string]string {
	out := make(map[string]string)

	query := `SELECT * FROM metrics_gauge
			  UNION ALL
		      SELECT * FROM metrics_counter;`
	rows, err := b.conn.Query(query)
	if err != nil {
		fmt.Println("error request AllMetrics")
		return out
	}
	defer rows.Close()

	var name, value string
	for rows.Next() {
		err = rows.Scan(&name, &value)
		if err != nil {
			fmt.Println("error AllMetrics read data ")
			return out
		}
		out[name] = value
	}
	err = rows.Err()
	if err != nil {
		fmt.Println("error AllMetrics rows data ")
	}
	return out
}
