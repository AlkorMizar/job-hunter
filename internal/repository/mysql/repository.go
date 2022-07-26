package mysql

import "github.com/jmoiron/sqlx"

type Repository struct {
	db *sqlx.DB
}

func NewMysqlRepository(db *sqlx.DB) *Repository {
	return &Repository{
		db: db,
	}
}

func (r *Repository) Close() error {
	return r.db.Close()
}
