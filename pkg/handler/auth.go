package handler

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

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

	err = h.services.Authorization.CreateUser(&newUser)

	if err != nil {
		log.Print(err)
		http.Error(w, "internal error", http.StatusInternalServerError)

		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = io.WriteString(w, `user created`)
}

// @Summary      Authorization
// @Description  if user exists sets cookie with JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        authInfo   body     model.AuthInfo true "Email and password"
// @Success      200  {object}  model.JSONResult{data=model.Token} "Message and token"
// @Failure      404  {string}  string
// @Failure      500  {string}  string
// @Router       /unauth/auth [post]
func (h *Handler) authorize(w http.ResponseWriter, r *http.Request) {
	validate := validator.New()

	decoder := json.NewDecoder(r.Body)

	var authInfo model.AuthInfo

	err := decoder.Decode(&authInfo)

	if err != nil {
		log.Print(err)
		http.Error(w, "incorrect data format", http.StatusBadRequest)

		return
	}

	err = validate.Struct(authInfo)

	if err != nil {
		log.Print(err)
		http.Error(w, "incorrect fields", http.StatusBadRequest)

		return
	}

	token, err := h.services.Authorization.CreateToken(authInfo)

	if err != nil {
		log.Print(err)
		http.Error(w, "internal error", http.StatusInternalServerError)

		return
	}

	body := model.JSONResult{
		Message: "Succesfully authorized",
		Data:    model.Token{Token: token},
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(body)
}
