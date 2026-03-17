package testdata

import (
	"context"
	"path/filepath"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
)

// Creates postgres container
func NewPostgresTestcontainer(ctx context.Context) (*postgres.PostgresContainer, error) {
	curDir := GetCurDirPath()

	pgContainer, err := postgres.Run(ctx, "postgres:18-alpine",
		postgres.WithInitScripts(filepath.Join(curDir, "test_init.sql")),
		postgres.WithDatabase("test-db"),
		postgres.WithUsername("postgres"),
		postgres.WithPassword("postgres"),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").WithOccurrence(2),
		),
	)
	if err != nil {
		return nil, err
	}
	return pgContainer, nil
}

// Stores the postgres container, manages snapshot and connections to isolate tests that require it.
type TestSuite struct {
	PgContainer  *postgres.PostgresContainer
	Db           *pgxpool.Pool
	snapshotName string
}

// Creates a TestSuite.
func NewTestSuite(pgContainer *postgres.PostgresContainer, ctx context.Context) (*TestSuite, error) {
	snapshotName := "base-snapshot"
	err := pgContainer.Snapshot(ctx, postgres.WithSnapshotName(snapshotName))
	if err != nil {
		return nil, err
	}

	connStr, err := pgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return nil, err
	}

	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return nil, err
	}

	config.MaxConns = 4
	config.MinConns = 1
	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return nil, err
	}

	return &TestSuite{
		PgContainer:  pgContainer,
		Db:           pool,
		snapshotName: snapshotName,
	}, nil
}

// Completes active connections, applies base snapshot, restoring the database to its original form.
func (s *TestSuite) SetupTestPg(ctx context.Context) error {
	if s.Db != nil {
		s.Db.Close()
	}

	err := s.PgContainer.Restore(ctx, postgres.WithSnapshotName(s.snapshotName))
	if err != nil {
		return err
	}

	connStr, err := s.PgContainer.ConnectionString(ctx, "sslmode=disable")
	if err != nil {
		return err
	}

	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		return err
	}

	config.MaxConns = 1
	config.MinConns = 1

	pool, err := pgxpool.NewWithConfig(ctx, config)
	if err != nil {
		return err
	}

	s.Db = pool

	return nil
}
