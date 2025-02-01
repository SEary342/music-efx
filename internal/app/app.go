package app

import (
	"music-efx/internal/menu"
	"music-efx/internal/player"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type GlobalModel struct {
	CurrentView string
	SharedData  map[string]interface{}
}

type Model struct {
	Global *GlobalModel
	Menu   menu.MenuModel
	Player player.PlayerModel
}

func NewAppModel(menuItems []menu.MenuItem, title string) Model {
	return Model{
		Global: &GlobalModel{CurrentView: "menu", SharedData: make(map[string]interface{})},
		Menu:   menu.New(menuItems, title, false, false),
		Player: player.PlayerModel{},
	}
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd

	switch m.Global.CurrentView {
	case "menu":
		updatedMenu, menuCmd := m.Menu.Update(msg)
		m.Menu = updatedMenu.(menu.MenuModel)
		cmd = menuCmd

		if m.Menu.Exiting {
			return m, tea.Quit
		}
		// TODO this need to be fully implemenated
		if m.Menu.Choice != "" {
			m.Menu.Choice = ""
			m.Global.CurrentView = "player"
			m.Player.TrackPath = "/home/sameary/Code/music-efx/sample/Test.mp3"
			m.Player.StartTrack()
			cmd = tea.Tick(500*time.Millisecond, func(t time.Time) tea.Msg {
				return player.TickMsg(t)
			})
		}

	case "player":
		updatedPlayer, playerCmd := m.Player.Update(msg)
		m.Player = updatedPlayer.(player.PlayerModel)
		if m.Player.Stopped {
			m.Global.CurrentView = "menu"
		}
		cmd = playerCmd
	}
	return m, cmd
}

func (m Model) View() string {
	switch m.Global.CurrentView {
	case "menu":
		return m.Menu.View()
	case "player":
		return m.Player.View()
	default:
		return "Unknown view"
	}
}
