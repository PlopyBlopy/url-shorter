package testdata

import (
	"fmt"
	"path/filepath"
	"runtime"
)

// Позволяет получить путь к текущему файлу
func GetCurDirPath() string {
	_, filename, _, _ := runtime.Caller(0)
	currentDir := filepath.Dir(filename)
	fmt.Printf("%s", currentDir)
	return currentDir
}
