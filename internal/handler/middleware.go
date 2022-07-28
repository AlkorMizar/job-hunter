package handler

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/AlkorMizar/job-hunter/internal/logging"
	"github.com/AlkorMizar/job-hunter/internal/services"
	"github.com/google/uuid"
	"go.uber.org/zap"
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
		log := h.log.WithCtx(r.Context())

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

		userInfo, err := h.auth.ParseToken(r.Context(), accessToken)
		if err != nil {
			if errors.Is(err, services.ErrExpiredToken) {
				http.Error(w, "Forbidden, please authorize again", http.StatusUnauthorized)
				return
			}
			http.Error(w, "Invalid token", http.StatusForbidden)
			return
		}

		ctx := context.WithValue(r.Context(), KeyUserInfo, userInfo)

		log.Debug("Token parsing result", zap.Any("userInfo", userInfo))

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *Handler) logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := logging.WithRqId(r.Context(), uuid.NewString())
		start := time.Now()
		defer func() {
			h.log.WithCtx(ctx).Debug("Request finished in time",
				zap.String("request", r.RequestURI),
				zap.Duration("executionTime", time.Since(start)))
		}()
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
