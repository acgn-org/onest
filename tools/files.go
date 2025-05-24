package tools

import (
	"fmt"
	"os"
	"path/filepath"
)

func CleanDirectory(pathname string, onDelete func(filename string) bool) error {
	dirContents, err := os.ReadDir(pathname)
	if err != nil {
		return fmt.Errorf("could not read directory contents: %w", err)
	}
	for _, entry := range dirContents {
		filename := entry.Name()
		if onDelete(filename) {
			pathname := filepath.Join(pathname, filename)
			err = os.RemoveAll(pathname)
			if err != nil {
				return fmt.Errorf("could not remove file '%s': %w", pathname, err)
			}
		}
	}
	return nil
}
