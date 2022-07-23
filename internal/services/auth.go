package services

import (
	"errors"
	"fmt"
	"time"

	"github.com/AlkorMizar/job-hunter/internal/handler/model"
	"github.com/AlkorMizar/job-hunter/internal/repository"
	"github.com/AlkorMizar/job-hunter/internal/util"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost = 8
	tokenTTL   = 3 * time.Hour
	issuer     = "job-hunter"
)

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
	repo       repository.UserManagment
	signingKey []byte
}

func NewAuthService(repo repository.UserManagment, sKey string) *AuthService {
	return &AuthService{
		repo:       repo,
		signingKey: []byte(sKey),
	}
}

func (s *AuthService) CreateUser(newUser *model.NewUser) (err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("in CreateUser: %w", err)
		}
	}()

	pwd, err := bcrypt.GenerateFromPassword([]byte(newUser.Password), bcryptCost)
	if err != nil {
		return fmt.Errorf("failed generate hash with error: %w", err)
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

func (s *AuthService) CreateToken(authInfo model.AuthInfo) (token string, err error) {
	defer util.Wrap(&err, "in CreateToken")

	user, err := s.repo.GetUserWithEamil(authInfo.Email)
	if err != nil {
		return "", err
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(authInfo.Password)); err != nil {
		return "", fmt.Errorf("failed compare passwords of user %d: %w", user.ID, err)
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
	token, err = jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(s.signingKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign claims with error: %w", err)
	}
	return token, err
}

func (s *AuthService) ParseToken(tokenStr string) (info model.UserInfo, err error) {
	defer func() {
		if err != nil {
			err = fmt.Errorf("in ParseToken : %w", err)
		}
	}()

	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return s.signingKey, nil
	})

	if err != nil {
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
