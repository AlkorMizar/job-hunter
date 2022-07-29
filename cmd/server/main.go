package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
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
	defaultLog := logging.NewDefaultLogger(logging.ErrorLevel)

	if err := initConfig(); err != nil {
		defaultLog.Fatal("error during initializaint config", zap.Error(err))
	}

	var repo services.Repository

	flag.Func("db", "server works with mysql/postgres", func(arg string) (err error) {
		repo, err = getRepo(arg, defaultLog)
		if err != nil {
			defaultLog.Error("Couldn't parse db args", zap.Error(err))
		}

		return err
	})

	flag.Parse()

	var log *logging.Logger
	if os.Getenv("APP_ENV") == "production" {
		log = logging.NewZapLogger(logging.ErrorLevel, logging.ErrorLevel, viper.GetString("log.path"))
	} else {
		log = logging.NewZapLogger(logging.DebugLeve, logging.DebugLeve, viper.GetString("log.path"))
	}

	if repo == nil {
		log.Fatal("Databse for server not set")
	}

	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		c := make(chan os.Signal, 1) // we need to reserve to buffer size 1, so the notifier are not blocked
		signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

		s := <-c
		log.Info("System received signal to shutdown", zap.Any("signal", s))
		cancel()
	}()

	if err := initConfig(); err != nil {
		log.Fatal("Error initializing configs", zap.Error(err))
	}

	bcryptCost, _ := strconv.Atoi(os.Getenv("BCRYPT_COST"))

	auth := services.NewAuthService(repo, os.Getenv("SIGNING_KEY"), bcryptCost, log)

	router := handler.NewHandler(log, auth)

	srv := server.NewServer(viper.GetString("server.host"), viper.GetString("server.port"), router.InitRoutes(), viper.GetInt("server.timeOutSEc"))

	g, gCtx := errgroup.WithContext(ctx)
	g.Go(func() error {
		log.Info("Running...")
		if err := srv.Run(); err != nil {
			log.Fatal("Error ocured during run", zap.Error(err))
			return fmt.Errorf("in main %w", err)
		}
		return nil
	})

	g.Go(func() error {
		log.Info("Waiting shutdown")
		<-gCtx.Done()
		log.Warn("Shutting down")

		var errAll error

		if err := srv.Shutdown(context.Background()); err != nil {
			log.Error("Error during shutdown server", zap.Error(err))
			errAll = fmt.Errorf("error in shutting down server;")
		}

		if err := repo.Close(); err != nil {
			log.Error("Error during shutdown db", zap.Error(err))
			errAll = fmt.Errorf("error in shutting down database %w", errAll)
		}
		return errAll
	})

	if err := g.Wait(); err != nil {
		log.Error("Error during shutdown", zap.Error(err))
	}
}

func initConfig() error {
	log := logging.NewDefaultLogger(logging.ErrorLevel)

	viper.AddConfigPath("configs")
	viper.SetConfigName("config")

	if err := viper.ReadInConfig(); err != nil {
		return fmt.Errorf("during reading configs %w", err)
	}

	log.Info("Config initialized", zap.Any("configFile", viper.ConfigFileUsed()))

	if err := godotenv.Load(".env"); err != nil {
		return fmt.Errorf("during reading .env %w", err)
	}

	log.Info(".env load")

	return nil
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
		log.Info("Mysql loaded")
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
		log.Info("Postgres loaded")
	}

	return repo, nil
}
