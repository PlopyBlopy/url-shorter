package main

import (
	"github.com/PlopyBlopy/url-shorter/internal/handlers"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func main() {
	log, err := zap.NewProduction()
	if err != nil {
		panic(err)
	}
	defer log.Sync()

	err = app(log)
	if err != nil {
		log.Fatal("app error", zap.Error(err))
	}
}

// TODO для использования значений: .env => config нужно params: config config.Config
func app(log *zap.Logger) error {
	r := gin.Default()

	v1 := r.Group("/v1")
	v1.POST("/", handlers.AddUriHandler())

	log.Info("Start listening to HTTP on the local host:8080")
	_ = log.Sync()

	if err := r.Run(":8080"); err != nil {
		return err
	}

	return nil
}
