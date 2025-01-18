package playlist

import (
	"fmt"
	"music-efx/pkg/model"
	"time"
)

func GenerateAndPlay(items []model.MP3Metadata, duration int, playbackHandler func(model.MP3Metadata)) {
	playlist := Generate(items, duration, 0)
	Play(playlist, playbackHandler)
}

func RandomPlay(playlist []model.MP3Metadata, playbackHandler func(model.MP3Metadata)) {
	items := randomizePlaylist(playlist)
	Play(items, playbackHandler)
}

func Play(playlist []model.MP3Metadata, playbackHandler func(model.MP3Metadata)) {
	if len(playlist) == 0 {
		fmt.Println("Playlist is empty. Nothing to play.")
		return
	}

	for _, item := range playlist {
		playbackHandler(item)
		time.Sleep(item.Length)
	}
}
