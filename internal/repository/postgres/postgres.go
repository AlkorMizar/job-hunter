package postgres

import (
	"fmt"
	"log"
	"time"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
)

const (
	connMaxLifetime = time.Minute * 3
	maxOpenConns    = 10
	maxIdleConns    = 10
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	DBName   string
	Options  string
}

func NewPodtgresDB(cfg *Config) (*sqlx.DB, error) {
	//postgresql://[user[:password]@][netloc][:port][/dbname]
	db, err := sqlx.Connect("pgx", fmt.Sprintf("postgres://%s:%s@%s:%s/%s?%s",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.DBName, cfg.Options))
	if err != nil {
		log.Fatalln(err)
	}

	db.SetConnMaxLifetime(connMaxLifetime)
	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)

	return db, nil
}
