package main

import (
	"log"

	"github.com/AlkorMizar/job-hunter/pkg"
	"github.com/AlkorMizar/job-hunter/pkg/handlers"
)

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.apikey ApiKeyAuth
//@in header
//@name Set-Cookie

func main() {
	router := handlers.InitRoutes()
	server := pkg.NewServer("localhost:8080", router)
	if err := server.Run(); err != nil {
		log.Fatalf("error ocured during run %s", err.Error())
	}
}
