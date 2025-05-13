package main

import (
	"action-detector-backend/config"
	"action-detector-backend/handler"
	"action-detector-backend/pkg/httpserver"
	"action-detector-backend/pkg/logger"
	"action-detector-backend/pkg/postgres"
	"action-detector-backend/repository"
	"action-detector-backend/storage"
	"action-detector-backend/usecase"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

// @title ActionDetector API
// @version 1.0
// @description API Server for Application
// @host localhost:4040
// @BasePath
// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization
func main() {
	logger := logger.GetLogger()
	cfg := config.GetConfig()

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable&timezone=Asia/Tashkent",
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresHost,
		cfg.PostgresPort,
		cfg.PostgresDatabase,
	)

	db, err := postgres.New(connStr, logger)
	if err != nil {
		logger.Fatal(err)
	}
	defer db.Close()

	storage, err := storage.NewStorage(cfg)
	if err != nil {
		logger.Fatal(err)
	}

	repo := repository.NewRepository(db, logger)
	usecase := usecase.NewUsecase(usecase.Dependencies{
		Repository: repo,
		Storage:    storage,
		Logger:     logger,
		Config:     cfg,
	})
	handler := handler.NewHandler(usecase, cfg, logger, db)

	srv := new(httpserver.Server)
	go func() {
		if err := srv.Run(cfg.HTTPHost, cfg.HTTPPort, handler.InitRoutes(cfg)); err != nil {
			logger.Fatalf("error occured while running http server: %s", err.Error())
		}
	}()

	logger.Info("App started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logger.Warn("App shutting down")

	if err := srv.Shutdown(context.Background()); err != nil {
		logger.Errorf("error occured on server shutting down: %s", err.Error())
	}
}
