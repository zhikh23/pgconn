package postgres

import (
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"time"
)

type PSQLClient struct {
	Db *sqlx.DB
}

func NewPostgresConnection(cfg Config) (client *PSQLClient, err error) {
	var database *sqlx.DB

	if err := DoWithTries(func() error {
		database, err = sqlx.Open("pgx", cfg.CreateUrl())
		if err != nil {
			return err
		}

		database.SetMaxOpenConns(cfg.Settings.MaxOpenConns)
		database.SetConnMaxLifetime(cfg.Settings.ConnMaxLifetime * time.Second)
		database.SetMaxIdleConns(cfg.Settings.MaxIdleConns)
		database.SetConnMaxIdleTime(cfg.Settings.ConnMaxIdleTime * time.Second)

		if err = database.Ping(); err != nil {
			return err
		}

		return nil
	}, 4, 2*time.Second); err != nil {
		return nil, err
	}

	return &PSQLClient{
		Db: database,
	}, nil
}

func DoWithTries(fn func() error, attempts int, delay time.Duration) (err error) {
	for attempts > 0 {
		if err = fn(); err != nil {
			time.Sleep(delay)
			attempts--
			continue
		}
		return nil
	}
	return
}
