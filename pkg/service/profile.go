package service

import (
	"github.com/AlkorMizar/job-hunter/pkg/handler/model"
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

func (u *UserServ) GetUser(id int) (user *model.User, err error) {
	return &model.User{}, nil
}
