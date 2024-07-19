package utility

import (
	"os"
)

// RemoveFile removes a file from the filesystem
func RemoveFile(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil
	}

	err := os.RemoveAll(path)
	if err != nil {
		return err
	}

	return nil
}
