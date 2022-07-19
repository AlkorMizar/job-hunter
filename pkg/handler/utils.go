package handler

import (
	"encoding/json"
	"net/http"

	"github.com/AlkorMizar/job-hunter/pkg/handler/model"
)

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
