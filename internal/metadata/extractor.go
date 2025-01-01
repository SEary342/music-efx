package metadata

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"music-efx/internal/player"
	"music-efx/pkg/model"
)

func ExtractMetadata(file string) (model.MP3Metadata, error) {
	track, err := player.LoadTrack(file)
	if err != nil {
		return model.MP3Metadata{}, err
	}
	defer track.Close()
	return model.MP3Metadata{
		Name:   strings.TrimSuffix(filepath.Base(file), filepath.Ext(file)),
		Length: formatLength(track.Length),
		Path:   file,
	}, nil
}

func formatLength(length time.Duration) string {
	// Convert length to seconds
	seconds := int(length / time.Second)
	minutes := seconds / 60
	seconds = seconds % 60
	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}
