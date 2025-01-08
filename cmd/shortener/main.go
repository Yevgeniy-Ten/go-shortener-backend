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
	fileStorage, err := storage.NewShortURLStorage(cfg.FilePath)
	if err != nil {
		return err
	}
	defer fileStorage.Close()
	pgxConnect, pgxConnectErr := database.NewConnection(context.TODO(), cfg.DatabaseURL)
	if pgxConnectErr != nil {
		myLogger.ErrorCtx(context.TODO(), "Failed to connect to database", zap.Error(pgxConnectErr))
	}
	defer pgxConnect.Close(context.TODO())
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
