package service

import (
	"github.com/AlkorMizar/job-hunter/pkg/handler/model"
	"github.com/AlkorMizar/job-hunter/pkg/repository"
	"golang.org/x/crypto/bcrypt"
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

	user := &repository.User{
		Login:    inf.Login,
		Email:    inf.Email,
		FullName: inf.FullName,
	}

	return u.repo.UpdateUserStr(id, user)
}

func (u *UserServ) UpdatePassword(id int, pwd model.Passwords) error {
	user, err := u.repo.GetUserById(id)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(pwd.CurrPassword)); err != nil {
		return err
	}

	pwdHash, err := bcrypt.GenerateFromPassword([]byte(pwd.NewPassword), bcryptCost)

	if err != nil {
		return err
	}

	if err := u.repo.SetPassword(id, pwdHash); err != nil {
		return err
	}

	return nil
}
