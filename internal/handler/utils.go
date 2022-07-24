package handler

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/AlkorMizar/job-hunter/internal/model/handl"
	"github.com/go-playground/validator"
)

func writeErrResp(w http.ResponseWriter, mess string, status int) {
	body := handl.JSONResult{
		Message: mess,
	}

	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Content-Type-Options", "nosniff")

	if err := json.NewEncoder(w).Encode(body); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func getFromBody(r *http.Request, v interface{}) (err error) {
	validate := validator.New()

	decoder := json.NewDecoder(r.Body)

	err = decoder.Decode(v)

	if err != nil {
		return fmt.Errorf("incorrect body structure")
	}

	err = validate.Struct(v)

	if err != nil {
		return fmt.Errorf("incorrect format")
	}

	return nil
}
