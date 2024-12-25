package main

import (
	"context"
	"log"
	"shorter/config"
	"shorter/internal/handlers"
	"shorter/internal/logger"
	"shorter/internal/storage"

	"go.uber.org/zap"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}

func run() error {
	cfg, err := config.NewConfig()
	if err != nil {
		return err
	}
	myLogger, err := logger.InitLogger()
	if err != nil {
		return err
	}
	fileStorage, err := storage.NewShortURLStorage(cfg.FilePath)
	if err != nil {
		return err
	}
	defer fileStorage.Close()
	h := handlers.NewHandler(cfg.Config, fileStorage, myLogger)
	r := h.GetRoutes()
	myLogger.InfoCtx(context.TODO(), "Server started", zap.String("address", cfg.Address))
	return r.Run(cfg.Address)
}
