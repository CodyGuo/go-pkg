package fs

import (
	"os"
	"path/filepath"
)

func Abs(path string) (string, error) {
	return filepath.Abs(path)
}

func Dir(path string) string {
	return filepath.Dir(path)
}

func IsFilePathValid(path string) bool {
	cleanInput := filepath.Clean(path)
	if cleanInput == "." || cleanInput == ".." {
		return false
	}
	return true
}

func FileExists(name string) bool {
	if name == "" {
		return false
	}
	info, err := os.Stat(name)

	return err == nil && !info.IsDir()
}
