package app

import (
	"context"
	"fmt"
	"music-efx/internal/files"
	"music-efx/internal/menu"
	"music-efx/internal/player"
	"music-efx/internal/playlist"
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
	Context     context.Context
	CancelFunc  context.CancelFunc
}

type Model struct {
	Global         *GlobalModel
	Menu           menu.MenuModel
	Player         player.PlayerModel
	TrackPicker    files.FileModel
	PlaylistPicker files.FileModel
}

func NewAppModel(title string) Model {
	pwd, _ := os.Getwd()

	ctx, cancel := context.WithCancel(context.Background())

	return Model{
		Global:         &GlobalModel{CurrentView: "menu", SharedData: make(map[string]interface{}), Context: ctx, CancelFunc: cancel},
		Menu:           menu.New(mainMenuItems, title, false, false),
		Player:         player.PlayerModel{},
		TrackPicker:    files.InitFilePicker(".mp3", pwd, false),
		PlaylistPicker: files.InitFilePicker(".yaml", pwd, false),
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
			m.Global.CancelFunc() // Cancel the context when exiting the app
			return m, tea.Quit
		}

		switch choice := m.Menu.Choice; choice {
		case "Folder Navigation":
			m.Menu.Choice = ""
			m.Global.CurrentView = "track-picker"
			updatedPicker, tpCmd := m.TrackPicker.Update(m.TrackPicker.Init()())
			m.TrackPicker = updatedPicker.(files.FileModel)
			cmd = tpCmd
		case "Playlist Selection":
			m.Menu.Choice = ""
			m.Global.CurrentView = "playlist-picker"
			updatedPicker, tpCmd := m.PlaylistPicker.Update(m.PlaylistPicker.Init()())
			m.PlaylistPicker = updatedPicker.(files.FileModel)
			cmd = tpCmd
		case "Auto-Playlist":
			m.Menu.Choice = ""
			m.Global.CurrentView = "player"
			m.Player = *playlist.StartAutoPlaylist(m.Global.Context, &m.Player)
		}

	case "player":
		updatedPlayer, playerCmd := m.Player.Update(msg)
		m.Player = updatedPlayer.(player.PlayerModel)
		if m.Player.Stopped {
			m.Global.CurrentView = "menu"
			m.Global.CancelFunc() // Cancel the context when playback is stopped
			ctx, cancel := context.WithCancel(context.Background())
			m.Global.Context = ctx
			m.Global.CancelFunc = cancel
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
	case "playlist-picker":
		updatedPicker, tpCmd := m.PlaylistPicker.Update(msg)
		m.PlaylistPicker = updatedPicker.(files.FileModel)
		if m.PlaylistPicker.Quitting {
			m.PlaylistPicker.Quitting = false
			m.Global.CurrentView = "menu"
		} else if len(m.PlaylistPicker.SelectedFile) > 0 {
			m.Global.CurrentView = "player"
			mp3Meta, err := playlist.LoadPlaylist(m.PlaylistPicker.SelectedFile)
			if err != nil {
				fmt.Println(err)
				m.Global.CurrentView = "menu"
			}
			m.PlaylistPicker.SelectedFile = ""
			m.Player = *playlist.RandomPlay(m.Global.Context, mp3Meta, &m.Player)
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
	case "playlist-picker":
		return m.PlaylistPicker.View()
	default:
		return "Unknown view"
	}
}
