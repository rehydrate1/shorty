package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
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
	store map[string]string
}

func main() {
	server := &Server{
		store: make(map[string]string),
	}
	router := gin.Default()

	router.POST("/shorten", server.handleShorten)
	router.GET("/:shortKey", server.handleRedirect)

	log.Printf("Server started at %s%s", addr, port)
	if err := router.Run(port); err != nil {
		log.Fatal("Failed to start server: ", err)
	}
}

func (s *Server) handleShorten(c *gin.Context) {
	var req ShortenRequest
	if err := c.BindJSON(&req); err != nil {
		return
	}

	shortKey := generateShortKey()
	s.store[shortKey] = req.URL

	baseURL := fmt.Sprintf("%s%s", addr, port)
	shortURL := fmt.Sprintf("%s/%s", baseURL, shortKey)

	resp := ShortenResponse{ShortURL: shortURL}
	c.JSON(http.StatusCreated, resp)
}

func (s *Server) handleRedirect(c *gin.Context) {
	shortKey := c.Param("shortKey")

	longURL, ok := s.store[shortKey]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
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
