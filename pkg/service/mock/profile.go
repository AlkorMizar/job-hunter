package mock

import "github.com/AlkorMizar/job-hunter/pkg/handler/model"

type UserServiceMock struct {
	MockGetUser    func(id int) (user *model.User, err error)
	MockUpdateUSer func(id int, inf model.UpdateInfo) error
}

func (u *UserServiceMock) GetUser(id int) (user *model.User, err error) {
	return u.MockGetUser(id)
}

func (u *UserServiceMock) UpdateUser(id int, inf model.UpdateInfo) error {
	return u.MockUpdateUSer(id, inf)
}
