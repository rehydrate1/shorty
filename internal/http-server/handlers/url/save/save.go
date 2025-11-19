package save

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rehydrate1/shorty/internal/lib/random"
)

type Request struct {
	URL string `json:"url" binding:"required,url"`
}

type Response struct {
	ShortURL string `json:"short_url"`
}

type URLSaver interface {
	SaveURL(ctx context.Context, alias, urlToSave string) error
}

func New(log *slog.Logger, urlSaver URLSaver, baseURL string) gin.HandlerFunc {
	return func(c *gin.Context) {
		const op = "handlers.url.save.New"

		ctx := c.Request.Context()
		var req Request
		if err := c.BindJSON(&req); err != nil {
			return
		}

		shortKey := random.NewRandomString()

		if err := urlSaver.SaveURL(ctx, shortKey, req.URL); err != nil {
			log.Error("failed to insert short link to DB", "op", op, "error", err)
			c.JSON(
				http.StatusInternalServerError,
				gin.H{"error": "failed to insert short link to DB"},
			)
			return
		}
		log.Info("Short link created", "short_key", shortKey, "original_url", req.URL)

		shortURL := fmt.Sprintf("%s/%s", baseURL, shortKey)

		resp := Response{ShortURL: shortURL}
		c.JSON(http.StatusCreated, resp)
	}
}