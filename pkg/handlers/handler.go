package handlers

import (
	"github.com/AlkorMizar/job-hunter/pkg/services"
	"github.com/gorilla/mux"
)

type Handler struct {
}

func NewHandler(serv *services.Service) *Handler {
	return nil
}

func InitRoutes() *mux.Router {
	r := mux.NewRouter()
	return r
}
