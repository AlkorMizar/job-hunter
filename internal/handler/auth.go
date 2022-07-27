package handler

import (
	"log"
	"net/http"

	"github.com/AlkorMizar/job-hunter/internal/model/handl"
)

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

	renderJSON(w, nil, "Successfully registered")
}

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

	renderJSON(w, handl.Token{Token: token}, "Successfully authorized")
}
