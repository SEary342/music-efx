package config

import (
	"fmt"
	"music-efx/pkg/model"
	"os"

	"gopkg.in/yaml.v3"
)

func LoadConfig(configPath string) (*[]model.PlaylistData, error) {
	data, err := os.ReadFile(configPath) // Replace with your file's path
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}

	// Create a custom unmarshaler for the End field
	var playlists []model.PlaylistData
	err = yaml.Unmarshal(data, &playlists)
	if err != nil {
		return nil, fmt.Errorf("failed to parse YAML: %w", err)
	}
	return &playlists, nil
}
