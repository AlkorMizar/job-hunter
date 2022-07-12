package main

import (
	"log"
	"os"

	"github.com/AlkorMizar/job-hunter/pkg"
	"github.com/AlkorMizar/job-hunter/pkg/handler"
	"github.com/AlkorMizar/job-hunter/pkg/repository"
	"github.com/AlkorMizar/job-hunter/pkg/service"
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
	db, err := repository.NewMySQLDB(repository.Config{
		Host:     "localhost",
		Port:     "3306",
		Username: "root",
		DBName:   "db",
		Protocol: "tcp",
		Password: os.Getenv("DB_PASSWORD"),
	})
	if err != nil {
		log.Fatalf("failed to initialize db: %s", err.Error())
	}
	repo := repository.NewRepository(db)
	service := service.NewService(repo)
	router := handler.NewHandler(service)
	server := pkg.NewServer("localhost:8080", router.InitRoutes())
	if err := server.Run(); err != nil {
		log.Fatalf("error ocured during run %s", err.Error())
	}
}
