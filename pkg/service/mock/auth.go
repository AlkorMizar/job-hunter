package mock

import (
	"github.com/AlkorMizar/job-hunter/pkg/handler/model"
)

type AuthServiceMock struct {
	MockCreateToken func(model.AuthInfo) (string, error)
	MockCreateUser  func(newUser *model.NewUser) error
}

func (s *AuthServiceMock) CreateUser(newUser *model.NewUser) error {
	return s.MockCreateUser(newUser)
}

func (s *AuthServiceMock) CreateToken(authInfo model.AuthInfo) (string, error) {
	return s.MockCreateToken(authInfo)
}

func (s *AuthServiceMock) ParseToken(tokenStr string) (id int, role map[string]struct{}, err error) {
	return 0, nil, nil
}
