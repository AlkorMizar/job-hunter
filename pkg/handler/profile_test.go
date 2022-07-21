package handler

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AlkorMizar/job-hunter/pkg/handler/model"
	"github.com/AlkorMizar/job-hunter/pkg/service"
	"github.com/AlkorMizar/job-hunter/pkg/service/mock"
)

func TestGetUser(t *testing.T) {
	tests := []struct {
		name               string
		userInf            interface{}
		mock               func(int) (*model.User, error)
		expectedStatusCode int
	}{
		{
			"ok",
			userInfo{
				1,
				make(map[string]struct{}),
			},
			func(i int) (*model.User, error) {
				return &model.User{}, nil
			},
			http.StatusOK,
		},
		{
			"internal error",
			userInfo{
				-1,
				make(map[string]struct{}),
			},
			func(i int) (*model.User, error) {
				return nil, fmt.Errorf("user doesn't exist")
			},
			http.StatusInternalServerError,
		},
		{
			"ok",
			&struct {
				id string
			}{
				"1",
			},
			func(i int) (*model.User, error) {
				return nil, fmt.Errorf("users' info is invalid")
			},
			http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			service := &service.Service{
				User: &mock.UserServiceMock{
					MockGetUser: test.mock,
				},
			}

			handler := Handler{service}

			req, err := http.NewRequest("GET", "/user", http.NoBody)

			if err != nil {
				t.Fatal(err)
			}

			ctx := context.WithValue(req.Context(), KeyUserInfo, test.userInf)

			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()

			reg := http.HandlerFunc(handler.getUser)

			reg.ServeHTTP(rr, req)

			// Check the status code is what we expect.
			if status := rr.Code; status != test.expectedStatusCode {
				t.Errorf("%s:handler returned wrong status code: got %v want %v",
					test.name, status, test.expectedStatusCode)
			}
		})
	}
}
