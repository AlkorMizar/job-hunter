package handler

import (
	"net/http"
	"strings"

	"github.com/AlkorMizar/job-hunter/internal/model/handl"
	"go.uber.org/zap"
)

func (h *Handler) getUser(w http.ResponseWriter, r *http.Request) {
	log := h.log.WithCtx(r.Context())

	userInf, ok := r.Context().Value(KeyUserInfo).(handl.UserInfo)
	if !ok {
		log.Debug("invalid userInfo type", zap.String("func", "getUser"), zap.Any("ctx_val", r.Context().Value(KeyUserInfo)))
		writeErrResp(w, "users' info is invalid", http.StatusBadRequest)

		return
	}

	res, err := h.profile.GetUser(r.Context(), userInf.ID)
	if err != nil {
		log.Error("during getting user profile", zap.Error(err))
		writeErrResp(w, "internal error", http.StatusInternalServerError)

		return
	}
	renderJSON(w, res, "Profile")
}

func (h *Handler) updateUser(w http.ResponseWriter, r *http.Request) {
	log := h.log.WithCtx(r.Context())

	userInf, ok := r.Context().Value(KeyUserInfo).(handl.UserInfo)
	if !ok {
		log.Debug("invalid userInfo type", zap.String("func", "getUser"), zap.Any("ctx_val", r.Context().Value(KeyUserInfo)))
		writeErrResp(w, "users' info is invalid", http.StatusBadRequest)

		return
	}

	var update handl.UpdateInfo

	if err := getFromBody(r, &update); err != nil {
		writeErrResp(w, err.Error(), http.StatusBadRequest)

		return
	}

	if update.FullName != "" {
		update.FullName = strings.TrimSpace(update.FullName)
		if update.FullName == "" {
			writeErrResp(w, "full name empty", http.StatusBadRequest)

			return
		}
	}

	if err := h.profile.UpdateUser(r.Context(), userInf.ID, update); err != nil {
		log.Error("during updating user profile", zap.Error(err))
		writeErrResp(w, "internal error", http.StatusInternalServerError)

		return
	}

	renderJSON(w, nil, "Profile info changed")
}

func (h *Handler) updatePassword(w http.ResponseWriter, r *http.Request) {
	log := h.log.WithCtx(r.Context())

	userInf, ok := r.Context().Value(KeyUserInfo).(handl.UserInfo)
	if !ok {
		log.Debug("invalid userInfo type", zap.String("func", "getUser"), zap.Any("ctx_val", r.Context().Value(KeyUserInfo)))
		writeErrResp(w, "users' info is invalid", http.StatusBadRequest)

		return
	}

	var pwds handl.Passwords

	if err := getFromBody(r, &pwds); err != nil {
		writeErrResp(w, err.Error(), http.StatusBadRequest)

		return
	}

	if err := h.profile.UpdatePassword(r.Context(), userInf.ID, pwds); err != nil {
		log.Error("during updating password", zap.Error(err))
		writeErrResp(w, "internal error", http.StatusInternalServerError)

		return
	}

	renderJSON(w, nil, "Password changed")
}
