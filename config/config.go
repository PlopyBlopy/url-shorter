package config

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type AppConfig struct {
	IsDev bool
	// HTTP
	Domain string `env:"DOMAIN"`
	Port   string `env:"PORT"`
	// DB
	DBConnString string `env:"DBConnString"`
	MinConn      int    `env:"MINCONN"`
	MaxConn      int    `env:"MAXCONN"`
}

func NewAppConfig() (*AppConfig, error) {
	c := &AppConfig{}

	envFlag := flag.String("env", "dev", "Environment: dev|prod")

	if ok := strings.EqualFold(*envFlag, "dev"); ok {
		c.IsDev = true
	}

	root, err := findProjectRoot()
	if err != nil {
		return nil, fmt.Errorf("Couldn't find the project root: %w", err)
	}

	envFile := filepath.Join(root, fmt.Sprintf(".env.%s", *envFlag))

	err = godotenv.Load(envFile)
	if err != nil {
		return nil, fmt.Errorf("Failed to load .env file: %w", err)
	}

	err = env.Parse(c)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse environment parameters into struct: %w", err)
	}

	return c, nil
}

// findProjectRoot searches for the directory containing go.mod, starting from the current folder.
func findProjectRoot() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(dir, "go.mod")); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return "", fmt.Errorf("The go.mod file was not found in the parent directories.")
}
