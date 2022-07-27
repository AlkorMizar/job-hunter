package handler

import (
	"net/http"

	"github.com/AlkorMizar/job-hunter/internal/model/handl"
	"github.com/gorilla/mux"
)

type Authorization interface {
	CreateUser(newUser *handl.NewUser) error
	CreateToken(authInfo handl.AuthInfo) (string, error)
	ParseToken(tokenStr string) (handl.UserInfo, error)
}

type Handler struct {
	auth Authorization
}

func NewHandler(auth Authorization) *Handler {
	return &Handler{auth: auth}
}

func (h *Handler) InitRoutes() *mux.Router {
	r := mux.NewRouter()
	unauth := r.PathPrefix("/").Subrouter()
	auth := r.PathPrefix("/").Subrouter()

	sh := http.StripPrefix("/swaggerui/", http.FileServer(http.Dir("./swaggerui/")))
	r.PathPrefix("/swaggerui/").Handler(sh)
	api := http.StripPrefix("/api/", http.FileServer(http.Dir("./api/")))
	r.PathPrefix("/api/").Handler(api)

	unauth.HandleFunc("/reg", h.register).Methods(http.MethodPost)
	unauth.HandleFunc("/auth", h.authenticate).Methods(http.MethodPost)

	auth.Use(h.authentication)

	return r
}
