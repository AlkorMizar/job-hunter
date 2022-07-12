package service

import "github.com/AlkorMizar/job-hunter/pkg/repository"

type Service struct {
}

func NewService(repo *repository.Repository) *Service {
	return &Service{}
}
