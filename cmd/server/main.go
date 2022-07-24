package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/AlkorMizar/job-hunter/internal/handler"
	"github.com/AlkorMizar/job-hunter/internal/repository"
	"github.com/AlkorMizar/job-hunter/internal/server"
	"github.com/AlkorMizar/job-hunter/internal/services"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"golang.org/x/sync/errgroup"
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
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		c := make(chan os.Signal, 1) // we need to reserve to buffer size 1, so the notifier are not blocked
		signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

		s := <-c
		log.Printf("Signal %v", s)
		cancel()
	}()

	if err := initConfig(); err != nil {
		log.Fatalf("error initializing configs: %s", err.Error())
	}

	if err := godotenv.Load(".env"); err != nil {
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

	auth := services.NewAuthService(repo, os.Getenv("SIGNING_KEY"))

	router := handler.NewHandler(auth)

	srv := server.NewServer(viper.GetString("server.host"), viper.GetString("server.port"), router.InitRoutes(), viper.GetInt("server.timeOutSEc"))

	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		log.Print("Run...")
		if err := srv.Run(); err != nil {
			log.Fatalf("error ocured during run %s", err.Error())
			return fmt.Errorf("in main %w", err)
		}
		return nil
	})
	g.Go(func() error {
		log.Print("Wait end")
		<-gCtx.Done()
		log.Print("Shuting down")
		err := db.Close()
		if err != nil {
			return fmt.Errorf("in main during shutdown db %w", err)
		}
		err = srv.Shutdown(context.Background())
		if err != nil {
			return fmt.Errorf("in main during shutdown server %w", err)
		}
		return nil
	})

	if err := g.Wait(); err != nil {
		fmt.Printf("exit reason: %s \n", err)
	}

}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")

	return viper.ReadInConfig()
}
