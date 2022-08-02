package mock

import "github.com/AlkorMizar/job-hunter/internal/model/repo"

type UserManagment struct {
	MockCreateUser       func(user *repo.User) error
	MockGetUserWithEamil func(email string) (repo.User, error)
	MockGetRoles         func(user *repo.User) (map[string]struct{}, error)
	MockSetRoles         func(user *repo.User) error
	MockGetUserByID      func(int) (*repo.User, error)
	MockUpdateUser       func(int, *repo.User) error
	MockUpdatePassword   func(int, []byte) error
}

func (u *UserManagment) CreateUser(user *repo.User) error {
	return u.MockCreateUser(user)
}

func (u *UserManagment) GetUserWithEamil(email string) (repo.User, error) {
	return u.MockGetUserWithEamil(email)
}

func (u *UserManagment) GetRoles(user *repo.User) (map[string]struct{}, error) {
	return u.MockGetRoles(user)
}

func (u *UserManagment) SetRoles(user *repo.User) error {
	return u.MockSetRoles(user)
}

func (u *UserManagment) GetUserByID(ID int) (*repo.User, error) {
	return u.MockGetUserByID(ID)
}

func (u *UserManagment) UpdateUser(ID int, user *repo.User) error {
	return u.MockUpdateUser(ID, user)
}

func (u *UserManagment) UpdatePassword(ID int, pwd []byte) error {
	return u.MockUpdatePassword(ID, pwd)
}
