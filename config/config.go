package config

import (
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

type AppConfig struct {
	IsDev bool
	// HTTP
	Host string `env:"HTTP_HOST"`
	Port string `env:"HTTP_PORT"`
	// DB
	DBConnString string `env:"APPDB_DBConnString"`
	MinConns     int32  `env:"APPDB_MINCONNS"`
	MaxConns     int32  `env:"APPDB_MAXCONNS"`
}

// Заполняет поля структуры AppConfig из файла .env в корне проекта или из уже установленных переменных окружения.
// Если в корне не было найдено файла .env - ошибки не будет.
func NewAppConfig() (*AppConfig, error) {
	c := &AppConfig{}

	envFlag := flag.String("env", "prod", "Environment: dev|prod")

	flag.Parse()

	if ok := strings.EqualFold(*envFlag, "dev"); ok {
		c.IsDev = true
	}

	root, err := findProjectRoot()
	if err != nil {
		if !errors.Is(err, os.ErrNotExist) {
			return nil, fmt.Errorf("Couldn't find the project root: %w", err)
		}
	} else {
		envFile := filepath.Join(root, fmt.Sprintf(".env.%s", *envFlag))

		err = godotenv.Load(envFile)
		if err != nil {
			if !errors.Is(err, fs.ErrNotExist) {
				return nil, fmt.Errorf("Failed to load .env file: %w", err)
			}
		}
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
	return "", err
}
