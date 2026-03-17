package testdata

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
)

type TestSuiteHTTP struct {
	*TestSuite
	BaseUrl  string
	hostPort string
	ctx      context.Context
	router   *gin.Engine
	server   *http.Server
	Client   *http.Client
}

type CancelFunc func() error

func NewTestSuiteHTTP(host, port string, pgContainer *postgres.PostgresContainer, ctx context.Context) (*TestSuiteHTTP, error) {
	ts, err := NewTestSuite(pgContainer, ctx)
	if err != nil {
		return nil, err
	}

	baseUrl := fmt.Sprintf("http://%s:%s/v1", host, port)

	return &TestSuiteHTTP{
		BaseUrl:   baseUrl,
		hostPort:  fmt.Sprintf("%s:%s", host, port),
		TestSuite: ts,
	}, nil
}

func (s *TestSuiteHTTP) SetupTestHTTP(r *gin.Engine, ctx context.Context) (error, CancelFunc) {
	s.server = &http.Server{
		Addr:    s.hostPort,
		Handler: r,
	}

	s.Client = &http.Client{
		Timeout: time.Second * 2,
	}

	go func() {
		_ = s.server.ListenAndServe()
	}()

	return nil, s.Shutdown(ctx)
}

func (s *TestSuiteHTTP) Shutdown(ctx context.Context) func() error {
	return func() error {
		err := s.server.Shutdown(ctx)
		if err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return nil
			}
			return err
		}

		return nil
	}
}
