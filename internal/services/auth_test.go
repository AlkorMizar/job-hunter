package services_test

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/AlkorMizar/job-hunter/internal/handler/model"
	"github.com/AlkorMizar/job-hunter/internal/repository"
	"github.com/AlkorMizar/job-hunter/internal/repository/mock"
	"github.com/AlkorMizar/job-hunter/internal/services"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

const (
	tokenTTL = time.Hour
)

var (
	validUserInf = model.UserInfo{
		ID:    1,
		Roles: map[string]struct{}{"role": {}},
	}
	invalidUserInf = model.UserInfo{}
	signingKey     = []byte("testSigningKey")
)

func TestParseToken(t *testing.T) {

	tests := []struct {
		name  string
		token func() string
		uInfo model.UserInfo
		err   error
	}{
		{
			"ok",
			func() string {
				claims := &services.Claims{
					UserInfo: validUserInf,
					RegisteredClaims: jwt.RegisteredClaims{
						ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenTTL)),
						IssuedAt:  jwt.NewNumericDate(time.Now()),
						Issuer:    "test",
					},
				}

				token, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(signingKey)
				return token
			},
			validUserInf,
			nil,
		},
		{
			"expired token",
			func() string {
				claims := &services.Claims{
					UserInfo: validUserInf,
					RegisteredClaims: jwt.RegisteredClaims{
						ExpiresAt: jwt.NewNumericDate(time.Now().Add(-tokenTTL)),
						IssuedAt:  jwt.NewNumericDate(time.Now()),
						Issuer:    "test",
					},
				}

				token, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(signingKey)
				return token
			},
			invalidUserInf,
			services.ErrExpiredToken,
		},
		{
			"empty token",
			func() string {
				return ""
			},
			invalidUserInf,
			services.ErrTokenInvalid,
		},
		{
			"invalid token",
			func() string {
				claims := &services.Claims{
					UserInfo: validUserInf,
					RegisteredClaims: jwt.RegisteredClaims{
						ExpiresAt: jwt.NewNumericDate(time.Now().Add(-tokenTTL)),
						IssuedAt:  jwt.NewNumericDate(time.Now()),
						Issuer:    "test",
					},
				}

				token, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(signingKey)
				token = "/" + token[1:]
				return token
			},
			invalidUserInf,
			services.ErrTokenInvalid,
		},
		{
			"invalid claim",
			func() string {
				claims := &jwt.RegisteredClaims{
					ExpiresAt: jwt.NewNumericDate(time.Now().Add(-tokenTTL)),
					IssuedAt:  jwt.NewNumericDate(time.Now()),
					Issuer:    "test",
				}

				token, _ := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(signingKey)
				return token
			},
			invalidUserInf,
			services.ErrExpiredToken,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			auth := services.NewAuthService(&mock.UserManagment{}, string(signingKey))

			info, err := auth.ParseToken(test.token())

			if !errors.Is(err, test.err) {
				t.Fatalf("got %v want %v", err, test.err)
			}

			if !reflect.DeepEqual(info, test.uInfo) {
				t.Fatalf("got %v want %v", err, test.err)
			}
		})
	}
}

func TestCreateUSer(t *testing.T) {
	validNewUser := &model.NewUser{
		Login:    "login",
		Email:    "email@gmail.com",
		Roles:    []string{"role1", "role2"},
		Password: "password",
	}

	errAlreadyExists := errors.New("user already exists")

	tests := []struct {
		name           string
		newUser        *model.NewUser
		mockCreateUser func(user *repository.User) error
		err            error
	}{
		{
			"ok",
			validNewUser,
			func(user *repository.User) error {
				if !compare(validNewUser, user) {
					return fmt.Errorf("incorrect repository.User data")
				}
				return nil
			},
			nil,
		},
		{
			"user already exist",
			validNewUser,
			func(user *repository.User) error {
				return errAlreadyExists
			},
			errAlreadyExists,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			auth := services.NewAuthService(&mock.UserManagment{
				MockCreateUser: test.mockCreateUser,
			}, string(signingKey))

			err := auth.CreateUser(test.newUser)

			if !errors.Is(err, test.err) {
				t.Fatalf("got %v want %v", err, test.err)
			}
		})
	}
}

func compare(newUser *model.NewUser, user *repository.User) bool {
	for _, v := range newUser.Roles {
		if _, ok := user.Roles[v]; !ok {
			return false
		}
	}

	if e := bcrypt.CompareHashAndPassword(user.Password, []byte(newUser.Password)); e != nil {
		return false
	}

	return newUser.Login == user.Login && newUser.Email == user.Email
}
