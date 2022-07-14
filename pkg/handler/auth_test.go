package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AlkorMizar/job-hunter/pkg/handler/model"
	"github.com/AlkorMizar/job-hunter/pkg/service"
)

const (
	notUniqueLogin = "notUniqueLogin"
)

func TestRegisterHandler(t *testing.T) {

	tests := []struct {
		name               string
		newUser            model.NewUser
		expectedStatusCode int
	}{
		{
			"ok",
			model.NewUser{
				Login:     "root",
				Email:     "root@root.com",
				Password:  "root1",
				CPassword: "root1",
				Roles:     "role",
			},
			http.StatusOK,
		},
		{
			"incorrect login",
			model.NewUser{
				Login:     "ro",
				Email:     "root@root.com",
				Password:  "root1",
				CPassword: "root1",
				Roles:     "role",
			},
			http.StatusBadRequest,
		},
		{
			"incorrect email",
			model.NewUser{
				Login:     "ro",
				Email:     "root@root",
				Password:  "root1",
				CPassword: "root1",
				Roles:     "role",
			},
			http.StatusBadRequest,
		},
		{
			"incorrect password",
			model.NewUser{
				Login:     "root",
				Email:     "root@root.com",
				Password:  "root",
				CPassword: "root",
				Roles:     "role",
			},
			http.StatusBadRequest,
		},
		{
			"incorrect confirm password",
			model.NewUser{
				Login:     "root",
				Email:     "root@root.com",
				Password:  "root1",
				CPassword: "root2",
				Roles:     "role",
			},
			http.StatusBadRequest,
		},
		{
			"internal error(not unique)",
			model.NewUser{
				Login:     notUniqueLogin,
				Email:     "root@root.com",
				Password:  "root1",
				CPassword: "root1",
				Roles:     "role",
			},
			http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			services := &service.Service{UserManagment: &userManagServiceMock{}}
			handler := Handler{services}
			body, err := json.Marshal(test.newUser)

			if err != nil {
				t.Fatal(err)
			}

			req, err := http.NewRequest("POST", "/unauth/reg", bytes.NewBuffer(body))

			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			reg := http.HandlerFunc(handler.register)

			reg.ServeHTTP(rr, req)

			// Check the status code is what we expect.
			if status := rr.Code; status != test.expectedStatusCode {
				t.Errorf("%s:handler returned wrong status code: got %v want %v",
					test.name, status, test.expectedStatusCode)
			}
		})
	}
}

func TestAuthHandler(t *testing.T) {

	tests := []struct {
		name               string
		authInfo           model.AuthInfo
		mock               func(model.AuthInfo) (string, error)
		expectedStatusCode int
		expectedCookie     string
	}{
		{
			"ok",
			model.AuthInfo{
				Email:    "root@root.com",
				Password: "root1",
			},
			func(ai model.AuthInfo) (string, error) {
				return "token", nil
			},
			http.StatusOK,
			"Token=token; Max-Age=3600000000000; HttpOnly",
		},
		{
			"incorrect data",
			model.AuthInfo{
				Email:    "root@.com",
				Password: "root",
			},
			func(ai model.AuthInfo) (string, error) {
				return "", nil
			},
			http.StatusBadRequest,
			"",
		},
		{
			"user not exists",
			model.AuthInfo{
				Email:    "root@root.com",
				Password: "root1",
			},
			func(ai model.AuthInfo) (string, error) {
				return "", fmt.Errorf("internal error")
			},
			http.StatusInternalServerError,
			"",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			services := &service.Service{UserManagment: &userManagServiceMock{
				mockCreateToken: test.mock,
			}}
			handler := Handler{services}
			body, err := json.Marshal(test.authInfo)

			if err != nil {
				t.Fatal(err)
			}

			req, err := http.NewRequest("POST", "/unauth/auth", bytes.NewBuffer(body))

			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			reg := http.HandlerFunc(handler.authorize)

			reg.ServeHTTP(rr, req)

			status := rr.Code

			if status != test.expectedStatusCode {
				t.Errorf("%s:handler returned wrong status code: got %v want %v",
					test.name, status, test.expectedStatusCode)
			}

			if cookie := rr.Header().Get("Set-Cookie"); cookie != test.expectedCookie {
				t.Errorf("%s:handler returned wrong token : got %v want %v",
					test.name, cookie, test.expectedCookie)
			}
		})
	}
}

type userManagServiceMock struct {
	mockCreateToken func(model.AuthInfo) (string, error)
}

func (s *userManagServiceMock) CreateUser(newUser model.NewUser) error {
	if newUser.Login == notUniqueLogin {
		return fmt.Errorf("Not unique")
	}
	return nil
}

func (s *userManagServiceMock) CreateToken(authInfo model.AuthInfo) (string, error) {
	return s.mockCreateToken(authInfo)
}
