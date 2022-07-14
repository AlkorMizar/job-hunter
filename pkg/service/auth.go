package service

import (
	"time"

	"github.com/AlkorMizar/job-hunter/pkg/handler/model"
	"github.com/AlkorMizar/job-hunter/pkg/repository"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost = 8
	signingKey = "qrkjk#4#%35FSFJlja#4353KSFjH"
	tokenTTL   = 1 * time.Hour
)

type UserService struct {
	repo repository.UserManagment
}

type tokenClaims struct {
	jwt.StandardClaims
	UserId int `json:"user_id"`
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

func (s *UserService) CreateToken(authInfo model.AuthInfo) (string, error) {
	pwd, err := bcrypt.GenerateFromPassword([]byte(authInfo.Password), bcryptCost)

	if err != nil {
		return "", err
	}

	user, err := s.repo.GetUser(authInfo.Email, pwd)
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &tokenClaims{
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(tokenTTL).Unix(),
			IssuedAt:  time.Now().Unix(),
		},
		user.Id,
	})

	return token.SignedString([]byte(signingKey))
}
