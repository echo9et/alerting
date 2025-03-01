package entities

import (
	"database/sql"
	"errors"
	"time"

	"github.com/jackc/pgerrcode"
	"github.com/jackc/pgx/v5/pgconn"
)

func Retry(fn func() error) error {
	var err error
	delays := []time.Duration{1 * time.Second, 2 * time.Second, 5 * time.Second}
	for _, delay := range delays {
		if err = fn(); err != nil {
			if !isRunReplay(err) {
				break
			}
			time.Sleep(delay)
		} else {
			return nil
		}
	}
	return err
}

func isRunReplay(err error) bool {
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
