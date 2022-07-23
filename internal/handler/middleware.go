package handler

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/AlkorMizar/job-hunter/internal/services"
)

type ctxKey string

var KeyUserInfo = ctxKey("userInfo")

const (
	authorizationHeaderKey  = "Authorization"
	authorizationTypeBearer = "bearer"
	authTokenFields         = 2
)

func (h *Handler) authentication(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, ok := r.Header[authorizationHeaderKey]
		if !ok || len(token) != 1 {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		fields := strings.Fields(token[0])
		if len(fields) < authTokenFields {
			http.Error(w, "Invalid authorization header format", http.StatusForbidden)
			return
		}

		authorizationType := strings.ToLower(fields[0])
		if authorizationType != authorizationTypeBearer {
			http.Error(w, "Unsupported authorization type "+authorizationType, http.StatusForbidden)
			return
		}

		accessToken := fields[1]

		userInfo, err := h.services.Authorization.ParseToken(accessToken)
		if err != nil {
			if errors.Is(err, services.ErrExpiredToken) {
				http.Error(w, "Forbidden, please authorize again", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Invalid token", http.StatusForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), KeyUserInfo, userInfo)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
