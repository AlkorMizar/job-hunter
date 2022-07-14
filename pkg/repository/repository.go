package repository

import "github.com/jmoiron/sqlx"

type UserManagment interface {
	CreateUser(user User) error
	GetUser(email string) (User, error)
	GetRoles(user User) (map[string]struct{}, error)
}

type Repository struct {
	UserManagment
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		UserManagment: NewUserManagMsql(db),
	}
}
