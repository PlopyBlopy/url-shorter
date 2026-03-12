package urls

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	"github.com/PlopyBlopy/url-shorter/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

func AddUrlHandler(u func(string, context.Context) (string, error)) func(*gin.Context) {
	return func(c *gin.Context) {
		var i domain.OrigUrl

		if err := c.ShouldBindJSON(&i); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		shortUrl, err := u(i.OrigUrl, c)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, shortUrl)
	}
}

func AddUrlUsecase(g domain.ShortURLGenerator, rep domain.UrlAddGetter) func(string, context.Context) (string, error) {
	return func(origUrl string, ctx context.Context) (string, error) {
		shortUrl, err := rep.GetShortUrl(origUrl, ctx)
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				shortUrl = g.GenerateShortUrl()

				url := domain.NewUrl(origUrl, shortUrl)

				err = rep.AddUrl(url, ctx)
				if err != nil {
					return "", fmt.Errorf("Failed add url: %w", err)
				}
			} else {
				return "", err
			}
		}

		return shortUrl, nil
	}
}
