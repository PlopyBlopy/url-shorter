package api

import (
	"fmt"

	"github.com/PlopyBlopy/url-shorter/internal/adapters"
	"github.com/PlopyBlopy/url-shorter/internal/domain"
	. "github.com/PlopyBlopy/url-shorter/internal/handlers/urls" // Point-to-point import in order not to specify the package urls
	"github.com/gin-gonic/gin"
)

func NewRouter(version int, g domain.ShortURLGenerator, rep adapters.RepositoryBase) *gin.Engine {
	// usecases
	addUrlUsecase := AddUrlUsecase(g, rep, rep)
	getUrlUsecase := GetUrlUsecase(rep)
	getUrlsUsecase := GetUrlsUsecase(rep)
	getOrigUrlUsecase := GetOrigUrlUsecase(rep)
	getShortUrlUsecase := GetShortUrlUsecase(rep)

	// handlers
	r := gin.Default()
	router := r.Group(fmt.Sprintf("/v%d", version))
	router.POST("/", AddUrlHandler(addUrlUsecase))
	router.GET("/url", GetUrlHandler(getUrlUsecase))
	router.GET("/urls", GetUrlsHandler(getUrlsUsecase))
	router.GET("/origurl", GetOrigUrlHandler(getOrigUrlUsecase))
	router.GET("/shorturl", GetShortUrlHandler(getShortUrlUsecase))

	return r
}
