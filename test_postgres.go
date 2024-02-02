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
)

func NewTestPostgresConnectionWithMigrations(connectionURL string, migrationURL string) (*sqlx.DB, error) {

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
