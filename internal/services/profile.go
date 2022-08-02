package services

import (
	"context"

	"github.com/AlkorMizar/job-hunter/internal/logging"
	"github.com/AlkorMizar/job-hunter/internal/model/handl"
	"github.com/AlkorMizar/job-hunter/internal/model/repo"
	"golang.org/x/crypto/bcrypt"
)

type UserServ struct {
	repo       UserManagment
	bcryptCost int
	log        *logging.Logger
}

func NewUserServ(repo UserManagment, bcryptCost int, log *logging.Logger) *UserServ {
	return &UserServ{
		repo:       repo,
		bcryptCost: bcryptCost,
		log:        log,
	}
}

func (u *UserServ) GetUser(ctx context.Context, id int) (*handl.User, error) {

	userDB, err := u.repo.GetUserByID(id)
	if err != nil {
		return nil, err
	}

	roles := make([]string, 0)
	for k := range userDB.Roles {
		roles = append(roles, k)
	}

	user := &handl.User{
		Login:    userDB.Login,
		Email:    userDB.Email,
		FullName: userDB.FullName,
		Roles:    roles,
	}
	return user, nil
}

func (u *UserServ) UpdateUser(ctx context.Context, id int, inf handl.UpdateInfo) error {

	user := &repo.User{
		Login:    inf.Login,
		Email:    inf.Email,
		FullName: inf.FullName,
	}

	return u.repo.UpdateUser(id, user)
}

func (u *UserServ) UpdatePassword(ctx context.Context, id int, pwd handl.Passwords) error {
	user, err := u.repo.GetUserByID(id)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword(user.Password, []byte(pwd.CurrPassword)); err != nil {
		return err
	}

	pwdHash, err := bcrypt.GenerateFromPassword([]byte(pwd.NewPassword), u.bcryptCost)

	if err != nil {
		return err
	}

	if err := u.repo.UpdatePassword(id, pwdHash); err != nil {
		return err
	}

	return nil
}
