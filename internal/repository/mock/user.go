package mock

import "github.com/AlkorMizar/job-hunter/internal/model/repo"

type UserManagment struct {
	MockCreateUser       func(user *repo.User) error
	MockGetUserWithEamil func(email string) (repo.User, error)
	MockGetRoles         func(user *repo.User) (map[string]struct{}, error)
	MockSetRoles         func(user *repo.User) error
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
