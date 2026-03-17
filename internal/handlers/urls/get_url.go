package urls

import (
	"context"
	"errors"
	"net/http"

	"github.com/PlopyBlopy/url-shorter/internal/domain"
	"github.com/gin-gonic/gin"
)

func GetUrlHandler(u func(url string, ctx context.Context) (domain.Url, error)) func(*gin.Context) {
	return func(c *gin.Context) {
		var anyurl domain.AnyUrl

		err := c.ShouldBindJSON(&anyurl)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		url, err := u(anyurl.Url, c)
		if err != nil {
			if errors.Is(err, domain.ErrURLSNotFound) {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, url)
	}
}

func GetUrlUsecase(rep domain.UrlAnyGetter) func(string, context.Context) (domain.Url, error) {
	return func(anyurl string, ctx context.Context) (domain.Url, error) {
		var url domain.Url

		url, err := rep.GetUrl(anyurl, ctx)
		if err != nil {
			return url, err
		}
		return url, nil
	}
}
