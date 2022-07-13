package service

import (
	"github.com/AlkorMizar/job-hunter/pkg/handler/model"
	"github.com/AlkorMizar/job-hunter/pkg/repository"
	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost = 8
)

type UserService struct {
	repo repository.UserManagment
}

func NewUserService(repo repository.UserManagment) *UserService {
	return &UserService{
		repo: repo,
	}
}

func (s *UserService) CreateUser(newUser model.NewUser) error {
	pwd, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcryptCost)

	if err != nil {
		return err
	}

	user := repository.User{
		Login:    newUser.Login,
		Password: pwd,
		Email:    newUser.Email,
		FullName: newUser.Login,
	}

	return s.repo.CreateUser(user)
}
