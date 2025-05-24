package tools

import (
	"fmt"
	"os"
	"path/filepath"
)

func CleanDirectory(pathname string) error {
	dirContents, err := os.ReadDir(pathname)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("could not read directory contents: %w", err)
	}
	for _, entry := range dirContents {
		filename := entry.Name()
		pathname := filepath.Join(pathname, filename)
		err = os.RemoveAll(pathname)
		if err != nil {
			return fmt.Errorf("could not remove file '%s': %w", pathname, err)
		}
	}
	return nil
}
