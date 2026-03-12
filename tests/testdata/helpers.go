package testdata

import (
	"fmt"
	"path/filepath"
	"runtime"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
)

// Позволяет получить путь к текущему файлу
func GetCurDirPath() string {
	_, filename, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(filename)
	fmt.Printf("%s", currentDir)
	return currentDir
}

type HTTPConfig struct {
	// HTTP
	Domain string `env:"DOMAIN"`
	Port   string `env:"PORT"`
}

func NewHTTPConfig() (*HTTPConfig, error) {
	c := &HTTPConfig{}

	curdir := GetCurDirPath()
	envFile := filepath.Join(curdir, ".env")

	err := godotenv.Load(envFile)
	if err != nil {
		return nil, fmt.Errorf("Failed to load .env file: %w", err)
	}

	err = env.Parse(c)
	if err != nil {
		return nil, fmt.Errorf("Failed to parse environment parameters into struct: %w", err)
	}

	return c, nil
}
