package urls

import (
	"context"
	"errors"
	"net/http"

	"github.com/PlopyBlopy/url-shorter/internal/domain"
	"github.com/gin-gonic/gin"
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
			if errors.Is(err, domain.ErrURLSNoAdded) {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, gin.H{"url": shortUrl})
	}
}

func AddUrlUsecase(g domain.ShortURLGenerator, urlRep domain.UrlAddGetter, counterRep domain.CounterGetter) func(string, context.Context) (string, error) {
	return func(origUrl string, ctx context.Context) (string, error) {
		shortUrl, err := urlRep.GetShortUrl(origUrl, ctx)
		if err != nil {
			if errors.Is(err, domain.ErrURLSNotFound) {
				counter, err := counterRep.GetCounter(ctx)
				if err != nil {
					return "", err
				}

				shortUrl = g.GenerateShortUrl(counter)

				url := domain.NewUrl(origUrl, shortUrl)

				err = urlRep.AddUrl(url, ctx)
				if err != nil {
					return "", err
				}
			} else {
				return "", err
			}
		}

		return shortUrl, nil
	}
}
