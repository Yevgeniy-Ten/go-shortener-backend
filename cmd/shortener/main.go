package main

import (
	"context"
	"log"
	"net/http"
	"shorter/config"
	"shorter/internal/domain"
	"shorter/internal/handlers"
	"shorter/internal/logger"
	"shorter/internal/urlstorage"
	"shorter/internal/urlstorage/database"
	"shorter/internal/urlstorage/filestorage"

	"github.com/gin-gonic/gin"

	_ "net/http/pprof"

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
	db, pgxConnectErr := database.New(ctx, myLogger, cfg.Config.DatabaseURL)
	if pgxConnectErr != nil {
		if cfg.Config.DatabaseURL == "" {
			myLogger.Log.Warn("Database URL is empty")
		} else {
			cfg.Config.DatabaseError = true
			myLogger.Log.Error("Failed to connect to database", zap.Error(pgxConnectErr))
		}
	}
	defer db.Close(ctx)
	fileStorage, err := filestorage.New(cfg.FilePath, myLogger)
	if err != nil {
		myLogger.Log.Info("Failed to create file storage", zap.Error(err))
	}
	defer fileStorage.Close()
	s := domain.Storage{
		User: nil,
		URLS: nil,
	}
	if db != nil {
		s.URLS = urlstorage.New(db.URLRepo)
		s.User = db.UsersRepo
	} else if fileStorage != nil {
		s.URLS = urlstorage.New(fileStorage)
	} else {
		s.URLS = urlstorage.New(nil)
	}
	defer fileStorage.Close()

	h := handlers.InitHandlers(cfg.Config, s, myLogger)
	h.GET("/debug/pprof/*any", gin.WrapH(http.DefaultServeMux))
	myLogger.Log.Info("Server started", zap.String("address", cfg.Address))
	return h.Run(cfg.Address)
}
