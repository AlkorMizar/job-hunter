package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/AlkorMizar/job-hunter/pkg/handler/model"
	"github.com/go-playground/validator"
)

// @Summary      Registration
// @Description  creates new user if unique login and email
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param   newUser      body     model.NewUser true "Login, email, password"
// @Success      200  {string}  string
// @Failure      400  {string}  string
// @Failure      500  {string}  string
// @Router       /unauth/reg [post]
func (h *Handler) register(w http.ResponseWriter, r *http.Request) {

	validate := validator.New()

	decoder := json.NewDecoder(r.Body)
	var newUser model.NewUser

	err := decoder.Decode(&newUser)

	if err != nil {
		log.Print(err)
		http.Error(w, "incorrect data format", http.StatusBadRequest)
		return
	}

	err = validate.Struct(newUser)

	if err != nil {
		log.Print(err)
		http.Error(w, "incorrect fields", http.StatusBadRequest)
		return
	}

	err = h.services.UserManagment.CreateUser(newUser)

	if err != nil {
		log.Print(err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	io.WriteString(w, `user created`)
}

// @Summary      Authorization
// @Description  if user exists sets cookie with JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param   authInfo   body     model.AuthInfo true "Email and password"
// @Success      200  {array}   string
// @Failure      400  {object}  string
// @Failure      404  {object}  string
// @Failure      500  {object}  string
// @Router       /unauth/auth [post]
func (h *Handler) authorize(w http.ResponseWriter, r *http.Request) {
	validate := validator.New()
	decoder := json.NewDecoder(r.Body)
	var authInfo model.AuthInfo

	err := decoder.Decode(&authInfo)

	if err != nil {
		http.Error(w, "incorrect data format", http.StatusBadRequest)
		return
	}

	err = validate.Struct(authInfo)

	if err != nil {
		http.Error(w, "incorrect fields", http.StatusBadRequest)
		return
	}

	token, err := h.services.UserManagment.CreateToken(authInfo)

	if err != nil {
		log.Print(err)
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

	tokenCookie := &http.Cookie{
		Name:     "Token",
		Value:    token,
		HttpOnly: true,
		MaxAge:   int(1 * time.Hour),
	}

	w.WriteHeader(http.StatusOK)
	http.SetCookie(w, tokenCookie)
	io.WriteString(w, `user created`)
}

// @Summary      Log out
// @Description  log out user, clear token
// @Tags         auth
// @Security ApiKeyAuth
// @Accept       json
// @Produce      json
// @Success 200 {object} string
// @Failure      400  {object}  string
// @Failure      404  {object}  string
// @Failure      500  {object}  string
// @Router       /auth/out [post]
func (h *Handler) logOut(w http.ResponseWriter, r *http.Request) {

}
