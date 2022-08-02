package handler

import (
	"context"
	"net/http"

	"github.com/AlkorMizar/job-hunter/internal/logging"
	"github.com/AlkorMizar/job-hunter/internal/model/handl"
	"github.com/gorilla/mux"
)

type Authorization interface {
	CreateUser(ctx context.Context, newUser *handl.NewUser) error
	CreateToken(ctx context.Context, authInfo handl.AuthInfo) (string, error)
	ParseToken(ctx context.Context, tokenStr string) (handl.UserInfo, error)
}

type Profile interface {
	GetUser(ctx context.Context, ID int) (*handl.User, error)
	UpdateUser(ctx context.Context, ID int, updateInf handl.UpdateInfo) error
	UpdatePassword(ctx context.Context, ID int, pwd handl.Passwords) error
}

type Handler struct {
	log     *logging.Logger
	auth    Authorization
	profile Profile
}

func NewHandler(log *logging.Logger, auth Authorization, profile Profile) *Handler {
	return &Handler{
		log:     log,
		auth:    auth,
		profile: profile,
	}
}

func (h *Handler) InitRoutes() *mux.Router {
	r := mux.NewRouter()

	r.Use(h.logging)

	unauth := r.PathPrefix("/").Subrouter()
	auth := r.PathPrefix("/").Subrouter()

	sh := http.StripPrefix("/swaggerui/", http.FileServer(http.Dir("./swaggerui/")))
	unauth.PathPrefix("/swaggerui/").Handler(sh)

	api := http.StripPrefix("/api/", http.FileServer(http.Dir("./api/")))
	unauth.PathPrefix("/api/").Handler(api)

	unauth.HandleFunc("/reg", h.register).Methods(http.MethodPost)
	unauth.HandleFunc("/auth", h.authenticate).Methods(http.MethodPost)

	auth.Use(h.authentication)

	user := auth.PathPrefix("/user").Subrouter()
	user.HandleFunc("", h.getUser).Methods(http.MethodGet)
	user.HandleFunc("", h.updateUser).Methods(http.MethodPut)
	user.HandleFunc("/password", h.updatePassword).Methods(http.MethodPut)

	h.log.Info("Routes inited")

	return r
}
