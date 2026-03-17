package urls

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/PlopyBlopy/url-shorter/internal/domain"
	"github.com/gin-gonic/gin"
)

func GetUrlsHandler(u func(int, context.Context) ([]domain.Url, error)) func(*gin.Context) {
	return func(c *gin.Context) {
		str := c.DefaultQuery("limit", "0")
		val, err := strconv.Atoi(str)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		urls, err := u(val, c)
		if err != nil {
			if errors.Is(err, domain.ErrURLSNotFound) {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, urls)
	}
}

func GetUrlsUsecase(rep domain.UrlsGetter) func(int, context.Context) ([]domain.Url, error) {
	return func(limit int, ctx context.Context) ([]domain.Url, error) {
		urls, err := rep.GetUrls(limit, ctx)
		if err != nil {
			return nil, fmt.Errorf("Failed add url: %w", err)
		}

		return urls, nil
	}
}
