package handler

import (
	"context"
	"net/http"
)

type userInfo struct {
	id    int
	roles map[string]struct{}
}

type ctxKey string

var keyUserInfo = ctxKey("userInfo")

func (h *Handler) authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := r.Cookie("Token")
		if err != nil {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		} else {
			id, roles, err := h.services.Authorization.ParseToken(token.Value)
			if err != nil {
				http.Error(w, "Forbidden, please authorize", http.StatusForbidden)
				return
			}
			ctx := context.WithValue(r.Context(), keyUserInfo, &userInfo{id, roles})
			next.ServeHTTP(w, r.WithContext(ctx))
		}
	})
}
