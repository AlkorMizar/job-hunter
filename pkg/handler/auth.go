package handler

import "net/http"

type NewUser struct {
	Login       string `json:"login" binding:"required" minimum:"5" maximum:"40"`
	Email       string `json:"email" binding:"required" maximum:"255"`
	Password    string `json:"password" binding:"required"  minimum:"5" maximum:"40"`
	PasswordRep string `json:"passwordRep" binding:"required"  minimum:"5" maximum:"40"`
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

}

// @Summary      Authorization
// @Description  if user exists sets cookie with JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param   email      body     string     true  "unique email"       Format(email) maxlength(40)
// @Param   password     body     string     true  "password"       minlength(5)  maxlength(40)
// @Success      200  {array}   string
// @Failure      400  {object}  string
// @Failure      404  {object}  string
// @Failure      500  {object}  string
// @Router       /unauth/auth [get]
func (h *Handler) authorize(w http.ResponseWriter, r *http.Request) {

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
// @Router       /auth/out [get]
func (h *Handler) logOut(w http.ResponseWriter, r *http.Request) {

}
