package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ilyakaznacheev/cleanenv"
	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	shortKeyLength = 6
	charset        = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

type Config struct {
	HTTPServer string `env:"HTTP_SERVER_ADDRESS" env-default:"localhost:8080"`
	BaseURL    string `env:"BASE_URL" env-default:"http://localhost:8080"`
	DB_DSN     string `env:"DATABASE_DSN" env-required:"true"`
}

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	ShortURL string `json:"short_url"`
}

type Server struct {
	db  *sql.DB
	cfg *Config
	log *slog.Logger
}

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	cfg, err := LoadConfig()
	if err != nil {
		log.Error("Failed to load configuration", "error", err)
		os.Exit(1)
	}

	dsn := cfg.DB_DSN
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Error("Failed to open DB connection", "error", err)
		os.Exit(1)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Error("Failed to ping DB", "error", err)
		os.Exit(1)
	}

	server := &Server{
		db:  db,
		cfg: cfg,
		log: log,
	}
	router := gin.Default()

	router.POST("/shorten", server.handleShorten)
	router.GET("/:shortKey", server.handleRedirect)

	log.Info("Starting server", "url", cfg.BaseURL)
	if err := router.Run(cfg.HTTPServer); err != nil {
		log.Error("Failed to start server", "error", err)
		os.Exit(1)
	}
}

func (s *Server) handleShorten(c *gin.Context) {
	ctx := c.Request.Context()
	var req ShortenRequest
	if err := c.BindJSON(&req); err != nil {
		return
	}

	shortKey := generateShortKey()

	query := `INSERT INTO links (short_key, original_url) VALUES ($1, $2)`
	if _, err := s.db.ExecContext(ctx, query, shortKey, req.URL); err != nil {
		s.log.Error("failed to insert short link to DB", "error", err)
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "failed to insert short link to DB"},
		)
		return
	}
	s.log.Info("Short link created", "short_key", shortKey, "original_url", req.URL)

	shortURL := fmt.Sprintf("%s/%s", s.cfg.BaseURL, shortKey)

	resp := ShortenResponse{ShortURL: shortURL}
	c.JSON(http.StatusCreated, resp)
}

func (s *Server) handleRedirect(c *gin.Context) {
	ctx := c.Request.Context()
	shortKey := c.Param("shortKey")

	query := `SELECT original_url FROM links WHERE short_key=$1`

	var longURL string
	if err := s.db.QueryRowContext(ctx, query, shortKey).Scan(&longURL); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
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

func generateShortKey() string {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	var sb strings.Builder
	sb.Grow(shortKeyLength)

	for range shortKeyLength {
		sb.WriteByte(charset[seededRand.Intn(len(charset))])
	}

	return sb.String()
}

func LoadConfig() (*Config, error) {
	var cfg Config
	if err := cleanenv.ReadConfig(".env", &cfg); err != nil {
		// TODO: add ReadEnv
		return nil, fmt.Errorf("failed to read config: %w", err)
	}
	return &cfg, nil
}
