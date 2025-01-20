package main

import (
	"fmt"
	"music-efx/internal/menu"
)

var title string = "Music-EFX"

var mainMenuItems = []menu.MenuItem{
	{Title: "Auto-Playlist", Description: "Start an automatic playlist based on the config yaml files"},
	{Title: "Playlist Selection", Description: "Open a pre-configured playlist"},
	{Title: "Folder Navigation", Description: "Find a song/directory to play from the file-system"},
}

func main() {
	//for {
	menuSelection := menu.Menu(mainMenuItems, title, false, false)
	fmt.Println(menuSelection)
	switch menuSelection {
	case "Auto-Playlist":
		fmt.Println("Playlist")
	case "Playlist Selection":
		fmt.Println("Playlist menu")
	case "Folder Navigation":
		fmt.Println("folder nav!")
	}
	//os.Exit(1)
	//}
}
