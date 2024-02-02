package postgres

import (
	"context"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/golang-migrate/migrate/v4/source/github"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"github.com/pkg/errors"
	"strings"
)

func NewTestPostgresConnectionWithString(connectionURL string, migrationURL string) (*sqlx.DB, error) {

	database, err := sqlx.ConnectContext(context.Background(), "pgx", connectionURL)
	if err != nil {
		return nil, err
	}

	MustUpTestMigrations(migrationURL, connectionURL)

	if err = database.Ping(); err != nil {
		return nil, err
	}

	return database, nil
}

func MustUpTestMigrations(sourceURL, connectionString string) {
	migrations, err := migrate.New(sourceURL, connectionString)
	if err != nil {
		panic(fmt.Sprintf("error connect to db: %s", err.Error()))
	}

	if err = migrations.Up(); err != nil {
		if errors.Is(err, migrate.ErrNoChange) {
			fmt.Println("no migrations to apply")
		} else {
			panic(fmt.Sprintf("error apply migrations: %s", err.Error()))
		}
	}
}

func NewTestPostgresConnection(host, port, user, password, dbName, absoluteLink string) (*sqlx.DB, error) {
	connectionURL := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s",
		host,
		port,
		user,
		password,
		dbName,
	)

	database, err := sqlx.ConnectContext(context.Background(), "pgx", connectionURL)
	if err != nil {
		return nil, err
	}
	err = UpMigrations(database, dbName, absoluteLink)
	if err != nil && !strings.Contains(err.Error(), "no change") {
		fmt.Println("Migration error: " + err.Error())
	}
	if err = database.Ping(); err != nil {
		return nil, err
	}

	return database, nil
}

const (
	host     = "localhost"
	port     = "5434"
	user     = "thestrikem"
	password = "123"
	dbName   = "test"
)

func MustTestPostgresInstance(host, port, user, password, dbName, absoluteLink string) *sqlx.DB {
	conn, err := NewTestPostgresConnection(host, port, user, password, dbName, absoluteLink)
	if err != nil {
		panic("Postgres instance: " + err.Error())
	}
	return conn
}

func MustTestDownMigrations(db *sqlx.DB, dbName, absoluteLink string) {
	if err := DownMigrations(db, dbName, absoluteLink); err != nil {
		panic("Down migration: " + err.Error())
	}
}
