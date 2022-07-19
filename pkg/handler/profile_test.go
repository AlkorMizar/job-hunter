package handler

import (
	"bytes"
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
			"invalid userIfo",
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

func nilMockUpdateUser(id int, inf model.UpdateInfo) error {
	return nil
}

func TestUpdateUser(t *testing.T) {
	uInfo := userInfo{
		1,
		make(map[string]struct{}),
	}

	tests := []struct {
		name               string
		userInf            interface{}
		newInf             model.UpdateInfo
		mock               func(id int, inf model.UpdateInfo) error
		expectedStatusCode int
	}{
		{
			"ok, all fields",
			uInfo,
			model.UpdateInfo{
				Login: "login", Email: "tesd@fsd.com", FullName: "Fluff Puff",
			},
			nilMockUpdateUser,
			http.StatusOK,
		},
		{
			"ok, one field",
			uInfo,
			model.UpdateInfo{
				Email: "tesd@fsd.com",
			},
			nilMockUpdateUser,
			http.StatusOK,
		},
		{
			"ok, two fields",
			uInfo,
			model.UpdateInfo{
				Login: "login", FullName: "Fluff Puff",
			},
			nilMockUpdateUser,
			http.StatusOK,
		},
		{
			"bad request body, incorrect all",
			userInfo{
				1,
				make(map[string]struct{}),
			},
			model.UpdateInfo{
				Login: "lo", Email: "tesdfsd.com", FullName: "f",
			},
			func(id int, inf model.UpdateInfo) error {
				return fmt.Errorf("bad request body")
			},
			http.StatusBadRequest,
		},
		{
			"bad request body, incorrect login",
			userInfo{
				1,
				make(map[string]struct{}),
			},
			model.UpdateInfo{
				Login: "l o", Email: "tesd@fsd.com",
			},
			func(id int, inf model.UpdateInfo) error {
				return fmt.Errorf("bad request body")
			},
			http.StatusBadRequest,
		},
		{
			"bad request body, incorrect email",
			userInfo{
				1,
				make(map[string]struct{}),
			},
			model.UpdateInfo{
				Email: "tesdfsd.com",
			},
			func(id int, inf model.UpdateInfo) error {
				return fmt.Errorf("bad request body")
			},
			http.StatusBadRequest,
		},
		{
			"bad request body, incorrect full name",
			userInfo{
				1,
				make(map[string]struct{}),
			},
			model.UpdateInfo{
				Login: "login", Email: "tesd@fsd.com", FullName: "f",
			},
			func(id int, inf model.UpdateInfo) error {
				return fmt.Errorf("bad request body")
			},
			http.StatusBadRequest,
		},
		{
			"bad request body, empty",
			userInfo{
				1,
				make(map[string]struct{}),
			},
			model.UpdateInfo{},
			func(id int, inf model.UpdateInfo) error {
				return fmt.Errorf("bad request body")
			},
			http.StatusBadRequest,
		},
		{
			"bad request body, spaces",
			userInfo{
				1,
				make(map[string]struct{}),
			},
			model.UpdateInfo{FullName: "     "},
			func(id int, inf model.UpdateInfo) error {
				return fmt.Errorf("bad request body")
			},
			http.StatusBadRequest,
		},
		{
			"internal error, incorrect userInfo",
			userInfo{
				-1,
				make(map[string]struct{}),
			},
			model.UpdateInfo{
				Login: "login", Email: "tesd@fsd.com", FullName: "Fluff Puff",
			},
			func(id int, inf model.UpdateInfo) error {
				return fmt.Errorf("bad request userInfo")
			},
			http.StatusInternalServerError,
		},
		{
			"invalid userInfo",
			&struct {
				id string
			}{
				"1",
			},
			model.UpdateInfo{
				Login: "login", Email: "tesd@fsd.com", FullName: "Fluff Puff",
			},
			nilMockUpdateUser,
			http.StatusBadRequest,
		},
		{
			"internal error during service process",
			userInfo{
				1,
				make(map[string]struct{}),
			},
			model.UpdateInfo{
				Login: "lddddo", Email: "tesdf@sd.com", FullName: "fdddddd",
			},
			func(id int, inf model.UpdateInfo) error {
				return fmt.Errorf("internal error")
			},
			http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			service := &service.Service{
				User: &mock.UserServiceMock{
					MockUpdateUSer: test.mock,
				},
			}

			handler := Handler{service}

			body, err := json.Marshal(test.newInf)

			if err != nil {
				t.Fatal(err)
			}

			req, err := http.NewRequest("PUT", "/user", bytes.NewBuffer(body))

			if err != nil {
				t.Fatal(err)
			}

			ctx := context.WithValue(req.Context(), KeyUserInfo, test.userInf)

			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()

			reg := http.HandlerFunc(handler.updateUser)

			reg.ServeHTTP(rr, req)

			// Check the status code is what we expect.
			if status := rr.Code; status != test.expectedStatusCode {
				t.Fatalf("%s:handler returned wrong status code: got %v want %v",
					test.name, status, test.expectedStatusCode)
			}
		})
	}
}

func TestUpdatePassword(t *testing.T) {
	uInfo := userInfo{
		1,
		make(map[string]struct{}),
	}

	tests := []struct {
		name               string
		userInf            interface{}
		pwds               model.Passwords
		mock               func(int, model.Passwords) error
		expectedStatusCode int
	}{
		{
			"ok",
			uInfo,
			model.Passwords{
				CurrPassword: "test",
				NewPassword:  "test2",
				CPassword:    "test2",
			},
			func(int, model.Passwords) error {
				return nil
			},
			http.StatusOK,
		},
		{
			"incorrect format of new password(empty)",
			uInfo,
			model.Passwords{
				CurrPassword: "test",
				NewPassword:  "   ",
				CPassword:    "   ",
			},
			func(int, model.Passwords) error {
				return fmt.Errorf("incorrect new password")
			},
			http.StatusBadRequest,
		},
		{
			"incorrect format of new password(less then min)",
			uInfo,
			model.Passwords{
				CurrPassword: "test",
				NewPassword:  "1 ",
				CPassword:    "1 ",
			},
			func(int, model.Passwords) error {
				return fmt.Errorf("incorrect new password")
			},
			http.StatusBadRequest,
		},
		{
			"new and confirm not the same",
			uInfo,
			model.Passwords{
				CurrPassword: "test",
				NewPassword:  "test1",
				CPassword:    "test2",
			},
			func(int, model.Passwords) error {
				return fmt.Errorf("incorrect new password")
			},
			http.StatusBadRequest,
		},
		{
			"wrong current password",
			uInfo,
			model.Passwords{
				CurrPassword: "wrong",
				NewPassword:  "test1",
				CPassword:    "test2",
			},
			func(int, model.Passwords) error {
				return fmt.Errorf("internal error")
			},
			http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			service := &service.Service{
				User: &mock.UserServiceMock{
					MockUpdatePwd: test.mock,
				},
			}

			handler := Handler{service}

			body, err := json.Marshal(test.pwds)

			if err != nil {
				t.Fatal(err)
			}

			req, err := http.NewRequest("PUT", "/user/passwords", bytes.NewBuffer(body))

			if err != nil {
				t.Fatal(err)
			}

			ctx := context.WithValue(req.Context(), KeyUserInfo, test.userInf)

			req = req.WithContext(ctx)

			rr := httptest.NewRecorder()

			reg := http.HandlerFunc(handler.updatePassword)

			reg.ServeHTTP(rr, req)

			// Check the status code is what we expect.
			if status := rr.Code; status != test.expectedStatusCode {
				t.Errorf("%s:handler returned wrong status code: got %v want %v",
					test.name, status, test.expectedStatusCode)
			}
		})
	}
}
