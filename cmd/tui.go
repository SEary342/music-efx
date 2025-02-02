package main

import (
	"music-efx/internal/app"

	tea "github.com/charmbracelet/bubbletea"
)

var title string = "Music-EFX"

func main() {
	m := app.NewAppModel(title)
	p := tea.NewProgram(&m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		panic(err)
	}
}
