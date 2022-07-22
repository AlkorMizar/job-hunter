package services_test

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/AlkorMizar/job-hunter/internal/handler/model"
	"github.com/AlkorMizar/job-hunter/internal/repository"
	"github.com/AlkorMizar/job-hunter/internal/repository/mock"
	"github.com/AlkorMizar/job-hunter/internal/services"
	"github.com/golang-jwt/jwt/v4"
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

			auth := services.NewAuthService(&mock.UserManagment{
				MockCreateUser: func(user *repository.User) error {
					return nil
				},
			}, string(signingKey))

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
