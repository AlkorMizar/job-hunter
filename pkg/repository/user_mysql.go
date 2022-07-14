package repository

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type User struct {
	Id       int    `json:"-" db:"id"`
	Login    string `json:"login" binding:"required"`
	FullName string `json:"fullName" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password []byte `json:"password" binding:"required"`
}

type UserManagMysql struct {
	db sqlx.DB
}

func NewUserManagMsql(db *sqlx.DB) *UserManagMysql {
	return &UserManagMysql{
		db: *db,
	}
}

func (r *UserManagMysql) CreateUser(user User) error {
	query := "INSERT INTO user (login, email, password, fullName) values (:login,:email,:password,:fullname)"
	res, err := r.db.NamedExec(query, user)

	if err != nil {
		return err
	}

	num, err := res.RowsAffected()

	if err != nil {
		return err
	}

	if num != 1 {
		return fmt.Errorf("couldn't insert")
	}

	return nil
}

func (r *UserManagMysql) GetUser(email string, password []byte) (User, error) {
	var user User
	query := "SELECT * FROM user WHERE email=? AND password=?"
	err := r.db.Get(&user, query, email, password)

	return user, err
}
