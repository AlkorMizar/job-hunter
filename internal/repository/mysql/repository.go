package mysql

import (
	"github.com/AlkorMizar/job-hunter/internal/logging"
	"github.com/jmoiron/sqlx"
)

type Repository struct {
	db  *sqlx.DB
	log *logging.Logger
}

func NewMysqlRepository(db *sqlx.DB, log *logging.Logger) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Close() error {
	return r.db.Close()
}
