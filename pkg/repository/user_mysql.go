package repository

import (
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"
)

type User struct {
	Id          int       `json:"-" db:"idUser"`
	Login       string    `json:"login" db:"login"`
	FullName    string    `json:"fullName" db:"fullName"`
	Email       string    `json:"email" db:"email"`
	Password    []byte    `json:"password" db:"password"`
	DateCreated time.Time `json:"datecreated" db:"dateCreated"`
	LastCheck   time.Time `json:"lastcheck" db:"lastCheck"`
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
	query := "INSERT INTO user (login, email, password, fullName) values (:login,:email,:password,:fullName)"
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

func (r *UserManagMysql) GetUser(email string) (User, error) {
	var user User
	query := "SELECT * FROM user WHERE email=?"
	err := r.db.Get(&user, query, email)
	return user, err
}
