package config

import (
	"fmt"
	"music-efx/internal/files"
	"music-efx/pkg/model"
	"os"

	"gopkg.in/yaml.v3"
)

func loadConfig(configPath string) (*[]model.PlaylistData, error) {
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

func LoadPlaylistYaml() []model.PlaylistData {
	// Load YAML Playlist configs
	var playlists []model.PlaylistData

	// Get the current working directory
	osWd, err := os.Getwd()
	if err != nil {
		fmt.Println("Error:", err)
		return playlists
	}

	// Find YAML files in the current directory
	yamlPaths, err := files.FindFiles(osWd, ".yaml")
	if err != nil {
		fmt.Println("Error finding YAML files:", err)
		return playlists
	}

	for _, file := range yamlPaths {
		cfgYaml, err := loadConfig(file)
		if err != nil {
			fmt.Printf("Unable to load playlist from: %s\n", file)
			continue
		}
		playlists = append(playlists, *cfgYaml...)
	}
	return playlists
}
