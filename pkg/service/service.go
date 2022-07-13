package service

import (
	"github.com/AlkorMizar/job-hunter/pkg/handler/model"
	"github.com/AlkorMizar/job-hunter/pkg/repository"
)

type UserManagment interface {
	CreateUser(newUser model.NewUser) error
}

type Service struct {
	UserManagment
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		UserManagment: NewUserService(repo.UserManagment),
	}
}
