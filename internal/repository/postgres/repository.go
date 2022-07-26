package postgres

import "github.com/jmoiron/sqlx"

type Repository struct {
	db *sqlx.DB
}

func NewPostgresRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Close() error {
	return r.db.Close()
}
