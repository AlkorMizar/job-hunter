package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/AlkorMizar/job-hunter/internal/model/handl"
	"github.com/AlkorMizar/job-hunter/internal/services/mock"
	"github.com/mitchellh/mapstructure"
)

const (
	notUniqueLogin = "notUniqueLogin"
	expectedToken  = "token"
)

func TestRegisterHandler(t *testing.T) {
	tests := []struct {
		name               string
		newUser            handl.NewUser
		expectedStatusCode int
	}{
		{
			"ok",
			handl.NewUser{
				Login:    "root",
				Email:    "root@root.com",
				Password: "root1",
				Roles:    []string{"mod"},
			},
			http.StatusOK,
		},
		{
			"incorrect login",
			handl.NewUser{
				Login:    "ro",
				Email:    "root@root.com",
				Password: "root1",
				Roles:    []string{"mod"},
			},
			http.StatusBadRequest,
		},
		{
			"incorrect email",
			handl.NewUser{
				Login:    "ro",
				Email:    "root@root",
				Password: "root1",
				Roles:    []string{"mod"},
			},
			http.StatusBadRequest,
		},
		{
			"incorrect password",
			handl.NewUser{
				Login:    "root",
				Email:    "root@root.com",
				Password: "root",
				Roles:    []string{"mod"},
			},
			http.StatusBadRequest,
		},
		{
			"internal error(not unique)",
			handl.NewUser{
				Login:    notUniqueLogin,
				Email:    "root@root.com",
				Password: "root1",
				Roles:    []string{"mod"},
			},
			http.StatusInternalServerError,
		},
	}

	handler := Handler{
		auth: &mock.AuthServiceMock{
			MockCreateUser: func(newUser *handl.NewUser) error {
				if newUser.Login == notUniqueLogin {
					return fmt.Errorf("Internal error")
				}
				return nil
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

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
		authInfo           handl.AuthInfo
		mock               func(handl.AuthInfo) (string, error)
		expectedStatusCode int
		expectedToken      string
	}{
		{
			"incorrect data",
			handl.AuthInfo{
				Email:    "root@.com",
				Password: "root",
			},
			func(handl.AuthInfo) (string, error) {
				return "", nil
			},
			http.StatusBadRequest,
			"",
		},
		{
			"ok",
			handl.AuthInfo{
				Email:    "root@root.com",
				Password: "root1",
			},
			func(handl.AuthInfo) (string, error) {
				return expectedToken, nil
			},
			http.StatusOK,
			expectedToken,
		},
		{
			"user not exists",
			handl.AuthInfo{
				Email:    "root@root.com",
				Password: "root1",
			},
			func(handl.AuthInfo) (string, error) {
				return "", fmt.Errorf("internal error")
			},
			http.StatusInternalServerError,
			"",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {

			handler := Handler{
				auth: &mock.AuthServiceMock{
					MockCreateToken: test.mock,
				}}
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

			bodyResp := handl.JSONResult{}

			err = json.NewDecoder(rr.Body).Decode(&bodyResp)
			if err != nil {
				t.Fatal(err)
			}

			if bodyResp.Data == nil {
				return
			}

			token := handl.Token{}
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
