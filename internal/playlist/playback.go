package playlist

import (
	"fmt"
	"music-efx/internal/config"
	"music-efx/internal/metadata"
	"music-efx/internal/player"
	"music-efx/pkg/model"
	"time"
)

func LoadPlaylist(playlistYamlPath string) ([]model.MP3Metadata, error) {
	playlist, err := config.LoadConfig(playlistYamlPath)
	if err != nil {
		return nil, err
	}
	lst := *playlist
	mp3Meta, err := metadata.LoadMp3Metadata(lst[0].Path)
	if err != nil {
		return nil, err
	}
	return mp3Meta, nil
}

func GenerateAndPlay(items []model.MP3Metadata, duration int, playerModel player.PlayerModel) {
	playlist := Generate(items, duration, 0)
	Play(playlist, playerModel)
}

func RandomPlay(playlist []model.MP3Metadata, playerModel player.PlayerModel) {
	items := randomizePlaylist(playlist)
	Play(items, playerModel)
}

func Play(playlist []model.MP3Metadata, playerModel player.PlayerModel) {
	if len(playlist) == 0 {
		fmt.Println("Playlist is empty. Nothing to play.")
		return
	}

	for _, item := range playlist {
		playerModel.TrackPath = item.Path
		playerModel.StartTrack()
		time.Sleep(item.Length)
	}
}
