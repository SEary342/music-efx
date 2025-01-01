package player

import (
	"fmt"
	"music-efx/internal/files"
	"time"

	"github.com/ebitengine/oto/v3"
)

var currentPlayer *oto.Player // Track the current player

// PlayMP3 starts playing an MP3 file and sets the global player instance
func PlayMP3(filePath string) {
	// Stop the current player if it's already playing
	if currentPlayer != nil {
		StopPlayback() // Stop the current playback before starting a new one
	}

	// Load the MP3 file
	decodedMp3, err := files.LoadFile(filePath)
	if err != nil {
		panic("mp3.NewDecoder failed: " + err.Error())
	}

	// Prepare an Oto context for playback
	op := &oto.NewContextOptions{
		SampleRate:   44100,
		ChannelCount: 2,
		Format:       oto.FormatSignedInt16LE,
	}
	otoCtx, readyChan, err := oto.NewContext(op)
	if err != nil {
		panic("oto.NewContext failed: " + err.Error())
	}
	<-readyChan

	// Create the player for the decoded MP3
	currentPlayer = otoCtx.NewPlayer(decodedMp3)
	currentPlayer.Play()

	// Keep the program running until the song finishes
	for currentPlayer.IsPlaying() {
		time.Sleep(time.Millisecond)
	}
	err = currentPlayer.Close()
	if err != nil {
		panic("player.Close failed: " + err.Error())
	}

	// Clear the player after playback finishes
	currentPlayer = nil
}

// StopPlayback stops the current playback if a song is playing
func StopPlayback() {
	// Check if there's a player instance to stop
	if currentPlayer != nil {
		err := currentPlayer.Close()
		if err != nil {
			fmt.Println("Error stopping playback:", err)
		} else {
			fmt.Println("Playback stopped.")
		}
		// Clear the player instance
		currentPlayer = nil
	}
}
