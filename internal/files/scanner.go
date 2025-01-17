package files

import (
	"fmt"
	"music-efx/pkg/model"
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

func CreateDirectoryTree(metadataList []model.MP3Metadata) map[string][]model.MP3Metadata {
	tree := make(map[string][]model.MP3Metadata)

	for _, file := range metadataList {
		// Extract the directory part of the file path
		dir := getDirectoryFromPath(file.Path)
		tree[dir] = append(tree[dir], file)
	}

	return tree
}

// getDirectoryFromPath extracts the directory portion of the file path.
func getDirectoryFromPath(path string) string {
	// Normalize path to the appropriate format
	dir := filepath.Dir(path)
	return dir
}

func GetDirectory() string {
	// Check if a directory argument is provided
	if len(os.Args) > 1 {
		return os.Args[1]
	}

	// If no directory is provided, prompt the user for a directory
	fmt.Println("Enter directory path to scan for MP3 files:")
	var directory string
	_, err := fmt.Scanln(&directory)
	if err != nil || directory == "" {
		fmt.Println("Invalid directory path.")
		os.Exit(1)
	}

	return directory
}
