package playlist

import (
	"fmt"
	"music-efx/internal/player"
	"music-efx/pkg/model"
	"time"
)

func GenerateAndPlay(items []model.MP3Metadata, duration int) {
	playlist := Generate(items, duration, 0)
	Play(playlist)
}

func RandomPlay(playlist []model.MP3Metadata) {
	items := randomizePlaylist(playlist)
	Play(items)
}

func Play(playlist []model.MP3Metadata) {
	if len(playlist) == 0 {
		fmt.Println("Playlist is empty. Nothing to play.")
		return
	}

	p := &player.Player{}

	for _, item := range playlist {
		// Load the current track
		currentTrack, err := player.LoadTrack(item.Path)
		if err != nil {
			fmt.Printf("Failed to load track %s: %v\n", item.Path, err)
			continue // Skip to the next track
		}
		p.PlayTrack(currentTrack)

		// Wait for the track to finish
		time.Sleep(currentTrack.Length)
	}

	// Stop the player after the last track
	p.Stop()
}
