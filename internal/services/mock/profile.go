package mock

import (
	"context"

	"github.com/AlkorMizar/job-hunter/internal/model/handl"
)

type UserServiceMock struct {
	MockGetUser    func(id int) (user *handl.User, err error)
	MockUpdateUSer func(id int, inf handl.UpdateInfo) error
	MockUpdatePwd  func(id int, pwd handl.Passwords) error
}

func (u *UserServiceMock) GetUser(ctx context.Context, id int) (user *handl.User, err error) {
	return u.MockGetUser(id)
}

func (u *UserServiceMock) UpdateUser(ctx context.Context, id int, inf handl.UpdateInfo) error {
	return u.MockUpdateUSer(id, inf)
}

func (u *UserServiceMock) UpdatePassword(ctx context.Context, id int, pwd handl.Passwords) error {
	return u.MockUpdatePwd(id, pwd)
}
