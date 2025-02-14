package app

import (
	"music-efx/internal/files"
	"music-efx/internal/menu"
	"music-efx/internal/player"
	"os"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

var mainMenuItems = []menu.MenuItem{
	{Title: "Auto-Playlist", Description: "Start an automatic playlist based on the config yaml files"},
	{Title: "Playlist Selection", Description: "Open a pre-configured playlist"},
	{Title: "Folder Navigation", Description: "Find a song/directory to play from the file-system"},
}

type GlobalModel struct {
	CurrentView string
	SharedData  map[string]interface{}
}

type Model struct {
	Global      *GlobalModel
	Menu        menu.MenuModel
	Player      player.PlayerModel
	TrackPicker files.FileModel
}

func NewAppModel(title string) Model {
	pwd, _ := os.Getwd()

	return Model{
		Global:      &GlobalModel{CurrentView: "menu", SharedData: make(map[string]interface{})},
		Menu:        menu.New(mainMenuItems, title, false, false),
		Player:      player.PlayerModel{},
		TrackPicker: files.InitFilePicker(".mp3", pwd, false),
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
		if m.Menu.Choice == "Folder Navigation" {
			m.Menu.Choice = ""
			m.Global.CurrentView = "track-picker"
			updatedPicker, tpCmd := m.TrackPicker.Update(m.TrackPicker.Init()())
			m.TrackPicker = updatedPicker.(files.FileModel)
			cmd = tpCmd
		} else if m.Menu.Choice != "" {
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

	case "track-picker":
		updatedPicker, tpCmd := m.TrackPicker.Update(msg)
		m.TrackPicker = updatedPicker.(files.FileModel)
		if m.TrackPicker.Quitting {
			m.TrackPicker.Quitting = false
			m.Global.CurrentView = "menu"
		} else if len(m.TrackPicker.SelectedFile) > 0 {
			m.Global.CurrentView = "player"
			m.Player.TrackPath = m.TrackPicker.SelectedFile
			m.TrackPicker.SelectedFile = ""
			m.Player.StartTrack()
			tpCmd = tea.Tick(500*time.Millisecond, func(t time.Time) tea.Msg {
				return player.TickMsg(t)
			})
		}
		cmd = tpCmd
	}

	return m, cmd
}

func (m Model) View() string {
	switch m.Global.CurrentView {
	case "menu":
		return m.Menu.View()
	case "player":
		return m.Player.View()
	case "track-picker":
		return m.TrackPicker.View()
	default:
		return "Unknown view"
	}
}
