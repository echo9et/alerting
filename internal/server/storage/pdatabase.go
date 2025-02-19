package storage

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/echo9et/alerting/internal/entities"
	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
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
		`CREATE TABLE IF NOT EXISTS metrics_counter (name varchar(255) PRIMARY KEY UNIQUE NOT NULL, value bigint NOT NULL);`)

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
		fmt.Println("ERROR GetCounter ", name, err)
		return "", false
	}

	return desc.String, true
}

func (b *Base) SetCounter(name string, iValue int64) {
	_, err := b.conn.Exec(
		`INSERT INTO metrics_counter (name, value) 
		VALUES ($1, $2) 
		ON CONFLICT (name) 
		DO UPDATE SET value = metrics_counter.value + EXCLUDED.value;`, name, iValue)
	if err != nil {
		fmt.Println("---", err)
	}
}

func (b *Base) GetGauge(name string) (string, bool) {
	query := `SELECT value FROM metrics_gauge WHERE name=$1;`
	var desc sql.NullString
	err := b.conn.QueryRow(query, name).Scan(&desc)

	if err != nil {
		fmt.Println("ERROR GetGauge ", name, err)
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

func (b *Base) AllMetricsJSON() []entities.MetricsJSON {
	out := make([]entities.MetricsJSON, 0)
	return out
}

func (b *Base) SetMetrics(mertics []entities.MetricsJSON) error {
	var err error
	for _, dealy := range []time.Duration{1 * time.Second, 2 * time.Second, 5 * time.Second} {
		if err = b.requestSaveMerics(mertics); err != nil {
			if !isRunReplay(err) {
				break
			}
			time.Sleep(dealy)
		} else {
			return nil
		}
	}
	return err
}

func (b *Base) requestSaveMerics(mertics []entities.MetricsJSON) error {
	tx, err := b.conn.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmtGauge, err := tx.Prepare(
		`INSERT INTO metrics_gauge (name, value) 
		VALUES ($1, $2) 
		ON CONFLICT (name) 
		DO UPDATE SET value = EXCLUDED.value;`)
	if err != nil {
		return err
	}
	defer stmtGauge.Close()

	stmtCounter, err := tx.Prepare(
		`INSERT INTO metrics_counter (name, value) 
		VALUES ($1, $2) 
		ON CONFLICT (name) 
		DO UPDATE SET value = metrics_counter.value + EXCLUDED.value;`)
	if err != nil {
		return err
	}
	defer stmtCounter.Close()

	for _, v := range mertics {
		if v.MType == entities.Gauge {
			_, err := stmtGauge.Exec(v.ID, v.Value)
			if err != nil {
				return err
			}
		} else if v.MType == entities.Counter {
			_, err := stmtCounter.Exec(v.ID, v.Delta)
			fmt.Println("updates/ Counter", v.ID, *v.Delta)
			if err != nil {
				return err
			}
		} else {
			fmt.Println("Неизвестный тип метрики ", err)
			return errors.New("Неизвестный тип метрики " + v.ID)
		}
	}
	fmt.Println("All right commit ")
	return tx.Commit()
}

func isRunReplay(err error) bool {
	fmt.Printf("ошибка при обработке запроса к postgres: %v\n", err)
	if pgErr, ok := err.(*pgconn.PgError); ok {
		switch pgErr.Code {
		case pgerrcode.SerializationFailure:
			return true
		case pgerrcode.LockNotAvailable:
			return true
		case pgerrcode.ConnectionException:
			return true
		case pgerrcode.AdminShutdown:
			return true
		case pgerrcode.CrashShutdown:
			return true
		case pgerrcode.CannotConnectNow:
			return true
		}
	}
	return errors.Is(err, sql.ErrConnDone)
}
