package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

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

	// http server
	errChan := make(chan error)

	router := api.NewRouter(1, g, rep)

	srv := http.Server{
		Addr:    fmt.Sprintf("%s:%s", c.Host, c.Port),
		Handler: router,
	}

	log.Info("Start listening to HTTP",
		zap.String("host", c.Host),
		zap.String("port", c.Port),
	)

	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			errChan <- err
		}
	}()

	// Graceful Shutdown
	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errChan:
		return err
	case <-stopChan:
		log.Info("The application is being Shutdown...")
		srvShutdownChan := make(chan struct{})

		var wg sync.WaitGroup

		wg.Go(func() {
			err := srv.Shutdown(ctx)
			if err != nil {
				errChan <- err
			}
		})
		wg.Go(func() {
			pool.Close()
		})

		go func() {
			wg.Wait()
			close(srvShutdownChan)
		}()

		timer := time.NewTimer(time.Second * 5)

		for {
			select {
			case <-timer.C:
				return errors.New("The timer has ended")
			case err := <-errChan:
				if !errors.Is(err, http.ErrServerClosed) {
					return err
				}
			case <-srvShutdownChan:
				log.Info("The application is completed")
				return nil
			}
		}
	}
}
