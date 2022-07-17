package service

import (
	"github.com/AlkorMizar/job-hunter/pkg/repository"
)

type UserServ struct {
	repo repository.UserManagment
}

func NewUserServ(repo repository.UserManagment) *UserServ {
	return &UserServ{
		repo: repo,
	}
}

func (u *UserServ) GetUser(id int) (user *repository.User, err error) {
	return &repository.User{}, nil
}
