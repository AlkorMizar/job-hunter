package mock

import (
	"github.com/AlkorMizar/job-hunter/internal/model/handl"
)

type AuthServiceMock struct {
	MockCreateToken func(handl.AuthInfo) (string, error)
	MockCreateUser  func(newUser *handl.NewUser) error
}

func (s *AuthServiceMock) CreateUser(newUser *handl.NewUser) error {
	return s.MockCreateUser(newUser)
}

func (s *AuthServiceMock) CreateToken(authInfo handl.AuthInfo) (string, error) {
	return s.MockCreateToken(authInfo)
}

func (s *AuthServiceMock) ParseToken(tokenStr string) (uInfo handl.UserInfo, err error) {
	return uInfo, nil
}
