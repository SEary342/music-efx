package playlist

import (
	"math/rand"
	"music-efx/pkg/model"
	"slices"
	"sort"
)

func Generate(items []model.MP3Metadata, duration int, crossfadeDuration int) []model.MP3Metadata {
	randomized := randomizePlaylist(items)
	if duration == 0 {
		return randomized
	}

	playlist := generateClosestPlaylist(randomized, duration, crossfadeDuration)
	return playlist
}

func generateClosestPlaylist(items []model.MP3Metadata, targetDuration int, crossfadeDuration int) []model.MP3Metadata {
	n := len(items)

	// Create a DP table
	dp := make([][]int, n+1)
	for i := range dp {
		dp[i] = make([]int, targetDuration+1)
	}

	// Fill the DP table
	for i := 1; i <= n; i++ {
		itemDuration := int(items[i-1].Length.Seconds())

		// Adjust item duration for crossfade if it's not the first item
		effectiveDuration := itemDuration
		if i > 1 {
			effectiveDuration -= crossfadeDuration
		}

		// Ensure effectiveDuration is at least 0
		if effectiveDuration < 0 {
			effectiveDuration = 0
		}

		for j := 0; j <= targetDuration; j++ {
			if effectiveDuration > j {
				dp[i][j] = dp[i-1][j]
			} else {
				dp[i][j] = max(dp[i-1][j], dp[i-1][j-effectiveDuration]+itemDuration)
			}
		}
	}

	// Backtrack to find the selected items
	selected := []model.MP3Metadata{}
	remainingDuration := targetDuration
	for i := n; i > 0 && remainingDuration > 0; i-- {
		if dp[i][remainingDuration] != dp[i-1][remainingDuration] {
			selected = append(selected, items[i-1])
			remainingDuration -= int(items[i-1].Length.Seconds())

			// Adjust remaining duration for crossfade
			if i > 1 {
				remainingDuration += crossfadeDuration
			}
		}
	}

	// If the total selected playlist duration is less than the targetDuration, try adding items
	totalDuration := 0
	for _, item := range selected {
		totalDuration += int(item.Length.Seconds())
	}

	// Add more items if total duration is less than the target duration
	if totalDuration < targetDuration {
		for i := 0; i < n && totalDuration < targetDuration; i++ {
			// Skip items already in the selected list
			if !contains(selected, items[i]) {
				selected = append(selected, items[i])
				totalDuration += int(items[i].Length.Seconds())

				// Adjust remaining duration for crossfade if needed
				if len(selected) > 1 {
					totalDuration += crossfadeDuration
				}
			}
		}
	}

	// If no items are selected, pick at least one item (even if it doesn't meet the target duration)
	if len(selected) == 0 {
		selected = append(selected, items[0]) // Pick the first item
	}

	// Reverse to maintain the original order of selected items
	slices.Reverse(selected)

	return selected
}

// Helper function to check if an item is already in the selected list
func contains(list []model.MP3Metadata, item model.MP3Metadata) bool {
	for _, selectedItem := range list {
		if selectedItem == item {
			return true
		}
	}
	return false
}

func randomizePlaylist(items []model.MP3Metadata) []model.MP3Metadata {
	randomized := append([]model.MP3Metadata(nil), items...)

	// Shuffle the copy
	sort.Slice(randomized, func(i, j int) bool {
		return rand.Intn(2) == 0
	})
	return randomized
}
