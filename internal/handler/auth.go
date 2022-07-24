package handler

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/AlkorMizar/job-hunter/internal/model/handl"
)

// @Summary      Registration
// @Description  Creates new user if unique login and email
// @Tags         auth
// @Accept       json
// @Produce      json
// @Param   newUser      body     handl.NewUser true "Login, email, password"
// @Success      200  {object}  handl.JSONResult
// @Failure      404  {object}  handl.JSONResult
// @Failure      500  {object}  handl.JSONResult
// @Router       /reg [post]
func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	var newUser handl.NewUser

	if err := getFromBody(r, &newUser); err != nil {
		writeErrResp(w, err.Error(), http.StatusBadRequest)

		return
	}

	if err := h.auth.CreateUser(&newUser); err != nil {
		log.Print(err)
		writeErrResp(w, "internal error", http.StatusInternalServerError)

		return
	}

	body := handl.JSONResult{
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
// @Param        authInfo   body     handl.AuthInfo true "Email and password"
// @Success      200  {object}  handl.JSONResult{data=handl.Token} "Message and token"
// @Failure      404  {object}  handl.JSONResult
// @Failure      500  {object}  handl.JSONResult
// @Router       /auth [post]
func (h *Handler) authenticate(w http.ResponseWriter, r *http.Request) {
	var authInfo handl.AuthInfo

	if err := getFromBody(r, &authInfo); err != nil {
		log.Print(err)
		writeErrResp(w, err.Error(), http.StatusBadRequest)

		return
	}

	token, err := h.auth.CreateToken(authInfo)

	if err != nil {
		log.Print(err)
		writeErrResp(w, "internal error", http.StatusInternalServerError)

		return
	}

	body := handl.JSONResult{
		Message: "Successfully authorized",
		Data:    handl.Token{Token: token},
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(body); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
