package services

import (
	"github.com/AlkorMizar/job-hunter/internal/handler/model"
	"github.com/AlkorMizar/job-hunter/internal/repository"
)

type Authorization interface {
	CreateUser(newUser *model.NewUser) error
	CreateToken(authInfo model.AuthInfo) (string, error)
	ParseToken(tokenStr string) (model.UserInfo, error)
}

type Service struct {
	Authorization
}

func NewService(repo *repository.Repository) *Service {
	return &Service{
		Authorization: NewAuthService(repo.UserManagment),
	}
}
