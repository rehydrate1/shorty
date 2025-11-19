package random

import (
	"math/rand"
	"strings"
	"time"
)

const (
	shortKeyLength = 6
	charset        = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func NewRandomString() string {
	seededRand := rand.New(rand.NewSource(time.Now().UnixNano()))

	var sb strings.Builder
	sb.Grow(shortKeyLength)

	for range shortKeyLength {
		sb.WriteByte(charset[seededRand.Intn(len(charset))])
	}

	return sb.String()
}