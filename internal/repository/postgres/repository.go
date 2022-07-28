package postgres

import (
	"github.com/AlkorMizar/job-hunter/internal/logging"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db  *sqlx.DB
	log *logging.Logger
}

func NewPostgresRepository(db *sqlx.DB, log *logging.Logger) *Repository {
	return &Repository{
		db:  db,
		log: log,
	}
}

func (r *Repository) Close() error {
	return r.db.Close()
}
