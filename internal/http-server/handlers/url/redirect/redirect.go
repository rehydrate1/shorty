package redirect

import (
	"context"
	"errors"
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/rehydrate1/shorty/internal/storage"
)

type URLGetter interface {
	GetURL(ctx context.Context, alias string) (string, error)
}

func New(log *slog.Logger, urlGetter URLGetter) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		shortKey := c.Param("shortKey")

		longURL, err := urlGetter.GetURL(ctx, shortKey)
		if err != nil {
			if errors.Is(err, storage.ErrURLNotFound) {
				log.Info("Link not found", "short_key", shortKey)
				c.JSON(http.StatusNotFound, gin.H{"error": "URL not found"})
				return
			}

			log.Error("failed to get original URL from DB", "error", err)
			c.JSON(
				http.StatusInternalServerError,
				gin.H{"error": "failed to get original URL from DB"},
			)
			return
		}

		c.Redirect(http.StatusFound, longURL)
	}
}
