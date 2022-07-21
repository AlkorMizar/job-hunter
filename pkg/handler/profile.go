package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/AlkorMizar/job-hunter/pkg/handler/model"
	"github.com/go-playground/validator"
)

// @Summary      Get User
// @Description  Returns users' login, email, full name,roles
// @Security     ApiKeyAuth
// @Tags         user
// @Produce      json
// @Success      200  {object}  model.JSONResult{data=model.User} "login, email, full name,roles"
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

// @Summary      Updates User
// @Description  Changes users' login, email, full name
// @Security     ApiKeyAuth
// @Tags         user
// @Accept       json
// @Produce      json
// @Param        newInfo    body 	   model.UpdateInfo true "Login, email, full name"
// @Success      200  		{object}   model.JSONResult
// @Failure      404  		{object}   model.JSONResult
// @Failure      500  		{object}   model.JSONResult
// @Router       /user [put]
func (h *Handler) updateUser(w http.ResponseWriter, r *http.Request) {
	userInf, ok := r.Context().Value(KeyUserInfo).(userInfo)
	if !ok {
		writeErrResp(w, "users' info is invalid", http.StatusBadRequest)

		return
	}

	validate := validator.New()

	decoder := json.NewDecoder(r.Body)

	var update model.UpdateInfo

	err := decoder.Decode(&update)

	if err != nil {
		writeErrResp(w, "incorrect data format", http.StatusBadRequest)

		return
	}

	err = validate.Struct(update)

	if err != nil {
		writeErrResp(w, "incorrect fields", http.StatusBadRequest)

		return
	}

	if update.FullName != "" {
		update.FullName = strings.TrimSpace(update.FullName)
		if update.FullName == "" {
			writeErrResp(w, "full name empty", http.StatusBadRequest)

			return
		}
	}

	err = h.services.UpdateUser(userInf.id, update)
	if err != nil {
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
