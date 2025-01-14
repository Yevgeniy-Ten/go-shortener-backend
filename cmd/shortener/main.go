package main

import (
	"context"
	"log"
	"shorter/config"
	"shorter/internal/handlers"
	"shorter/internal/logger"
	"shorter/internal/storage"
	"shorter/internal/storage/database"
	"shorter/internal/storage/filestorage"

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
	ctx := context.TODO()
	db, pgxConnectErr := database.New(ctx, myLogger, cfg.DatabaseURL)
	if pgxConnectErr != nil {
		myLogger.Log.Info("Failed to connect to database", zap.Error(pgxConnectErr))
	}
	defer db.Close(ctx)
	fileStorage, err := filestorage.New(cfg.FilePath, myLogger)
	if err != nil {
		myLogger.Log.Info("Failed to create file storage", zap.Error(err))
	}
	defer fileStorage.Close()
	var store *storage.ShortURLStorage
	if db != nil {
		store = storage.New(db.URLRepo)
	} else if fileStorage != nil {
		store = storage.New(fileStorage)
	} else {
		store = storage.New(nil)
	}
	defer fileStorage.Close()
	h := handlers.InitHandlers(cfg.Config, store, myLogger)
	myLogger.Log.Info("Server started", zap.String("address", cfg.Address))
	return h.Run(cfg.Address)
}
