package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/PlopyBlopy/url-shorter/config"
	"github.com/PlopyBlopy/url-shorter/internal"
	"github.com/PlopyBlopy/url-shorter/internal/adapters"
	"github.com/PlopyBlopy/url-shorter/internal/api"
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

	err = app(*c, log)
	if err != nil {
		log.Fatal("app error", zap.Error(err))
	}
}

func app(c config.AppConfig, log *zap.Logger) error {

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

	//router
	router := api.NewRouter(1, g, rep)

	//httpserver
	log.Info("Start listening to HTTP",
		zap.String("domain", c.Domain),
		zap.String("port", c.Port),
	)

	srv := http.Server{
		Addr:    fmt.Sprintf("%s:%s", c.Domain, c.Port),
		Handler: router,
	}

	// горутина для ListenAndServe
	// Ожидание os события с завершением приложения
	// Shutdown для srv с context.WithTimeout
	err = srv.ListenAndServe()
	if err != nil {
		return err
	}

	// Graceful Shutdown
	// in the future

	return nil
}
