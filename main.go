package main

import (
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strings"
	"time"
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

func handleShorten(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	defer r.Body.Close()
	longURLBytes, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Failed to read request body", http.StatusInternalServerError)
		return
	}

	shortKey := generateShortKey()
	urlStore[shortKey] = string(longURLBytes)

	shortURL := fmt.Sprintf("%s%s/%s", addr, port, shortKey)
	w.WriteHeader(http.StatusCreated)
	w.Write([]byte(shortURL))
} 

func handleRedirect(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	shortKey := r.URL.Path[1:]
	if shortKey == "" || shortKey == "shorten" {
		http.NotFound(w, r)
		return
	}

	longURL, ok := urlStore[shortKey]
	if !ok {
		http.NotFound(w, r)
		return
	}

	http.Redirect(w, r, longURL, http.StatusFound)
}

func main() {
	http.HandleFunc("/shorten", handleShorten)
	http.HandleFunc("/", handleRedirect)

	fmt.Printf("Starting server on %s%s", addr, port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
