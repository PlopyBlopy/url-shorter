package urls

import (
	"context"
	"errors"
	"net/http"

	"github.com/PlopyBlopy/url-shorter/internal/domain"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func GetUrlHandler(u func(url string, ctx context.Context) (domain.Url, error)) func(*gin.Context) {
	return func(c *gin.Context) {
		anyurl := c.DefaultQuery("search", "")

		validate := validator.New()

		err := validate.Var(anyurl, "required,url")
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		url, err := u(anyurl, c)
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
