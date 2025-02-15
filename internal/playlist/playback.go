package playlist

import (
	"context"
	"fmt"
	"music-efx/internal/config"
	"music-efx/internal/metadata"
	"music-efx/internal/player"
	"music-efx/pkg/model"
	"sort"
	"time"
)

// LoadPlaylist now loads metadata for all tracks listed in the YAML file.
func LoadPlaylist(playlistYamlPath string) ([]model.MP3Metadata, error) {
	playlistConfig, err := config.LoadConfig(playlistYamlPath)
	if err != nil {
		return nil, err
	}

	lst := *playlistConfig
	var mp3Meta []model.MP3Metadata
	for _, item := range lst {
		meta, err := metadata.LoadMp3Metadata(item.Path)
		if err != nil {
			return nil, err
		}
		// Append all metadata items from this file.
		mp3Meta = append(mp3Meta, meta...)
	}
	return mp3Meta, nil
}

func GenerateAndPlay(ctx context.Context, items []model.MP3Metadata, duration int, playerModel *player.PlayerModel) *player.PlayerModel {
	// Generate the playlist based on the given duration.
	playlist := Generate(items, duration, 0)
	return playPlaylist(ctx, playlist, playerModel)
}

func RandomPlay(ctx context.Context, playlist []model.MP3Metadata, playerModel *player.PlayerModel) *player.PlayerModel {
	// Randomize the playlist.
	items := randomizePlaylist(playlist)
	return playPlaylist(ctx, items, playerModel)
}

func playPlaylist(ctx context.Context, playlist []model.MP3Metadata, playerModel *player.PlayerModel) *player.PlayerModel {
	if len(playlist) == 0 {
		fmt.Println("Playlist is empty. Nothing to play.")
		return playerModel
	}

	// Create a channel to receive updates
	playerModelChan := make(chan *player.PlayerModel)

	// Start the first track synchronously
	playerModel.TrackPath = playlist[0].Path
	updatedPlayer := playerModel.StartTrack()

	// Start the playback loop asynchronously
	go func() {
		time.Sleep(playlist[0].Length)
		playNextTrack(ctx, 1, playlist, updatedPlayer, playerModelChan)
	}()

	// Listen for updates from the channel
	go func() {
		for updated := range playerModelChan {
			// Update the player model based on feedback
			playerModel = updated
		}
	}()

	return updatedPlayer
}

func playNextTrack(ctx context.Context, i int, playlist []model.MP3Metadata, playerModel *player.PlayerModel, playerModelChan chan *player.PlayerModel) {
	// Exit if we've reached the end of the playlist, playback has been stopped, or context is canceled.
	if i >= len(playlist) || playerModel.Stopped {
		close(playerModelChan) // Close the channel when playback ends
		return
	}
	select {
	case <-ctx.Done():
		close(playerModelChan)
		return
	default:
	}

	// Set and start the current track
	playerModel.TrackPath = playlist[i].Path
	updatedPlayer := playerModel.StartTrack()

	// Send the updated model to the channel
	playerModelChan <- updatedPlayer

	// Schedule the next track
	go func(trackDuration time.Duration, nextPlayer *player.PlayerModel) {
		time.Sleep(trackDuration)
		playNextTrack(ctx, i+1, playlist, nextPlayer, playerModelChan)
	}(playlist[i].Length, updatedPlayer)
}

func StartAutoPlaylist(ctx context.Context, m *player.PlayerModel) *player.PlayerModel {
	playlists := config.LoadPlaylistYaml()
	playlistMeta := make([]model.PlaylistData, len(playlists))
	copy(playlistMeta, playlists)

	// Sort playlists by their end time
	sort.Slice(playlistMeta, func(i, j int) bool {
		return playlistMeta[i].End.Before(playlistMeta[j].End.Time)
	})

	mp3MetaMap := make(map[string][]model.MP3Metadata)
	for _, lst := range playlistMeta {
		mp3Meta, err := metadata.LoadMp3Metadata(lst.Path)
		if err != nil {
			fmt.Println("Failed to load playlist mp3 files:", lst.Name)
			continue
		}
		mp3MetaMap[lst.Name] = mp3Meta
	}

	// Create a channel to receive the updated PlayerModel
	updatedModelCh := make(chan *player.PlayerModel)

	go func() {
		defer close(updatedModelCh) // Close the channel when the goroutine finishes
		for _, lst := range playlistMeta {
			duration := int(time.Until(lst.End.Time).Seconds())
			if duration <= 0 {
				continue
			}

			// Stop the currently playing track, if any
			m.Player.Stop()

			// Generate and play the playlist
			select {
			case <-ctx.Done():
				fmt.Println("Context canceled, stopping auto-playlist.")
				m.Player.Stop()
				return
			default:
				// Call GenerateAndPlay and update PlayerModel
				m = GenerateAndPlay(ctx, mp3MetaMap[lst.Name], duration, m)
				// Send the updated model back through the channel
				updatedModelCh <- m
			}

			// Wait for the duration or context cancellation
			select {
			case <-ctx.Done():
				fmt.Println("Context canceled during wait.")
				m.Player.Stop()
				return
			case <-time.After(time.Duration(duration) * time.Second):
				// Continue to the next playlist after the duration
			}
		}
	}()

	// Listen for updates to the PlayerModel in the main goroutine
	for updatedModel := range updatedModelCh {
		return updatedModel
	}
	return m
}
