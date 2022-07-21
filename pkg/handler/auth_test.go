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
	"github.com/mitchellh/mapstructure"
)

const (
	notUniqueLogin = "notUniqueLogin"
	expectedToken  = "token"
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
				Roles:     []string{"mod"},
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
				Roles:     []string{"mod"},
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
				Roles:     []string{"mod"},
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
				Roles:     []string{"mod"},
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
				Roles:     []string{"mod"},
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
				Roles:     []string{"mod"},
			},
			http.StatusInternalServerError,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			services := &service.Service{Authorization: &userManagServiceMock{}}
			handler := Handler{services}
			body, err := json.Marshal(test.newUser)

			if err != nil {
				t.Fatal(err)
			}

			req, err := http.NewRequest("POST", "/reg", bytes.NewBuffer(body))

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
		expectedToken      string
	}{
		{
			"incorrect data",
			model.AuthInfo{
				Email:    "root@.com",
				Password: "root",
			},
			func(model.AuthInfo) (string, error) {
				return "", nil
			},
			http.StatusBadRequest,
			"",
		},
		{
			"ok",
			model.AuthInfo{
				Email:    "root@root.com",
				Password: "root1",
			},
			func(model.AuthInfo) (string, error) {
				return expectedToken, nil
			},
			http.StatusOK,
			expectedToken,
		},
		{
			"user not exists",
			model.AuthInfo{
				Email:    "root@root.com",
				Password: "root1",
			},
			func(model.AuthInfo) (string, error) {
				return "", fmt.Errorf("internal error")
			},
			http.StatusInternalServerError,
			"",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			services := &service.Service{Authorization: &userManagServiceMock{
				mockCreateToken: test.mock,
			}}
			handler := Handler{services}
			body, err := json.Marshal(test.authInfo)

			if err != nil {
				t.Fatal(err)
			}

			req, err := http.NewRequest("POST", "/auth", bytes.NewBuffer(body))
			if err != nil {
				t.Fatal(err)
			}

			rr := httptest.NewRecorder()
			reg := http.HandlerFunc(handler.authenticate)

			reg.ServeHTTP(rr, req)

			status := rr.Code
			if status != test.expectedStatusCode {
				t.Fatalf("%s:handler returned wrong status code: got %v want %v",
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

			token := model.Token{}
			err = mapstructure.Decode(bodyResp.Data, &token)

			if err != nil {
				t.Fatalf("%s:incorrect data format get %v", test.name, bodyResp)
			}

			if token.Token != expectedToken {
				t.Fatalf("%s:handler returned wrong token: got %v want %v",
					test.name, token, expectedToken)
			}
		})
	}
}

type userManagServiceMock struct {
	mockCreateToken func(model.AuthInfo) (string, error)
}

func (s *userManagServiceMock) CreateUser(newUser *model.NewUser) error {
	if newUser.Login == notUniqueLogin {
		return fmt.Errorf("Not unique")
	}

	return nil
}

func (s *userManagServiceMock) CreateToken(authInfo model.AuthInfo) (string, error) {
	return s.mockCreateToken(authInfo)
}

func (s *userManagServiceMock) ParseToken(tokenStr string) (id int, role map[string]struct{}, err error) {
	return 0, nil, nil
}
