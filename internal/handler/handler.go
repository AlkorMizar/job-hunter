package handler

import (
	"net/http"

	_ "github.com/AlkorMizar/job-hunter/api/docs" //nolint:blank-imports // for swagger documentation page
	"github.com/AlkorMizar/job-hunter/internal/model/handl"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
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
	unauth.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/swagger/doc.json"),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	)).Methods(http.MethodGet)

	auth.Use(h.authentication)

	return r
}
