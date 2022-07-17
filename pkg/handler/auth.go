package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/AlkorMizar/job-hunter/pkg/handler/model"
	"github.com/go-playground/validator"
)

// @Summary      Registration
// @Description  Creates new user if unique login and email
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param   newUser      body     model.NewUser true "Login, email, password"
// @Success      200  {object}  model.JSONResult
// @Failure      404  {object}  model.JSONResult
// @Failure      500  {object}  model.JSONResult
// @Router       /reg [post]
func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	validate := validator.New()

	decoder := json.NewDecoder(r.Body)

	var newUser model.NewUser

	err := decoder.Decode(&newUser)

	if err != nil {
		log.Print(err)
		writeErrResp(w, "incorrect data format", http.StatusBadRequest)

		return
	}

	err = validate.Struct(newUser)

	if err != nil {
		log.Print(err)
		writeErrResp(w, "incorrect fields", http.StatusBadRequest)

		return
	}

	err = h.services.Authorization.CreateUser(&newUser)

	if err != nil {
		log.Print(err)
		writeErrResp(w, "internal error", http.StatusInternalServerError)

		return
	}

	body := model.JSONResult{
		Message: "Successfully authorized",
		Data:    nil,
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
<<<<<<< HEAD
	_, _ = io.WriteString(w, `user created`)
=======

	if err := json.NewEncoder(w).Encode(body); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
>>>>>>> 64981498c71ed2b0e99f0cbf40d29b02705ad0ea
}

// @Summary      Authentication
// @Description  If user exists returns JWT token
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param        authInfo   body     model.AuthInfo true "Email and password"
// @Success      200  {object}  model.JSONResult{data=model.Token} "Message and token"
// @Failure      404  {object}  model.JSONResult
// @Failure      500  {object}  model.JSONResult
// @Router       /auth [post]
func (h *Handler) authenticate(w http.ResponseWriter, r *http.Request) {
	validate := validator.New()

	decoder := json.NewDecoder(r.Body)

	var authInfo model.AuthInfo

	err := decoder.Decode(&authInfo)

	if err != nil {
		log.Print(err)
		writeErrResp(w, "incorrect data format", http.StatusBadRequest)

		return
	}

	err = validate.Struct(authInfo)

	if err != nil {
		log.Print(err)
		writeErrResp(w, "incorrect fields", http.StatusBadRequest)

		return
	}

	token, err := h.services.Authorization.CreateToken(authInfo)

	if err != nil {
		log.Print(err)
		writeErrResp(w, "internal error", http.StatusInternalServerError)

		return
	}

	body := model.JSONResult{
		Message: "Successfully authorized",
		Data:    model.Token{Token: token},
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(body); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
