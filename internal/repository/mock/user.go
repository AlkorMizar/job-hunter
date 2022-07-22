package mock

import "github.com/AlkorMizar/job-hunter/internal/repository"

type UserManagment struct {
	MockCreateUser func(user *repository.User) error
	MockGetUser    func(email string) (repository.User, error)
	MockGetRoles   func(user *repository.User) (map[string]struct{}, error)
	MockSetRoles   func(user *repository.User) error
}

func (u *UserManagment) CreateUser(user *repository.User) error {
	return u.MockCreateUser(user)
}
func (u *UserManagment) GetUser(email string) (repository.User, error) {
	return u.MockGetUser(email)
}
func (u *UserManagment) GetRoles(user *repository.User) (map[string]struct{}, error) {
	return u.MockGetRoles(user)
}
func (u *UserManagment) SetRoles(user *repository.User) error {
	return u.MockSetRoles(user)
}
