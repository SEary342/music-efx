package metadata

import (
	"fmt"
	"path/filepath"
	"strings"

	"music-efx/internal/files"
	"music-efx/pkg/model"
)

func ExtractMetadata(file string) (model.MP3Metadata, error) {
	decodedMp3, err := files.LoadFile(file)
	if err != nil {
		return model.MP3Metadata{}, err
	}
	return model.MP3Metadata{
		Name:   strings.TrimSuffix(filepath.Base(file), filepath.Ext(file)),
		Length: formatLength(decodedMp3.Length()),
		Path:   file,
	}, nil
}

func formatLength(length int64) string {
	minutes := int(length) / 60
	seconds := int(length) % 60
	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}
