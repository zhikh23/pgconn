package pgconn

import (
	"fmt"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/jmoiron/sqlx"
)

func ConnectWithTries(cfg *ConnConfig, settings *ConnSettings, attempts int, attemptDuration time.Duration) (db *sqlx.DB, err error) {
	for ; attempts > 0; attempts-- {
		db, err = Connect(cfg, settings)
		if err == nil {
			return
		}
		time.Sleep(attemptDuration)
	}
	return
}

func Connect(cfg *ConnConfig, settings *ConnSettings) (db *sqlx.DB, err error) {
	db, err = sqlx.Open("pgx", cfg.Url())
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	setupConnection(db, settings)

	return db, nil
}

func MigrateUp(cfg ConnConfig, migrationUrl string) error {
	dbUrl := cfg.Url()

	migrator, err := migrate.New(migrationUrl, dbUrl)
	if err != nil {
		return fmt.Errorf("error migrate: %s", err)
	}

	err = migrator.Up()

	return err
}

func setupConnection(db *sqlx.DB, settings *ConnSettings) {
	db.SetMaxOpenConns(settings.MaxOpenConns)
	db.SetConnMaxLifetime(settings.ConnMaxLifetime * time.Second)
	db.SetMaxIdleConns(settings.MaxIdleConns)
	db.SetConnMaxIdleTime(settings.ConnMaxIdleTime * time.Second)
}