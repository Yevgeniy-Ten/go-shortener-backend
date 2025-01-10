package main

import (
	"context"
	"log"
	"net/http"
	"shorter/config"
	"shorter/internal/handlers"
	"shorter/internal/logger"
	"shorter/internal/storage"
	"shorter/internal/storage/database"
	"shorter/internal/storage/filestorage"

	"github.com/gin-gonic/gin"
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
	db, pgxConnectErr := database.NewDatabase(ctx, cfg.DatabaseURL)
	if pgxConnectErr != nil {
		myLogger.Log.Error("Failed to connect to database", zap.Error(pgxConnectErr))
	}
	defer db.Close(ctx)
	fileStorage, err := filestorage.NewFileStorage(cfg.FilePath)
	if err != nil {
		myLogger.Log.Error("Failed to create file storage", zap.Error(err))
	}
	defer fileStorage.Close()
	var store *storage.ShortURLStorage
	if db != nil {
		store = storage.NewShortURLStorage(db.URLRepo)
	} else if fileStorage != nil {
		store = storage.NewShortURLStorage(fileStorage)
	} else {
		store = storage.NewShortURLStorage(nil)
	}
	defer fileStorage.Close()
	h := handlers.InitHandlers(cfg.Config, store, myLogger)
	h.GET("/ping", func(c *gin.Context) {
		if pgxConnectErr != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		c.Status(http.StatusOK)
	})
	myLogger.Log.Info("Server started", zap.String("address", cfg.Address))
	return h.Run(cfg.Address)
}
