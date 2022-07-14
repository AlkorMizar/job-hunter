package handler

import (
	"context"
	"net/http"
)

type userInfo struct {
	id    int
	roles map[string]struct{}
}

func (h *Handler) authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, err := r.Cookie("Token")
		if err != nil {
			http.Error(w, "Forbidden", http.StatusForbidden)
		} else {
			id, roles, err := h.services.UserManagment.ParseToken(token.Value)
			if err != nil {
				http.Error(w, "Forbidden, please authorize", http.StatusForbidden)
			}
			ctx := context.WithValue(r.Context(), "userInfo", &userInfo{id, roles})
			next.ServeHTTP(w, r.WithContext(ctx))
			next.ServeHTTP(w, r)
		}
	})
}
