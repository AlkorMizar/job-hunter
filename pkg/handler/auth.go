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

	w.WriteHeader(http.StatusOK)
	_, _ = io.WriteString(w, `user created`)
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

func writeErrResp(w http.ResponseWriter, mess string, status int) {
	body := model.JSONResult{
		Message: mess,
	}

	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	if err := json.NewEncoder(w).Encode(body); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
