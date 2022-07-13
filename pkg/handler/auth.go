package handler

import (
	"encoding/json"
	"net/http"
)

type NewUser struct {
	Login     string `json:"login" binding:"required" minimum:"5" maximum:"40" validate:"required,min=3,max=40"`
	Email     string `json:"email" binding:"required" maximum:"255" validate:"required,email"`
	Password  string `json:"password" binding:"required"  minimum:"5" maximum:"40" validate:"required,eqfield=CPassword"`
	CPassword string `json:"cPassword" binding:"required"  minimum:"5" maximum:"40" validate:"required,min=5,max=40"`
}

type AuthInfo struct {
	Email    string `json:"email" binding:"required" maximum:"255" validate:"required,email"`
	Password string `json:"password" binding:"required"  minimum:"5" maximum:"40" validate:"required,email,min=5,max=40"`
}

// @Summary      Registration
// @Description  creates new user if unique login and email
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param   newUser      body     handler.NewUser true "Login, email, password"
// @Success      200  {array}   string
// @Failure      400  {object}  string
// @Failure      500  {object}  string
// @Router       /unauth/reg [post]
func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var newUser NewUser

	err := decoder.Decode(&newUser)

	if err != nil {
		http.Error(w, "incorrect data format", http.StatusBadRequest)
		return
	}
}

// @Summary      Authorization
// @Description  if user exists sets cookie with JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param   authInfo   body     handler.AuthInfo true "Email and password"
// @Success      200  {array}   string
// @Failure      400  {object}  string
// @Failure      404  {object}  string
// @Failure      500  {object}  string
// @Router       /unauth/auth [post]
func (h *Handler) authorize(w http.ResponseWriter, r *http.Request) {
	decoder := json.NewDecoder(r.Body)
	var authInfo AuthInfo

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
