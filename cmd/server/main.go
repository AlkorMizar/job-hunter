package main

import (
	"log"

	"github.com/AlkorMizar/job-hunter/pkg"
)

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.apikey ApiKeyAuth
//@in header
//@name access_token

func main() {
	server := pkg.NewServer("localhost:5080")
	if err := server.Run(); err != nil {
		log.Fatalf("error ocured during run %s", err.Error())
	}
}
