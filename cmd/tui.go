package main

import (
	"music-efx/internal/app"
	"music-efx/internal/menu"

	tea "github.com/charmbracelet/bubbletea"
)

var title string = "Music-EFX"

var mainMenuItems = []menu.MenuItem{
	{Title: "Auto-Playlist", Description: "Start an automatic playlist based on the config yaml files"},
	{Title: "Playlist Selection", Description: "Open a pre-configured playlist"},
	{Title: "Folder Navigation", Description: "Find a song/directory to play from the file-system"},
}

func main() {
	m := app.NewAppModel(mainMenuItems, title)
	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
