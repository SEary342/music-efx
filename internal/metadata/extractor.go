package metadata

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"music-efx/internal/files"
	"music-efx/internal/player"
	"music-efx/pkg/model"
)

func ExtractMetadata(file string) (model.MP3Metadata, error) {
	track, err := player.LoadTrack(file, true)
	if err != nil {
		return model.MP3Metadata{}, err
	}
	defer track.Close()
	return model.MP3Metadata{
		Name:   strings.TrimSuffix(filepath.Base(file), filepath.Ext(file)),
		Length: track.Length,
		Path:   file,
	}, nil
}

func FormatLength(length time.Duration) string {
	// Convert length to seconds
	seconds := int(length / time.Second)
	minutes := seconds / 60
	seconds = seconds % 60
	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}

func LoadMp3Metadata(directory string) ([]model.MP3Metadata, error) {
	// Discover MP3 files in the specified directory
	paths, err := files.FindFiles(directory, ".mp3")
	if err != nil {
		fmt.Println("Error loading mp3 data:", err)
		return nil, err
	}

	// Extract metadata for the MP3 files
	var metadataList []model.MP3Metadata
	for _, path := range paths {
		meta, err := ExtractMetadata(path)
		if err == nil {
			metadataList = append(metadataList, meta)
		}
	}
	return metadataList, nil
}
