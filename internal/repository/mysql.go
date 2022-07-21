package repository

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql" //nolint:blank-imports // only for sqlx
	"github.com/jmoiron/sqlx"
)

const (
	ConnMaxLifetime = time.Minute * 3
	MaxOpenConns    = 10
	MaxIdleConns    = 10
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	Protocol string
	Options  string
}

func NewMySQLDB(cfg *Config) (*sqlx.DB, error) {
	db, err := sqlx.Open("mysql", fmt.Sprintf("%s:%s@%s(%s:%s)/%s?%s",
		cfg.Username, cfg.Password, cfg.Protocol, cfg.Host, cfg.Port, cfg.DBName, cfg.Options))
	if err != nil {
		return nil, err
	}

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(ConnMaxLifetime)
	db.SetMaxOpenConns(MaxOpenConns)
	db.SetMaxIdleConns(MaxIdleConns)

	return db, nil
}
