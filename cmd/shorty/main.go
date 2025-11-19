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


	"github.com/gin-gonic/gin"
)

type Server struct {
	storage storage.URLSaver
	cfg *config.Config
	log *slog.Logger
}

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	cfg, err := config.LoadConfig()
	if err != nil {
		log.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}

	db, err := postgres.New(cfg.DB_DSN)
	if err != nil {
		log.Error("Failed to open DB connection", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	server := &Server{
		storage:  db,
		cfg: cfg,
		log: log,
	}
	router := gin.Default()

	saveHandler := save.New(log, db, cfg.BaseURL)

	router.POST("/shorten", saveHandler)
	router.GET("/:shortKey", server.handleRedirect)

	log.Info("Starting server", "url", cfg.BaseURL)
	if err := router.Run(cfg.HTTPServer); err != nil {
		log.Error("Failed to start server", "error", err)
		os.Exit(1)
	}
}

func (s *Server) handleRedirect(c *gin.Context) {
	ctx := c.Request.Context()
	shortKey := c.Param("shortKey")

	longURL, err := s.storage.GetURL(ctx, shortKey)
	if err != nil {
		if errors.Is(err, storage.ErrURLNotFound) {
			s.log.Info("Link not found", "short_key", shortKey)
			c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
			return
		}

		s.log.Error("failed to get original URL from DB", "error", err)
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "failed to get original URL from DB"},
		)
		return
	}

	c.Redirect(http.StatusFound, longURL)
}
