package services

import "github.com/AlkorMizar/job-hunter/internal/model/repo"

type UserManagment interface {
	CreateUser(user *repo.User) error
	GetUserWithEamil(email string) (repo.User, error)
	GetRoles(user *repo.User) (map[string]struct{}, error)
	SetRoles(user *repo.User) error
	GetUserByID(ID int) (*repo.User, error)
	UpdateUser(ID int, user *repo.User) error
	UpdatePassword(ID int, pwd []byte) error
}

type Repository interface {
	UserManagment
	Close() error
}
