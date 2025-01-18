package playlist

import (
	"fmt"
	"music-efx/internal/player"
	"music-efx/pkg/model"
	"time"
)

func GenerateAndPlay(items []model.MP3Metadata, duration int, stopChan chan bool) {
	playlist := Generate(items, duration, 0)
	Play(playlist, stopChan)
}

func RandomPlay(playlist []model.MP3Metadata, stopChan chan bool) {
	items := randomizePlaylist(playlist)
	Play(items, stopChan)
}

func Play(playlist []model.MP3Metadata, stopChan chan bool) {
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

		// Wait for the track to finish or until a stop signal is received
		select {
		case <-time.After(currentTrack.Length): // Wait for the track duration
			// Continue to the next track
		case <-stopChan: // Stop if the signal is received
			fmt.Println("Playback stopped early.")
			p.Stop() // Ensure the player stops
			return
		}
	}

	// Stop the player after the last track
	p.Stop()
}
