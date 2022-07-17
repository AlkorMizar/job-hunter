package handler

import (
	"encoding/json"
	"net/http"

	"github.com/AlkorMizar/job-hunter/pkg/handler/model"
)

func (h *Handler) getUser(w http.ResponseWriter, r *http.Request) {
	userInf, ok := r.Context().Value(KeyUserInfo).(userInfo)
	if !ok {
		writeErrResp(w, "User info is invalid", http.StatusBadRequest)
	}

	res, err := h.services.GetUser(userInf.id)
	if err != nil {
		writeErrResp(w, "internal error", http.StatusInternalServerError)
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
