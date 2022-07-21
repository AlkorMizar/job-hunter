package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AlkorMizar/job-hunter/pkg/handler/model"
	"github.com/AlkorMizar/job-hunter/pkg/service"
	"github.com/AlkorMizar/job-hunter/pkg/service/mock"
	"github.com/mitchellh/mapstructure"
)

func TestGetUser(t *testing.T) {
	userNeed := model.User{
		Login: "bob",
	}

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
				return &userNeed, nil
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

			if rr.Body.Len() <= 0 {
				return
			}

			bodyResp := model.JSONResult{}

			err = json.NewDecoder(rr.Body).Decode(&bodyResp)
			if err != nil {
				t.Fatal(err)
			}

			if bodyResp.Data == nil {
				return
			}

			user := model.User{}
			err = mapstructure.Decode(bodyResp.Data, &user)

			if err != nil {
				t.Fatalf("%s:incorrect data format get %v", test.name, bodyResp)
			}

			if user.Login != userNeed.Login {
				t.Fatalf("%s:incorrect data format : got %v want %v", test.name, user, userNeed)
			}
		})
	}
}
