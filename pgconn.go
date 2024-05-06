package pgconn

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/jmoiron/sqlx"
)

func Connect(context context.Context, url string, settingsOrNil *ConnSettings) (db *sqlx.DB, err error) {
	db, err = sqlx.ConnectContext(context, "pgx", url)
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

func ConnectWithTries(context context.Context, url string, settings *ConnSettings, attempts int, attemptDuration time.Duration) (db *sqlx.DB, err error) {
	for ; attempts > 0; attempts-- {
		db, err = Connect(context, url, settings)
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

	defer func(migrator *migrate.Migrate) {
		errSource, errDatabase := migrator.Close()
		if errSource != nil || errDatabase != nil {
			slog.Error(
				fmt.Sprintf(
					"MIGRATION ERROR: source:%s;database:%s",
					errSource.Error(),
					errDatabase.Error(),
				),
			)
		}
	}(migrator)

	err = migrator.Up()
	return err
}

func setupConnection(db *sqlx.DB, settings *ConnSettings) {
	db.SetMaxOpenConns(settings.MaxOpenConns)
	db.SetConnMaxLifetime(settings.ConnMaxLifetime * time.Second)
	db.SetMaxIdleConns(settings.MaxIdleConns)
	db.SetConnMaxIdleTime(settings.ConnMaxIdleTime * time.Second)
}
