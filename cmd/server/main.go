package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/AlkorMizar/job-hunter/internal/handler"
	"github.com/AlkorMizar/job-hunter/internal/logging"
	"github.com/AlkorMizar/job-hunter/internal/repository/mysql"
	"github.com/AlkorMizar/job-hunter/internal/repository/postgres"
	"github.com/AlkorMizar/job-hunter/internal/server"
	"github.com/AlkorMizar/job-hunter/internal/services"
	"github.com/joho/godotenv"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

func main() {
	logger := logging.NewLogger()

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		c := make(chan os.Signal, 1) // we need to reserve to buffer size 1, so the notifier are not blocked
		signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

		s := <-c
		logger.Info("System recived signal to shutdown", zap.Any("signal", s))
		cancel()
	}()

	if err := initConfig(); err != nil {
		logger.Fatal("Error initializing configs", zap.Error(err))
	}

	logger.Info("Config init", zap.Any("configFile", viper.ConfigFileUsed()))

	if err := godotenv.Load(".env"); err != nil {
		logger.Fatal("Error loading env variables", zap.Error(err))
	}

	logger.Info(".env load")

	var repo services.Repository
	flag.Func("db", "server works with mysql/postgres", func(arg string) (err error) {
		repo, err = getRepo(arg, logger)
		if err != nil {
			logger.Error("Couldn't parse db args", zap.Error(err))
		}

		return err
	})

	flag.Parse()

	if repo == nil {
		logger.Fatal("Databse for server not set")
	}

	auth := services.NewAuthService(repo, os.Getenv("SIGNING_KEY"), logger)

	router := handler.NewHandler(logger, auth)

	srv := server.NewServer(viper.GetString("server.host"), viper.GetString("server.port"), router.InitRoutes(), viper.GetInt("server.timeOutSEc"))

	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		logger.Info("Runing...")
		if err := srv.Run(); err != nil {
			logger.Fatal("Error ocured during run", zap.Error(err))
			return fmt.Errorf("in main %w", err)
		}
		return nil
	})
	g.Go(func() error {
		logger.Info("Waiting shutdown")
		<-gCtx.Done()
		logger.Warn("Shutting down")

		var errAll error
		errDB := repo.Close()
		if errDB != nil {
			logger.Error("Error during shutdown db", zap.Error(errDB))
			errAll = fmt.Errorf("error in shutting down database ;")
		}

		errSrv := srv.Shutdown(context.Background())
		if errSrv != nil {
			logger.Error("Error during shutdown server", zap.Error(errSrv))
			errAll = fmt.Errorf("error in shutting down server %w", errAll)
		}

		return errAll
	})

	if err := g.Wait(); err != nil {
		logger.Error("Error during shutdown", zap.Error(err))
	}

}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")

	return viper.ReadInConfig()
}

func getRepo(dbType string, log *logging.Logger) (repo services.Repository, err error) {
	log.Info("Parsing args for db", zap.String("dbType", dbType))

	if dbType == "mysql" {
		db, err := mysql.NewMySQLDB(&mysql.Config{
			Host:     viper.GetString("mysql.host"),
			Port:     viper.GetString("mysql.port"),
			Username: viper.GetString("mysql.username"),
			DBName:   viper.GetString("mysql.dbname"),
			Protocol: viper.GetString("mysql.protocol"),
			Options:  viper.GetString("mysql.options"),
			Password: os.Getenv("MYSQL_PASSWORD"),
		})

		if err != nil {
			return nil, fmt.Errorf("failed to initialize mysql db: %w", err)
		}

		repo = mysql.NewMysqlRepository(db, log)
	}

	if dbType == "postgres" {
		db, err := postgres.NewPodtgresDB(&postgres.Config{
			Host:     viper.GetString("postgres.host"),
			Port:     viper.GetString("postgres.port"),
			Username: viper.GetString("postgres.username"),
			DBName:   viper.GetString("postgres.dbname"),
			Options:  viper.GetString("postgres.options"),
			Password: os.Getenv("POSTGRES_PASSWORD"),
		})

		if err != nil {
			return nil, fmt.Errorf("failed to initialize postgres db: %w", err)
		}

		repo = postgres.NewPostgresRepository(db, log)
	}

	return repo, nil
}
