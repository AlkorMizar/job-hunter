package mock

import "github.com/AlkorMizar/job-hunter/internal/repository"

type UserManagment struct {
	MockCreateUser       func(user *repository.User) error
	MockGetUserWithEamil func(email string) (repository.User, error)
	MockGetRoles         func(user *repository.User) (map[string]struct{}, error)
	MockSetRoles         func(user *repository.User) error
}

func (u *UserManagment) CreateUser(user *repository.User) error {
	return u.MockCreateUser(user)
}
func (u *UserManagment) GetUserWithEamil(email string) (repository.User, error) {
	return u.MockGetUserWithEamil(email)
}
func (u *UserManagment) GetRoles(user *repository.User) (map[string]struct{}, error) {
	return u.MockGetRoles(user)
}
func (u *UserManagment) SetRoles(user *repository.User) error {
	return u.MockSetRoles(user)
}
