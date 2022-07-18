package handler

import (
	"encoding/json"
	"net/http"

	"github.com/AlkorMizar/job-hunter/pkg/handler/model"
)

// @Summary      Get User
// @Description  Returns users' login, email, full name
// @Security     ApiKeyAuth
// @Tags         user
// @Produce      json
// @Success      200  {object}  model.JSONResult{data=model.User} "login, email, full name"
// @Failure      404  {object}  model.JSONResult
// @Failure      500  {object}  model.JSONResult
// @Router       /user [get]
func (h *Handler) getUser(w http.ResponseWriter, r *http.Request) {
	userInf, ok := r.Context().Value(KeyUserInfo).(userInfo)
	if !ok {
		writeErrResp(w, "users' info is invalid", http.StatusBadRequest)

		return
	}

	res, err := h.services.GetUser(userInf.id)
	if err != nil {
		writeErrResp(w, "internal error", http.StatusInternalServerError)

		return
	}

	body := model.JSONResult{
		Message: "Successfully authorized",
		Data:    res,
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(body); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}
