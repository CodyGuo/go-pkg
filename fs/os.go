package fs

import (
	"errors"
	"io/fs"
	"os"
)

func MkdirAll(path string, mode os.FileMode) error {
	_, err := os.Stat(path)
	if err != nil {
		if errors.Is(err, fs.ErrNotExist) {
			return os.MkdirAll(path, mode)
		}
		return err
	}
	return nil
}
