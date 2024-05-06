package pgconn

import (
	"fmt"
	"time"
)

const (
	SslModeDisable = "disable"
	SslModeEnable  = "enable"
)

type ConnConfig struct {
	Port     string
	Host     string
	User     string
	Password string
	DbName   string
	SslMode  string
}

type ConnSettings struct {
	MaxOpenConns    int
	ConnMaxLifetime time.Duration
	MaxIdleConns    int
	ConnMaxIdleTime time.Duration
}

func (cfg *ConnConfig) Url() string {
	return fmt.Sprintf(
		"%s://%s:%s/%s?sslmode=%s&user=%s&password=%s",
		"postgres", cfg.Host, cfg.Port, cfg.DbName, cfg.SslMode, cfg.User, cfg.Password,
	)
}
