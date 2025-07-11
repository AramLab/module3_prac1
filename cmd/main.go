package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"github.com/pkg/errors"

	"github.com/AramLab/module3_prac1/internal/api"
	"github.com/AramLab/module3_prac1/internal/config"
	customLogger "github.com/AramLab/module3_prac1/internal/logger"
	"github.com/AramLab/module3_prac1/internal/repo"
	"github.com/AramLab/module3_prac1/internal/service"
	"github.com/AramLab/module3_prac1/pkg/migration"
)

func main() {
	if err := godotenv.Load(config.EnvPath); err != nil {
		log.Fatal("Ошибка загрузки env файла:", err)
	}

	var cfg config.AppConfig
	if err := envconfig.Process("", &cfg); err != nil {
		log.Fatal(errors.Wrap(err, "failed to load configuration"))
	}

	logger, err := customLogger.NewLogger(cfg.LogLevel)
	if err != nil {
		log.Fatal(errors.Wrap(err, "error initializing logger"))
	}

	if err := migration.RunMigrations(cfg.PostgreSQL); err != nil {
		logger.Fatal(errors.Wrap(err, "failed to run migration"))
	}

	repository, err := repo.NewRepository(context.Background(), cfg.PostgreSQL)
	if err != nil {
		logger.Fatal(errors.Wrap(err, "failed to create repository"))
	}

	serviceInstance := service.NewService(repository, logger)

	app := api.NewRouters(&api.Routers{Service: serviceInstance})

	// Запуск HTTP-сервера в отдельной горутине
	go func() {
		logger.Infof("Starting server on %s", cfg.Rest.ListenAddress)
		if err := app.Listen(cfg.Rest.ListenAddress); err != nil {
			log.Fatal(errors.Wrap(err, "failed to start server"))
		}
	}()

	// Ожидание системных сигналов для корректного завершения работы
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt, syscall.SIGTERM)
	<-signalChan

	logger.Info("Shutting down gracefully...")
}
