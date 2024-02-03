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

func Connect(url string, settingsOrNil *ConnSettings) (db *sqlx.DB, err error) {
	db, err = sqlx.Open("pgx", url)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	if settingsOrNil != nil {
		setupConnection(db, settingsOrNil)
	}

	return db, nil
}

func ConnectWithTries(url string, settings *ConnSettings, attempts int, attemptDuration time.Duration) (db *sqlx.DB, err error) {
	for ; attempts > 0; attempts-- {
		db, err = Connect(url, settings)
		if err == nil {
			return
		}
		time.Sleep(attemptDuration)
	}
	return
}

func MigrateUp(dbUrl, migrationUrl string) error {
	migrator, err := migrate.New(migrationUrl, dbUrl)
	if err != nil {
		return fmt.Errorf("error migrate: %s", err)
	}
	defer migrator.Close()

	err = migrator.Up()

	return err
}

func setupConnection(db *sqlx.DB, settings *ConnSettings) {
	db.SetMaxOpenConns(settings.MaxOpenConns)
	db.SetConnMaxLifetime(settings.ConnMaxLifetime * time.Second)
	db.SetMaxIdleConns(settings.MaxIdleConns)
	db.SetConnMaxIdleTime(settings.ConnMaxIdleTime * time.Second)
}
