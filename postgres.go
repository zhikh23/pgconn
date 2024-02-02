package postgres

import (
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	"log"
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

func NewPostgresConnectionWithMigrations(cfg Config, absoluteLink string) (client *PSQLClient, err error) {
	client, err = NewPostgresConnection(cfg)
	if err != nil {
		return nil, err
	}

	if err := UpMigrations(client.Db, cfg.DBName, absoluteLink); err != nil {
		log.Fatal(err)
	}

	return client, nil
}

func UpMigrations(sqlxDb *sqlx.DB, dbName string, absoluteLink string) error {
	driver, err := postgres.WithInstance(sqlxDb.DB, &postgres.Config{})
	if err != nil {
		return err
	}

	//pwd, _ := os.Getwd()
	m, err := migrate.NewWithDatabaseInstance(
		absoluteLink,
		dbName, driver)
	if err != nil {
		return err
	}
	if err := m.Up(); err != nil {
		log.Fatal(err)
	}
	return nil
}

func DownMigrations(sqlxDb *sqlx.DB, dbName string, absoluteLink string) error {
	driver, err := postgres.WithInstance(sqlxDb.DB, &postgres.Config{})
	if err != nil {
		return err
	}
	//pwd, _ := os.Getwd()
	m, err := migrate.NewWithDatabaseInstance(
		absoluteLink,
		dbName, driver)
	if err != nil {
		return err
	}
	if err := m.Down(); err != nil {
		return err
	}
	return nil
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
