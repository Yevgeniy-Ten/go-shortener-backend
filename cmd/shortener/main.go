package main

import (
	"context"
	"log"
	"net/http"
	"shorter/config"
	"shorter/internal/database"
	"shorter/internal/handlers"
	"shorter/internal/logger"
	"shorter/internal/storage"

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
	db, pgxConnectErr := database.NewConnection(context.TODO(), cfg.DatabaseURL)
	if pgxConnectErr != nil {
		myLogger.ErrorCtx(ctx, "Failed to connect to database", zap.Error(pgxConnectErr))
	}

	defer db.Close(ctx)
	var fileStorage *storage.ShortURLStorage
	if db != nil {
		fileStorage, err = storage.NewShortURLStorage(cfg.FilePath, db)
	} else {
		fileStorage, err = storage.NewShortURLStorage(cfg.FilePath, nil)
	}
	if err != nil {
		return err
	}
	defer fileStorage.Close()
	h := handlers.NewHandler(cfg.Config, fileStorage, myLogger)
	r := h.GetRoutes()
	r.GET("/ping", func(c *gin.Context) {
		if pgxConnectErr != nil {
			c.Status(http.StatusInternalServerError)
			return
		}
		c.Status(http.StatusOK)
	})
	myLogger.InfoCtx(context.TODO(), "Server started", zap.String("address", cfg.Address))
	return r.Run(cfg.Address)
}
