package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/rehydrate1/shorty/internal/config"
	"github.com/rehydrate1/shorty/internal/http-server/handlers/url/redirect"
	"github.com/rehydrate1/shorty/internal/http-server/handlers/url/save"
	"github.com/rehydrate1/shorty/internal/storage/postgres"

	"github.com/gin-gonic/gin"
)

func main() {
	// logger init
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// config init
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}

	// storage init
	db, err := postgres.New(cfg.DB_DSN)
	if err != nil {
		log.Error("Failed to open DB connection", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	// router init
	router := gin.Default()

	// handlers
	saveHandler := save.New(log, db, cfg.BaseURL)
	redirectHandler := redirect.New(log, db)

	// routs
	router.POST("/shorten", saveHandler)
	router.GET("/:shortKey", redirectHandler)

	srv := &http.Server{
		Addr:    cfg.HTTPServer,
		Handler: router,
	}

	// server start
	log.Info("Starting server", "url", cfg.BaseURL)
	go func() {
		err := srv.ListenAndServe()
		if err != nil && err != http.ErrServerClosed {
			log.Error("Failed to start server", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	<-quit
	log.Info("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Error("Server forced to shutdown", "error", err)
	}

	log.Info("Server exiting")
}
