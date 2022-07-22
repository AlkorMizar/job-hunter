package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/AlkorMizar/job-hunter/internal/handler/model"
	"github.com/AlkorMizar/job-hunter/internal/repository"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost = 8
	tokenTTL   = 3 * time.Hour
	issuer     = "job-hunter"
)

var signingKey = []byte("dontforgettochange")

var (
	ErrExpiredToken = errors.New("token expired")
	ErrTokenInvalid = errors.New("token has invalid format or couldn't handle it")
	ErrClaimsInvald = errors.New("claims invalid")
)

type Claims struct {
	model.UserInfo
	jwt.RegisteredClaims
}

type AuthService struct {
	repo repository.UserManagment
}

func NewAuthService(repo repository.UserManagment) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

func (s *AuthService) CreateUser(newUser *model.NewUser) error {
	pwd, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcryptCost)

	if err != nil {
		return err
	}

	//TODO: maybe move to func
	roles := make(map[string]struct{})
	for _, v := range newUser.Roles {
		roles[v] = struct{}{}
	}

	user := repository.User{
		Login:    newUser.Login,
		Password: pwd,
		Email:    newUser.Email,
		FullName: newUser.Login,
		Roles:    roles,
	}

	return s.repo.CreateUser(&user)
}

func (s *AuthService) CreateToken(authInfo model.AuthInfo) (string, error) {
	user, err := s.repo.GetUser(authInfo.Email)
	if err != nil {
		return "", err
	}

	if e := bcrypt.CompareHashAndPassword(user.Password, []byte(authInfo.Password)); e != nil {
		return "", e
	}

	user.Roles, err = s.repo.GetRoles(&user)
	if err != nil {
		return "", err
	}

	// Create the Claims
	claims := &Claims{
		UserInfo: model.UserInfo{
			ID:    user.ID,
			Roles: user.Roles,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    issuer,
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(signingKey)
}

func (s *AuthService) ParseToken(tokenStr string) (info model.UserInfo, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("during ParseToken %w", err)
		}
	}()

	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return signingKey, nil
	})

	if !token.Valid {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorExpired != 0 {
				return model.UserInfo{}, ErrExpiredToken
			}
		}
		return model.UserInfo{}, ErrTokenInvalid
	}

	claims, ok := token.Claims.(*Claims)
	if !ok {
		return model.UserInfo{}, ErrClaimsInvald
	}

	info = claims.UserInfo
	return info, nil
}
