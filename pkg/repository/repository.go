package repository

import "github.com/jmoiron/sqlx"

type UserManagment interface {
	CreateUser(user *User) error
	GetUser(email string) (User, error)
	GetRoles(user *User) (map[string]struct{}, error)
	SetRoles(user *User) error
	GetUserFromId(id int) (user *User, err error)
	UpdateUserStr(id int, updateU *User) (err error)
	SetPassword(id int, pwd string) error
	GetUserById(id int) (User, error)
}

type Repository struct {
	UserManagment
}

func NewRepository(db *sqlx.DB) *Repository {
	return &Repository{
		UserManagment: NewUserManagMsql(db),
	}
}
