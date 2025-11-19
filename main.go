package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
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
}

func main() {
	cfg, err := LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	dsn := cfg.DB_DSN
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		log.Fatalf("ERROR: failed to open DB connection: %v", err)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		log.Fatalf("ERROR: Failed to ping DB: %v", err)
	}

	server := &Server{
		db: db,
		cfg: cfg,
	}
	router := gin.Default()

	router.POST("/shorten", server.handleShorten)
	router.GET("/:shortKey", server.handleRedirect)

	log.Printf("Starting server at %s", cfg.BaseURL)
	if err := router.Run(cfg.HTTPServer); err != nil {
		log.Fatalf("ERROR: Failed to start server: %v", err)
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
		log.Printf("ERROR: failed to insert short link to DB: %v", err)
		c.JSON(
			http.StatusInternalServerError,
			gin.H{"error": "failed to insert short link to DB"},
		)
		return
	}

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
			log.Printf("ERROR: URL not found: %v", err)
			c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
			return
		}

		log.Printf("ERROR: failed to get original URL from DB: %v", err)
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
