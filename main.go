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

var urlStore = make(map[string]string)

const (
	shortKeyLength = 6
	addr = "http://localhost"
	port = ":8080"
)

func generateShortKey() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	var sb strings.Builder
	for range shortKeyLength {
		sb.WriteByte(charset[seededRand.Intn(len(charset))])
	}
	return sb.String()
}

func handleShorten(c *gin.Context) {
	longURLBytes, err := c.GetRawData()
	if err != nil {
		c.String(http.StatusInternalServerError, "Failed to read request body")
		return
	}

	shortKey := generateShortKey()
	urlStore[shortKey] = string(longURLBytes)

	shortURL := fmt.Sprintf("%s%s/%s", addr, port, shortKey)
	c.String(http.StatusCreated, shortURL)
} 

func handleRedirect(c *gin.Context) {
	shortKey := c.Param("shortKey")

	longURL, ok := urlStore[shortKey]
	if !ok {
		c.String(http.StatusNotFound, "URL not found")
		return
	}

	c.Redirect(http.StatusFound, longURL)
}

func main() {
	router := gin.Default()

	router.POST("/shorten", handleShorten)
	router.GET("/:shortKey", handleRedirect)

	fmt.Printf("Starting server on %s%s\n", addr, port)
	if err := router.Run(port); err != nil {
		log.Fatal("Failed to run server: ", err)
	}
}
