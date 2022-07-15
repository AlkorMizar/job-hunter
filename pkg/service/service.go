package service

import (
	"github.com/AlkorMizar/job-hunter/pkg/handler/model"
	"github.com/AlkorMizar/job-hunter/pkg/repository"
)

type Authorization interface {
	CreateUser(newUser *model.NewUser) error
	CreateToken(authInfo model.AuthInfo) (string, error)
	ParseToken(tokenStr string) (int, map[string]struct{}, error)
}

type Service struct {
	Authorization
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repo.UserManagment),
	}
}
