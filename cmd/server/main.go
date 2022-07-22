package main

import (
	"log"
	"os"

	"github.com/AlkorMizar/job-hunter/internal/handler"
	"github.com/AlkorMizar/job-hunter/internal/repository"
	"github.com/AlkorMizar/job-hunter/internal/server"
	"github.com/AlkorMizar/job-hunter/internal/services"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @host      localhost:8080
// @BasePath  /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("error initializing configs: %s", err.Error())
	}

	if err := godotenv.Load("../../.env"); err != nil {
		log.Fatalf("error loading env variables: %s", err.Error())
	}

	db, err := repository.NewMySQLDB(&repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		DBName:   viper.GetString("db.dbname"),
		Protocol: viper.GetString("db.protocol"),
		Options:  viper.GetString("db.options"),
		Password: os.Getenv("DB_PASSWORD"),
	})
	if err != nil {
		log.Fatalf("failed to initialize db: %s", err.Error())
	}

	repo := repository.NewRepository(db)

	service := services.NewService(repo, os.Getenv("SIGNING_KEY"))

	router := handler.NewHandler(service)

	srv := server.NewServer(viper.GetString("adr.host"), viper.GetString("adr.port"), router.InitRoutes())
	if err := srv.Run(); err != nil {
		log.Fatalf("error ocured during run %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("../../configs")
	viper.SetConfigName("config")

	return viper.ReadInConfig()
}
