package playlist

import (
	"math/rand"
	"music-efx/pkg/model"
	"sort"
	"time"
)

func Generate(items []model.MP3Metadata, duration time.Duration) []model.MP3Metadata {
	randomized := randomizePlaylist(items)
	if duration == 0 {
		return randomized
	}

	playlist := generateClosestPlaylist(randomized, duration)
	return playlist
}

func generateClosestPlaylist(items []model.MP3Metadata, duration time.Duration) []model.MP3Metadata {
	n := len(items)
	maxDuration := int(duration.Seconds())

	// Create a DP table
	dp := make([][]int, n+1)
	for i := range dp {
		dp[i] = make([]int, maxDuration+1)
	}

	// Fill the DP table
	for i := 1; i <= n; i++ {
		itemDuration := int(items[i-1].Length.Seconds())
		for j := 0; j <= maxDuration; j++ {
			if itemDuration > j {
				dp[i][j] = dp[i-1][j]
			} else {
				dp[i][j] = max(dp[i-1][j], dp[i-1][j-itemDuration]+itemDuration)
			}
		}
	}

	// Backtrack to find the selected items
	selected := []model.MP3Metadata{}
	remainingDuration := maxDuration
	for i := n; i > 0 && remainingDuration > 0; i-- {
		if dp[i][remainingDuration] != dp[i-1][remainingDuration] {
			selected = append(selected, items[i-1])
			remainingDuration -= int(items[i-1].Length.Seconds())
		}
	}

	return selected
}

func randomizePlaylist(items []model.MP3Metadata) []model.MP3Metadata {
	randomized := append([]model.MP3Metadata(nil), items...)

	// Shuffle the copy
	sort.Slice(randomized, func(i, j int) bool {
		return rand.Intn(2) == 0
	})
	return randomized
}
