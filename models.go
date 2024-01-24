package postgres

import "time"

type Config struct {
	Port     string
	Host     string
	User     string
	Password string
	DBName   string
	SSLMode  string
	PGDriver string
	Settings Settings
}

type Settings struct {
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
	MaxIdleConns    int
	ConnMaxIdleTime time.Duration
}
