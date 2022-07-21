package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/AlkorMizar/job-hunter/internal/handler/model"
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
	var newUser model.NewUser

	if err := getFromBody(r, &newUser); err != nil {
		log.Print(err)
		writeErrResp(w, err.Error(), http.StatusBadRequest)

		return
	}

	if err := h.services.Authorization.CreateUser(&newUser); err != nil {
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

	if err := json.NewEncoder(w).Encode(body); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
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
	var authInfo model.AuthInfo

	if err := getFromBody(r, &authInfo); err != nil {
		log.Print(err)
		writeErrResp(w, err.Error(), http.StatusBadRequest)

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
