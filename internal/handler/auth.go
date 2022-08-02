package handler

import (
	"net/http"

	"github.com/AlkorMizar/job-hunter/internal/model/handl"
	"go.uber.org/zap"
)

func (h *Handler) register(w http.ResponseWriter, r *http.Request) {
	log := h.log.WithCtx(r.Context())

	var newUser handl.NewUser

	if err := getFromBody(r, &newUser); err != nil {
		writeErrResp(w, err.Error(), http.StatusBadRequest)

		return
	}

	log.Info("Registering user")

	if err := h.auth.CreateUser(r.Context(), &newUser); err != nil {
		log.Error("Error during creating user", zap.Error(err))
		writeErrResp(w, "internal error", http.StatusInternalServerError)

		return
	}

	log.Info("User successfully registered")

	renderJSON(w, nil, "Successfully registered")
}

func (h *Handler) authenticate(w http.ResponseWriter, r *http.Request) {
	log := h.log.WithCtx(r.Context())

	var authInfo handl.AuthInfo

	if err := getFromBody(r, &authInfo); err != nil {
		writeErrResp(w, err.Error(), http.StatusBadRequest)

		return
	}

	token, err := h.auth.CreateToken(r.Context(), authInfo)

	if err != nil {
		log.Error("Error during user signing in", zap.Error(err))
		writeErrResp(w, "internal error", http.StatusInternalServerError)

		return
	}

	renderJSON(w, handl.Token{Token: token}, "Successfully authenticated")
}
