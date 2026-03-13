package urls

import (
	"context"
	"errors"
	"net/http"

	"github.com/PlopyBlopy/url-shorter/internal/domain"
	"github.com/gin-gonic/gin"
)

func GetShortUrlHandler(u func(string, context.Context) (string, error)) func(*gin.Context) {
	return func(c *gin.Context) {
		var i domain.OrigUrl

		if err := c.ShouldBindJSON(&i); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		shortUrl, err := u(i.OrigUrl, c)
		if err != nil {
			if errors.Is(err, domain.ErrURLSNotFound) {
				c.JSON(http.StatusBadRequest, gin.H{"error": err})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"url": shortUrl})
	}
}

func GetShortUrlUsecase(rep domain.UrlShortGetter) func(string, context.Context) (string, error) {
	return func(origUrl string, ctx context.Context) (string, error) {
		shortUrl, err := rep.GetShortUrl(origUrl, ctx)
		if err != nil {
			return "", err
		}

		return shortUrl, err
	}
}
