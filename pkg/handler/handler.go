package handler

import (
	"net/http"

	_ "github.com/AlkorMizar/job-hunter/api/docs"
	"github.com/AlkorMizar/job-hunter/pkg/service"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Handler struct {
	services *service.Service
}

func NewHandler(serv *service.Service) *Handler {
	return &Handler{services: serv}
}

func (h *Handler) InitRoutes() *mux.Router {
	r := mux.NewRouter()
	unauth := r.PathPrefix("/unauth").Subrouter()
	auth := r.PathPrefix("/auth").Subrouter()

	unauth.HandleFunc("/reg", h.register).Methods(http.MethodPost)
	unauth.HandleFunc("/auth", h.authorize).Methods(http.MethodPost)
	unauth.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("http://localhost:8080/unauth/swagger/doc.json"), //The url pointing to API definition
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
		httpSwagger.DomID("swagger-ui"),
	)).Methods(http.MethodGet)

	auth.Use(authentication)
	auth.HandleFunc("/out", h.logOut).Methods(http.MethodPost)

	return r
}
