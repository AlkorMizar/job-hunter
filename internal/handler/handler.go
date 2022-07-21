package handler

import (
	"net/http"

	_ "github.com/AlkorMizar/job-hunter/api/docs" //nolint:blank-imports // for swagger documentation page
	"github.com/AlkorMizar/job-hunter/internal/services"
	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

type Handler struct {
	services *services.Service
}

func NewHandler(serv *services.Service) *Handler {
	return &Handler{services: serv}
}

func (h *Handler) InitRoutes() *mux.Router {
	r := mux.NewRouter()
	unauth := r.PathPrefix("/").Subrouter()
	auth := r.PathPrefix("/").Subrouter()

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
