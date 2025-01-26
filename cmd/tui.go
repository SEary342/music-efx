package main

import (
	"fmt"
	"music-efx/internal/menu"
	"music-efx/internal/player"
	"os"
)

var title string = "Music-EFX"

var mainMenuItems = []menu.MenuItem{
	{Title: "Auto-Playlist", Description: "Start an automatic playlist based on the config yaml files"},
	{Title: "Playlist Selection", Description: "Open a pre-configured playlist"},
	{Title: "Folder Navigation", Description: "Find a song/directory to play from the file-system"},
}

func main() {
	for {
		menuSelection := menu.Menu(mainMenuItems, title, false, false)
		if menuSelection.Exiting {
			os.Exit(0)
		}
		switch menuSelection.Choice {
		case "Auto-Playlist":
			fmt.Println("Playlist")
		case "Playlist Selection":
			fmt.Println("Playlist menu")
		case "Folder Navigation":
			fmt.Println("folder nav!")
		}
		// TODO can we unify these components? Right now they are running as separate programs
		// https://leg100.github.io/en/posts/building-bubbletea-programs/
		// TODO This is not the final implementation:
		player.PlayUI("/home/sameary/Code/music-efx/sample/Test.mp3")
	}
}
