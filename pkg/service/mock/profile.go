package mock

import "github.com/AlkorMizar/job-hunter/pkg/handler/model"

type UserServiceMock struct {
	MockGetUser func(id int) (user *model.User, err error)
}

func (u *UserServiceMock) GetUser(id int) (user *model.User, err error) {
	return u.MockGetUser(id)
}
