package addUrl

import (
	"context"
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
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, shortUrl)
	}
}
