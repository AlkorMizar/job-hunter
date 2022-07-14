package repository

import "github.com/jmoiron/sqlx"

type UserManagment interface {
	CreateUser(user User) error
	GetUser(email string, password []byte) (User, error)
}

type Repository struct {
	UserManagment
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		UserManagment: NewUserManagMsql(db),
	}
}
