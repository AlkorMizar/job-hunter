package service

import (
	"github.com/AlkorMizar/job-hunter/pkg/handler/model"
	"github.com/AlkorMizar/job-hunter/pkg/repository"
)

type UserServ struct {
	repo repository.UserManagment
}

func NewUserServ(repo repository.UserManagment) *UserServ {
	return &UserServ{
		repo: repo,
	}
}

func (u *UserServ) GetUser(id int) (*model.User, error) {

	userDB, err := u.repo.GetUserFromId(id)
	if err != nil {
		return nil, err
	}

	roles := make([]string, 0)
	for k := range userDB.Roles {
		roles = append(roles, k)
	}

	user := &model.User{
		Login:    userDB.Login,
		Email:    userDB.Email,
		FullName: userDB.FullName,
		Roles:    roles,
	}
	return user, nil
}

func (u *UserServ) UpdateUser(id int, inf model.UpdateInfo) error {
	return nil
}
