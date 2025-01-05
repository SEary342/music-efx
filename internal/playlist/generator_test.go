package playlist

import (
	"music-efx/pkg/model"
	"testing"
	"time"
)

func TestGenerate(t *testing.T) {
	sampleItems := []model.MP3Metadata{
		{Name: "Song1", Length: 30 * time.Second, Path: "/path/to/song1"},
		{Name: "Song2", Length: 60 * time.Second, Path: "/path/to/song2"},
		{Name: "Song3", Length: 120 * time.Second, Path: "/path/to/song3"},
	}

	t.Run("Generates random playlist without duration limit", func(t *testing.T) {
		randomized := Generate(sampleItems, 0)
		if len(randomized) != len(sampleItems) {
			t.Errorf("expected playlist length %d, got %d", len(sampleItems), len(randomized))
		}
	})

	t.Run("Generates closest playlist with duration limit", func(t *testing.T) {
		durationLimit := 90 * time.Second
		playlistResult := Generate(sampleItems, durationLimit)

		var totalDuration time.Duration
		for _, item := range playlistResult {
			totalDuration += item.Length
		}

		if totalDuration > durationLimit {
			t.Errorf("expected total duration <= %v, got %v", durationLimit, totalDuration)
		}
	})
}

func TestGenerateClosestPlaylist(t *testing.T) {
	sampleItems := []model.MP3Metadata{
		{Name: "Song1", Length: 30 * time.Second, Path: "/path/to/song1"},
		{Name: "Song2", Length: 60 * time.Second, Path: "/path/to/song2"},
		{Name: "Song3", Length: 120 * time.Second, Path: "/path/to/song3"},
	}

	t.Run("Finds exact match for duration", func(t *testing.T) {
		playlistResult := Generate(sampleItems, 90*time.Second)

		var totalDuration time.Duration
		for _, item := range playlistResult {
			totalDuration += item.Length
		}

		if totalDuration != 90*time.Second {
			t.Errorf("expected total duration 90s, got %v", totalDuration)
		}
	})

	t.Run("Handles zero duration", func(t *testing.T) {
		playlistResult := generateClosestPlaylist(sampleItems, 0)
		if len(playlistResult) != 0 {
			t.Errorf("expected empty playlist, got %d items", len(playlistResult))
		}
	})

	t.Run("Handles playlist shorter than duration limit", func(t *testing.T) {
		playlistResult := generateClosestPlaylist(sampleItems, 300*time.Second)
		if len(playlistResult) != len(sampleItems) {
			t.Errorf("expected playlist with all items, got %d items", len(playlistResult))
		}
	})
}

func TestRandomizePlaylist(t *testing.T) {
	sampleItems := []model.MP3Metadata{
		{Name: "Song1", Length: 30 * time.Second, Path: "/path/to/song1"},
		{Name: "Song2", Length: 60 * time.Second, Path: "/path/to/song2"},
		{Name: "Song3", Length: 120 * time.Second, Path: "/path/to/song3"},
	}

	t.Run("Shuffles playlist", func(t *testing.T) {
		randomized := randomizePlaylist(sampleItems)

		if len(randomized) != len(sampleItems) {
			t.Errorf("expected shuffled playlist length %d, got %d", len(sampleItems), len(randomized))
		}

		if randomized[0] == sampleItems[0] &&
			randomized[1] == sampleItems[1] &&
			randomized[2] == sampleItems[2] {
			t.Error("playlist was not shuffled")
		}
	})
}
