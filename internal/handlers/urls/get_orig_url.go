package urls

import (
	"context"
	"errors"
	"net/http"

	"github.com/PlopyBlopy/url-shorter/internal/domain"
	"github.com/gin-gonic/gin"
)

func GetOrigUrlHandler(u func(string, context.Context) (string, error)) func(*gin.Context) {
	return func(c *gin.Context) {
		var i domain.ShortUrl

		if err := c.ShouldBindJSON(&i); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		origUrl, err := u(i.ShortUrl, c)
		if err != nil {
			if errors.Is(err, domain.ErrURLSNotFound) {
				c.JSON(http.StatusBadRequest, gin.H{"error": err})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"url": origUrl})
	}
}

func GetOrigUrlUsecase(rep domain.UrlOrigGetter) func(string, context.Context) (string, error) {
	return func(shortUrl string, ctx context.Context) (string, error) {
		origUrl, err := rep.GetOrigUrl(shortUrl, ctx)
		if err != nil {
			return "", err
		}

		return origUrl, err
	}
}
