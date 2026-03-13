package main

import (
	"context"
	"fmt"

	"github.com/PlopyBlopy/url-shorter/config"
	"github.com/PlopyBlopy/url-shorter/internal"
	"github.com/PlopyBlopy/url-shorter/internal/adapters"
	. "github.com/PlopyBlopy/url-shorter/internal/handlers/urls" // Point-to-point import in order not to specify the package urls
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"go.uber.org/zap"
)

func main() {
	c, err := config.NewAppConfig()
	if err != nil {
		panic(err)
	}

	var log *zap.Logger
	if ok := c.IsDev; ok {
		log, err = zap.NewDevelopment()
		if err != nil {
			panic(err)
		}
	} else {
		log, err = zap.NewProduction()
		if err != nil {
			panic(err)
		}
	}

	defer log.Sync()

	err = app(c, log)
	if err != nil {
		log.Fatal("app error", zap.Error(err))
	}
}

func app(c *config.AppConfig, log *zap.Logger) error {
	r := gin.Default()

	// PostgreSQL connection
	ctx := context.Background()
	config, err := pgxpool.ParseConfig(c.DBConnString)
	if err != nil {
		return fmt.Errorf("Failed create pgxpool.Config: %w", err)
	}

	config.MaxConns = c.MaxConns
	config.MinConns = c.MinConns

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return fmt.Errorf("Failed create new pgxpool: %w", err)
	}

	// dependencies
	rep := adapters.NewRepository(pool)

	g, err := internal.NewGenerator(rep, ctx)
	if err != nil {
		return err
	}

	// usecases
	addUrlUsecase := urls.AddUrlUsecase(g, rep)
	getUrlsUsecase := urls.GetUrlsUsecase(rep)

	// handlers
	v1 := r.Group("/v1")
	v1.POST("/", urls.AddUrlHandler(addUrlUsecase))
	v1.GET("/urls", urls.GetUrlsHandler(getUrlsUsecase))

	// HTTP Server
	log.Info("Start listening to HTTP",
		zap.String("domain", c.Domain),
		zap.String("port", c.Port),
	)

	if err := r.Run(fmt.Sprintf("%s:%s", c.Domain, c.Port)); err != nil {
		return err
	}

	// Graceful Shutdown
	// in the future

	return nil
}
