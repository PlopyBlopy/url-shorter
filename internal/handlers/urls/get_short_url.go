package urls

import (
	"context"
	"errors"
	"net/http"

	"github.com/PlopyBlopy/url-shorter/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func GetShortUrlHandler(u func(string, context.Context) (string, error)) func(*gin.Context) {
	return func(c *gin.Context) {
		origUrl := c.DefaultQuery("orig", "")

		validate := validator.New()

		err := validate.Var(origUrl, "required,url")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		shortUrl, err := u(origUrl, c)
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
