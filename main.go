package main

import (
	"database/sql"
	"errors"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
)

const (
	shortKeyLength = 6
	charset        = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	addr = "http://localhost"
	port = ":8080"
)

type ShortenRequest struct {
	URL string `json:"url"`
}

type ShortenResponse struct {
	ShortURL string `json:"short_url"`
}

type Server struct {
	db *sql.DB
}

func main() {
	dsn := os.Getenv("DATABASE_DSN")
	if dsn == "" {
		log.Fatal("DATABASE_DSN env is not set")
	}

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
	}
	router := gin.Default()

	router.POST("/shorten", server.handleShorten)
	router.GET("/:shortKey", server.handleRedirect)

	log.Printf("Starting server at %s%s", addr, port)
	if err := router.Run(port); err != nil {
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

	baseURL := fmt.Sprintf("%s%s", addr, port)
	shortURL := fmt.Sprintf("%s/%s", baseURL, shortKey)

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
