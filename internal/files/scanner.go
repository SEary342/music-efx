package files

import (
	"os"
	"path/filepath"
	"strings"
)

func FindFiles(root string, extension string) ([]string, error) {
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(strings.ToLower(info.Name()), extension) {
			files = append(files, path)
		}
		return nil
	})
	return files, err
}
