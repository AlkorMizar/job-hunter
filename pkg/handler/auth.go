package handler

import (
	"encoding/json"
	"net/http"

	"github.com/AlkorMizar/job-hunter/pkg/handler/model"
)

// @Summary      Registration
// @Description  creates new user if unique login and email
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param   newUser      body     model.NewUser true "Login, email, password"
// @Success      200  {array}   string
// @Failure      400  {object}  string
// @Failure      500  {object}  string
// @Router       /unauth/reg [post]
func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var newUser model.NewUser

	err := decoder.Decode(&newUser)

	if err != nil {
		http.Error(w, "incorrect data format", http.StatusBadRequest)
		return
	}

	err = h.services.UserManagment.CreateUser(newUser)

	if err != nil {
		http.Error(w, "internal error", http.StatusInternalServerError)
		return
	}

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
	decoder := json.NewDecoder(r.Body)
	var authInfo model.AuthInfo

	err := decoder.Decode(&authInfo)

	if err != nil {
		http.Error(w, "incorrect data format", http.StatusBadRequest)
		return
	}
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
