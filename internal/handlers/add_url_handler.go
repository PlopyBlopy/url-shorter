package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type AddUrlInput struct {
	Url string `json:"url" binding:"required,url"`
}

func AddUriHandler() func(*gin.Context) {
	return func(c *gin.Context) {
		var input AddUrlInput

		if err := c.ShouldBindJSON(&input); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusCreated, input)
	}
}
