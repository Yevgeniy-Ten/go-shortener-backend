// Shortener url service
// Used for short urls
package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	_ "net/http/pprof"
	"shorter/config"
	"shorter/internal/domain"
	"shorter/internal/handlers"
	"shorter/internal/logger"
	"shorter/internal/shutdown"
	"shorter/internal/urlstorage"
	"shorter/internal/urlstorage/database"
	"shorter/internal/urlstorage/filestorage"
	"time"

	"go.uber.org/zap"
)

// TimeForShutdown is the time for shutdown
const TimeForShutdown = 3

var buildVersion = "N/A"
var buildDate = "N/A"
var buildCommit = "N/A"

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
	switch {
	case db != nil:
		s.URLS = urlstorage.New(db.URLRepo)
		s.User = db.UsersRepo
	case fileStorage != nil:
		s.URLS = urlstorage.New(fileStorage)
	default:
		s.URLS = urlstorage.New(nil)
	}

	defer fileStorage.Close()

	h := handlers.InitHandlers(cfg.Config, s, myLogger)
	h.GET("/debug/pprof/*any", gin.WrapH(http.DefaultServeMux))
	myLogger.Log.Info("Server started:", zap.String("address", cfg.Address))
	myLogger.Log.Info("Build version:", zap.String("version", buildVersion))
	myLogger.Log.Info("Build date:", zap.String("date", buildDate))
	myLogger.Log.Info("Build commit:", zap.String("commit", buildCommit))
	server := &http.Server{
		Addr:    cfg.Address,
		Handler: h,
	}
	if cfg.HTTPS {
		server := &http.Server{
			Addr:    cfg.Address,
			Handler: h,
		}
		if err := server.ListenAndServeTLS("./certs/cert.pem", "./certs/key.pem"); err != nil {
			return err
		}
	} else {
		if err := server.ListenAndServe(); err != nil {

			return err
		}
	}
	go func() {
		shutdown.GraceFullShutDown(server, TimeForShutdown*time.Second)
	}()
	return err
}
