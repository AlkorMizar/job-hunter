package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AlkorMizar/job-hunter/internal/logging"
	"github.com/AlkorMizar/job-hunter/internal/model/handl"
	"github.com/AlkorMizar/job-hunter/internal/services/mock"
	"github.com/mitchellh/mapstructure"
)

func TestGetUser(t *testing.T) {
	userNeed := handl.User{
		Login: "bob",
	}

	tests := []struct {
		name               string
		userInf            interface{}
		mock               func(int) (*handl.User, error)
		expectedStatusCode int
	}{
		{
			"ok",
			handl.UserInfo{
				ID:    1,
				Roles: make(map[string]struct{}),
			},
			func(i int) (*handl.User, error) {
				return &userNeed, nil
			},
			http.StatusOK,
		},
		{
			"internal error",
			handl.UserInfo{
				ID:    -1,
				Roles: make(map[string]struct{}),
			},
			func(i int) (*handl.User, error) {
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
			func(i int) (*handl.User, error) {
				return nil, fmt.Errorf("users' info is invalid")
			},
			http.StatusBadRequest,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			handler := Handler{
				profile: &mock.UserServiceMock{
					MockGetUser: test.mock,
				},
				log: logging.NewDefaultLogger(logging.DebugLeve),
			}

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

			bodyResp := handl.JSONResult{}

			err = json.NewDecoder(rr.Body).Decode(&bodyResp)
			if err != nil {
				t.Fatal(err)
			}

			if bodyResp.Data == nil {
				return
			}

			user := handl.User{}
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

func nilMockUpdateUser(id int, inf handl.UpdateInfo) error {
	return nil
}

func TestUpdateUser(t *testing.T) {
	uInfo := handl.UserInfo{
		ID:    1,
		Roles: make(map[string]struct{}),
	}

	tests := []struct {
		name               string
		userInf            interface{}
		newInf             handl.UpdateInfo
		mock               func(id int, inf handl.UpdateInfo) error
		expectedStatusCode int
	}{
		{
			"ok, all fields",
			uInfo,
			handl.UpdateInfo{
				Login: "login", Email: "tesd@fsd.com", FullName: "Fluff Puff",
			},
			nilMockUpdateUser,
			http.StatusOK,
		},
		{
			"ok, one field",
			uInfo,
			handl.UpdateInfo{
				Email: "tesd@fsd.com",
			},
			nilMockUpdateUser,
			http.StatusOK,
		},
		{
			"ok, two fields",
			uInfo,
			handl.UpdateInfo{
				Login: "login", FullName: "Fluff Puff",
			},
			nilMockUpdateUser,
			http.StatusOK,
		},
		{
			"bad request body, incorrect all",
			handl.UserInfo{
				ID:    1,
				Roles: make(map[string]struct{}),
			},
			handl.UpdateInfo{
				Login: "lo", Email: "tesdfsd.com", FullName: "f",
			},
			func(id int, inf handl.UpdateInfo) error {
				return fmt.Errorf("bad request body")
			},
			http.StatusBadRequest,
		},
		{
			"bad request body, incorrect login",
			handl.UserInfo{
				ID:    1,
				Roles: make(map[string]struct{}),
			},
			handl.UpdateInfo{
				Login: "l o", Email: "tesd@fsd.com",
			},
			func(id int, inf handl.UpdateInfo) error {
				return fmt.Errorf("bad request body")
			},
			http.StatusBadRequest,
		},
		{
			"bad request body, incorrect email",
			handl.UserInfo{
				ID:    1,
				Roles: make(map[string]struct{}),
			},
			handl.UpdateInfo{
				Email: "tesdfsd.com",
			},
			func(id int, inf handl.UpdateInfo) error {
				return fmt.Errorf("bad request body")
			},
			http.StatusBadRequest,
		},
		{
			"bad request body, incorrect full name",
			handl.UserInfo{
				ID:    1,
				Roles: make(map[string]struct{}),
			},
			handl.UpdateInfo{
				Login: "login", Email: "tesd@fsd.com", FullName: "f",
			},
			func(id int, inf handl.UpdateInfo) error {
				return fmt.Errorf("bad request body")
			},
			http.StatusBadRequest,
		},
		{
			"bad request body, empty",
			handl.UserInfo{
				ID:    1,
				Roles: make(map[string]struct{}),
			},
			handl.UpdateInfo{},
			func(id int, inf handl.UpdateInfo) error {
				return fmt.Errorf("bad request body")
			},
			http.StatusBadRequest,
		},
		{
			"bad request body, spaces",
			handl.UserInfo{
				ID:    1,
				Roles: make(map[string]struct{}),
			},
			handl.UpdateInfo{FullName: "     "},
			func(id int, inf handl.UpdateInfo) error {
				return fmt.Errorf("bad request body")
			},
			http.StatusBadRequest,
		},
		{
			"internal error, incorrect handl.UserInfo",
			handl.UserInfo{
				ID:    -1,
				Roles: make(map[string]struct{}),
			},
			handl.UpdateInfo{
				Login: "login", Email: "tesd@fsd.com", FullName: "Fluff Puff",
			},
			func(id int, inf handl.UpdateInfo) error {
				return fmt.Errorf("bad request handl.UserInfo")
			},
			http.StatusInternalServerError,
		},
		{
			"invalid handl.UserInfo",
			&struct {
				id string
			}{
				"1",
			},
			handl.UpdateInfo{
				Login: "login", Email: "tesd@fsd.com", FullName: "Fluff Puff",
			},
			nilMockUpdateUser,
			http.StatusBadRequest,
		},
		{
			"internal error during service process",
			handl.UserInfo{
				ID:    1,
				Roles: make(map[string]struct{}),
			},
			handl.UpdateInfo{
				Login: "lddddo", Email: "tesdf@sd.com", FullName: "fdddddd",
			},
			func(id int, inf handl.UpdateInfo) error {
				return fmt.Errorf("internal error")
			},
			http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			handler := Handler{
				profile: &mock.UserServiceMock{
					MockUpdateUSer: test.mock,
				},
			}

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
	uInfo := handl.UserInfo{
		ID:    1,
		Roles: make(map[string]struct{}),
	}

	tests := []struct {
		name               string
		userInf            interface{}
		pwds               handl.Passwords
		mock               func(int, handl.Passwords) error
		expectedStatusCode int
	}{
		{
			"ok",
			uInfo,
			handl.Passwords{
				NewPassword:  "test2",
				CurrPassword: "test2",
			},
			func(int, handl.Passwords) error {
				return nil
			},
			http.StatusOK,
		},
		{
			"incorrect format of new password(empty)",
			uInfo,
			handl.Passwords{
				NewPassword:  "     ",
				CurrPassword: "     ",
			},
			func(int, handl.Passwords) error {
				return fmt.Errorf("incorrect new password")
			},
			http.StatusBadRequest,
		},
		{
			"incorrect format of new password(less then min)",
			uInfo,
			handl.Passwords{
				NewPassword:  "1 ",
				CurrPassword: "1 ",
			},
			func(int, handl.Passwords) error {
				return fmt.Errorf("incorrect new password")
			},
			http.StatusBadRequest,
		},
		{
			"new and confirm not the same",
			uInfo,
			handl.Passwords{
				NewPassword:  "test2",
				CurrPassword: "test3",
			},
			func(int, handl.Passwords) error {
				return fmt.Errorf("incorrect new password")
			},
			http.StatusBadRequest,
		},
		{
			"wrong current password",
			uInfo,
			handl.Passwords{
				NewPassword:  "test1",
				CurrPassword: "test1",
			},
			func(int, handl.Passwords) error {
				return fmt.Errorf("internal error")
			},
			http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			handler := Handler{
				profile: &mock.UserServiceMock{
					MockUpdatePwd: test.mock,
				},
			}

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
