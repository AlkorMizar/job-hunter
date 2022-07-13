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

type userManagServiceMock struct {
}

func (s *userManagServiceMock) CreateUser(newUser model.NewUser) error {
	if newUser.Login == notUniqueLogin {
		return fmt.Errorf("Not unique")
	}
	return nil
}
