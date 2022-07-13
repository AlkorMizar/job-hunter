package repository

import "github.com/jmoiron/sqlx"

type User struct {
	Id       int    `json:"-" db:"id"`
	Login    string `json:"login" binding:"required"`
	FullName string `json:"fullName" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password []byte `json:"password" binding:"required"`
}

type UserManagment struct {
	db sqlx.DB
}

func NewUserManagment(db *sqlx.DB) *UserManagment {
	return &UserManagment{
		db: *db,
	}
}

func (um *UserManagment) CreateUser(user User) error {
	return nil
}
