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
		randomized := Generate(sampleItems, 0, 5)
		if len(randomized) != len(sampleItems) {
			t.Errorf("expected playlist length %d, got %d", len(sampleItems), len(randomized))
		}
	})

	t.Run("Generates closest playlist with duration limit and crossfade", func(t *testing.T) {
		durationLimit := 90
		crossfade := 5
		playlistResult := Generate(sampleItems, durationLimit, crossfade)

		var totalDuration int
		for i, item := range playlistResult {
			totalDuration += int(item.Length.Seconds())
			if i > 0 {
				totalDuration -= crossfade // Adjust for crossfade
			}
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

	t.Run("Handles crossfade adjustment", func(t *testing.T) {
		crossfade := 5
		maxDuration := 90
		playlistResult := generateClosestPlaylist(sampleItems, maxDuration, crossfade)

		var totalDuration int
		for i, item := range playlistResult {
			totalDuration += int(item.Length.Seconds())
			if i > 0 {
				totalDuration -= crossfade // Adjust for crossfade
			}
		}

		if totalDuration > maxDuration {
			t.Errorf("expected total duration <= %v, got %v", maxDuration, totalDuration)
		}
	})

	t.Run("Handles crossfade exceeding song duration", func(t *testing.T) {
		sampleItems := []model.MP3Metadata{
			{Name: "ShortSong", Length: 3 * time.Second, Path: "/path/to/shortsong"},
			{Name: "LongSong", Length: 120 * time.Second, Path: "/path/to/longsong"},
		}

		crossfade := 5
		playlistResult := generateClosestPlaylist(sampleItems, 120, crossfade)

		// Ensure the crossfade doesn't result in a negative duration
		for i, item := range playlistResult {
			if i > 0 && int(item.Length.Seconds())-crossfade < 0 {
				t.Errorf("invalid crossfade adjustment for item %v", item.Name)
			}
		}
	})

	t.Run("Handles exact match for duration with crossfade", func(t *testing.T) {
		playlistResult := generateClosestPlaylist(sampleItems, 90, 5)

		// Calculate the actual total duration considering crossfade
		var totalDuration int
		for i, item := range playlistResult {
			duration := int(item.Length.Seconds())
			if i > 0 { // Apply crossfade to all but the first item
				duration -= 5
			}
			totalDuration += duration
		}

		// The expected duration should match the limit
		if totalDuration != 85 {
			t.Errorf("expected total duration 85s, got %v", totalDuration)
		}
	})
}

func TestRandomizePlaylist(t *testing.T) {
	sampleItems := []model.MP3Metadata{
		{Name: "Song1", Length: 30 * time.Second, Path: "/path/to/song1"},
		{Name: "Song2", Length: 60 * time.Second, Path: "/path/to/song2"},
		{Name: "Song3", Length: 120 * time.Second, Path: "/path/to/song3"},
		{Name: "Song4", Length: 120 * time.Second, Path: "/path/to/song4"},
	}

	maxRetries := 5 // Number of retries allowed for the shuffle test

	t.Run("Shuffles playlist", func(t *testing.T) {
		var identical bool
		for attempt := 1; attempt <= maxRetries; attempt++ {
			randomized := randomizePlaylist(sampleItems)

			if len(randomized) != len(sampleItems) {
				t.Errorf("expected shuffled playlist length %d, got %d", len(sampleItems), len(randomized))
				return // Critical failure, no need to retry
			}

			// Check if any element has changed position
			identical = true
			for i := range sampleItems {
				if sampleItems[i] != randomized[i] {
					identical = false
					break
				}
			}

			if !identical {
				return // Test passed, no further retries needed
			}

			// Retry if the playlist appears unshuffled
			if attempt < maxRetries {
				t.Logf("Retrying shuffle test (attempt %d of %d)", attempt+1, maxRetries)
			} else {
				t.Error("playlist was not shuffled after maximum retries")
			}
		}
	})

	t.Run("Handles empty playlist", func(t *testing.T) {
		randomized := randomizePlaylist([]model.MP3Metadata{})
		if len(randomized) != 0 {
			t.Errorf("expected empty playlist, got %d items", len(randomized))
		}
	})
}
