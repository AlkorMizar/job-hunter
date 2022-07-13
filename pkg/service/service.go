package service

import "github.com/AlkorMizar/job-hunter/pkg/repository"

type Authorization interface {
	CreateUser(user repository.User) error
}

type Service struct {
	Authorization
}

func NewService(repo *repository.Repository) *Service {
	return &Service{}
}
