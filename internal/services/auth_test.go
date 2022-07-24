package services_test

import (
	"errors"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/AlkorMizar/job-hunter/internal/model/handl"
	"github.com/AlkorMizar/job-hunter/internal/model/repo"
	"github.com/AlkorMizar/job-hunter/internal/repository/mock"
	"github.com/AlkorMizar/job-hunter/internal/services"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

const (
	tokenTTL = time.Hour
)

var (
	signingKey = []byte("testSigningKey")
)

func TestParseToken(t *testing.T) {

	validUserInf := handl.UserInfo{
		ID:    1,
		Roles: map[string]struct{}{"role": {}},
	}
	invalidUserInf := handl.UserInfo{}

	tests := []struct {
		name  string
		token func() string
		uInfo handl.UserInfo
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
	validNewUser := &handl.NewUser{
		Login:    "login",
		Email:    "email@gmail.com",
		Roles:    []string{"role1", "role2"},
		Password: "password",
	}

	errAlreadyExists := errors.New("user already exists")

	tests := []struct {
		name           string
		newUser        *handl.NewUser
		mockCreateUser func(user *repo.User) error
		err            error
	}{
		{
			"ok",
			validNewUser,
			func(user *repo.User) error {
				if !compare(validNewUser, user) {
					return fmt.Errorf("incorrect repo.User data")
				}
				return nil
			},
			nil,
		},
		{
			"user already exist",
			validNewUser,
			func(user *repo.User) error {
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

func TestCreateToken(t *testing.T) {
	errNotFound := errors.New("not found")
	errGetRolesFailed := errors.New("faileg GetRoles")

	pwd, _ := bcrypt.GenerateFromPassword([]byte("valid"), 8)

	user := repo.User{
		ID:       1,
		Password: pwd,
	}

	roles := map[string]struct{}{"role1": {}, "role2": {}}

	invalUInfo := handl.UserInfo{}

	tests := []struct {
		name         string
		authInfo     handl.AuthInfo
		mockGetUser  func(email string) (repo.User, error)
		mockGetRoles func(user *repo.User) (map[string]struct{}, error)
		err          error
		uInfo        handl.UserInfo
	}{
		{
			"ok",
			handl.AuthInfo{
				Email:    "test@gmail.com",
				Password: "valid",
			},
			func(email string) (repo.User, error) {
				return user, nil
			},
			func(user *repo.User) (map[string]struct{}, error) {
				return roles, nil
			},
			nil,
			handl.UserInfo{
				ID:    user.ID,
				Roles: roles,
			},
		},
		{
			"ok,empty roles",
			handl.AuthInfo{
				Email:    "test@gmail.com",
				Password: "valid",
			},
			func(email string) (repo.User, error) {
				return user, nil
			},
			func(user *repo.User) (map[string]struct{}, error) {
				return map[string]struct{}{}, nil
			},
			nil,
			handl.UserInfo{
				ID:    user.ID,
				Roles: map[string]struct{}{},
			},
		},
		{
			"user not found",
			handl.AuthInfo{
				Email:    "not@found.com",
				Password: "valid",
			},
			func(email string) (repo.User, error) {
				return repo.User{}, errNotFound
			},
			func(user *repo.User) (map[string]struct{}, error) {
				return nil, fmt.Errorf("impossible error")
			},
			errNotFound,
			invalUInfo,
		},
		{
			"invalid password",
			handl.AuthInfo{
				Email:    "test@gmail.com",
				Password: "invalid",
			},
			func(email string) (repo.User, error) {
				return user, nil
			},
			func(user *repo.User) (map[string]struct{}, error) {
				return nil, fmt.Errorf("impossible error")
			},
			bcrypt.ErrMismatchedHashAndPassword,
			invalUInfo,
		},
		{
			"error in roles",
			handl.AuthInfo{
				Email:    "test@gmail.com",
				Password: "valid",
			},
			func(email string) (repo.User, error) {
				return user, nil
			},
			func(user *repo.User) (map[string]struct{}, error) {
				return nil, errGetRolesFailed
			},
			errGetRolesFailed,
			invalUInfo,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			authServ := services.NewAuthService(&mock.UserManagment{
				MockGetRoles:         test.mockGetRoles,
				MockGetUserWithEamil: test.mockGetUser,
			}, string(signingKey))

			tokenStr, err := authServ.CreateToken(test.authInfo)

			if !errors.Is(err, test.err) {
				t.Fatalf("got %v want %v", err, test.err)
			}

			if reflect.DeepEqual(test.uInfo, invalUInfo) {
				return
			}

			token, err := jwt.ParseWithClaims(tokenStr, &services.Claims{}, func(token *jwt.Token) (interface{}, error) {
				return signingKey, nil
			})

			if err != nil {
				t.Fatalf("got unexpected error %v", err)
			}

			claims, ok := token.Claims.(*services.Claims)
			if !ok {
				t.Fatalf("got unexpected token %v", token)
			}

			if !reflect.DeepEqual(test.uInfo, claims.UserInfo) {
				t.Fatalf("got %v want %v", claims.UserInfo, test.uInfo)
			}
		})
	}
}

func compare(newUser *handl.NewUser, user *repo.User) bool {
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
