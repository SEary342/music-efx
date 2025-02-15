package player

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

const (
	padding  = 2
	maxWidth = 80
)

var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("#626262")).Render

type TickMsg time.Time

type PlayerModel struct {
	TrackPath     string
	progress      progress.Model
	Player        *Player
	Stopped       bool
	title         string
	timeRemaining time.Duration
}

func (m *PlayerModel) StartTrack() *PlayerModel {
	// Set up the track and start playing
	fileName := filepath.Base(m.TrackPath)
	m.progress = progress.New(progress.WithDefaultGradient())
	m.title = strings.TrimSuffix(fileName, filepath.Ext(fileName))
	trk, _ := LoadTrack(m.TrackPath, false)
	m.Player = &Player{track: trk}
	m.Stopped = false
	m.progress.ShowPercentage = false

	m.timeRemaining = trk.Length
	m.Player.PlayTrack()
	return m
}

func (m PlayerModel) Init() tea.Cmd {
	return nil
}

func (m PlayerModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "esc":
			m.Player.Stop()
			m.Stopped = true
			return m, nil
		}
		return m, nil

	case tea.WindowSizeMsg:
		m.progress.Width = msg.Width - padding*2 - 4
		if m.progress.Width > maxWidth {
			m.progress.Width = maxWidth
		}
		return m, nil

	case TickMsg:
		if m.progress.Percent() == 1.0 {
			return m, nil
		}
		var cmd tea.Cmd
		if m.Player != nil && m.Player.track != nil {
			position := float64(m.Player.track.Stream.Position()) / float64(m.Player.track.Format.SampleRate)
			total := float64(m.Player.track.Stream.Len()) / float64(m.Player.track.Format.SampleRate)

			progress := position / total
			cmd = m.progress.SetPercent(progress)
			m.timeRemaining = time.Duration((total - position) * float64(time.Second))
		}
		return m, tea.Batch(tickCmd(), cmd)

	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd

	default:
		return m, nil
	}
}

func fmtDuration(d time.Duration) string {
	d = d.Round(time.Second)
	m := d / time.Minute
	d -= m * time.Minute
	s := d / time.Second

	return fmt.Sprintf("%02d:%02d", m, s)
}

func (m PlayerModel) View() string {
	pad := strings.Repeat(" ", padding)

	return "\n" +
		pad + m.title + "\n\n" +
		pad + m.progress.View() + pad + fmtDuration(m.timeRemaining) + "\n\n" +
		pad + helpStyle("Press any key to quit")
}

func tickCmd() tea.Cmd {
	return tea.Tick(500*time.Millisecond, func(t time.Time) tea.Msg {
		return TickMsg(t)
	})
}
