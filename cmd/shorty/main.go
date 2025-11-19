package main

import (
	"errors"
	"log/slog"
	"net/http"
	"os"

	"github.com/rehydrate1/shorty/internal/config"
	"github.com/rehydrate1/shorty/internal/storage/postgres"
	"github.com/rehydrate1/shorty/internal/storage"
	"github.com/rehydrate1/shorty/internal/http-server/handlers/url/save"
	"github.com/rehydrate1/shorty/internal/http-server/handlers/url/redirect"


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

	// server start
	log.Info("Starting server", "url", cfg.BaseURL)
	if err := router.Run(cfg.HTTPServer); err != nil {
		log.Error("Failed to start server", "error", err)
		os.Exit(1)
	}
}
