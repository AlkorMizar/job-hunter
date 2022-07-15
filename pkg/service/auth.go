package service

import (
	"fmt"
	"log"
	"time"

	"github.com/AlkorMizar/job-hunter/pkg/handler/model"
	"github.com/AlkorMizar/job-hunter/pkg/repository"
	"github.com/golang-jwt/jwt/v4"
	"golang.org/x/crypto/bcrypt"
)

const (
	bcryptCost = 8
	tokenTTL   = 15 * time.Minute
)

var signingKey = []byte("dontforgettochange")

type Claims struct {
	UserID int                 `json:"userId"`
	Roles  map[string]struct{} `json:"roles"`
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

	/////////////////////////////////////////
	roles := make(map[string]struct{})
	for _, v := range newUser.Roles {
		roles[v] = struct{}{}
	}
	////////////////////////////////////////

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

	log.Printf("User credentials %d %v", user.ID, user.Roles)

	// Create the Claims
	claims := &Claims{
		UserID: user.ID,
		Roles:  user.Roles,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(tokenTTL)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "api_token",
		},
	}

	// Declare the token with the algorithm used for signing, and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	return token.SignedString(signingKey)
}

func (s *AuthService) ParseToken(tokenStr string) (id int, roles map[string]struct{}, err error) {
	token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return signingKey, nil
	})

	if err != nil {
		return 0, nil, err
	}

	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims.UserID, claims.Roles, nil
	}

	return 0, nil, fmt.Errorf("invalid token")
}
