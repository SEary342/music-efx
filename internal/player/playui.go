package player

import (
	"fmt"
	"os"
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

func PlayUI(trackPath string) {
	fileName := filepath.Base(trackPath)

	m := model{
		progress: progress.New(progress.WithDefaultGradient()),
		title:    strings.TrimSuffix(fileName, filepath.Ext(fileName)),
		player:   &Player{},
		stopChan: make(chan bool),
	}
	m.progress.ShowPercentage = false

	trk, _ := LoadTrack(trackPath)
	m.timeRemaining = trk.Length
	m.player.PlayTrack(trk)

	if _, err := tea.NewProgram(m).Run(); err != nil {
		fmt.Println("Oh no!", err)
		os.Exit(1)
	}
}

type tickMsg time.Time

type model struct {
	progress      progress.Model
	player        *Player
	stopChan      chan bool
	title         string
	timeRemaining time.Duration
}

func (m model) Init() tea.Cmd {
	return tea.Batch(tickCmd(), tea.ClearScreen)
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		return m, tea.Quit

	case tea.WindowSizeMsg:
		m.progress.Width = msg.Width - padding*2 - 4
		if m.progress.Width > maxWidth {
			m.progress.Width = maxWidth
		}
		return m, nil

	case tickMsg:
		if m.progress.Percent() == 1.0 {
			return m, tea.Quit
		}

		position := float64(m.player.track.Stream.Position()) / float64(m.player.track.Format.SampleRate)
		total := float64(m.player.track.Stream.Len()) / float64(m.player.track.Format.SampleRate)

		// Calculate the progress percentage
		progress := position / total
		cmd := m.progress.SetPercent(progress)
		m.timeRemaining = time.Duration((total - position) * float64(time.Second))
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

func (m model) View() string {
	pad := strings.Repeat(" ", padding)

	return "\n" +
		pad + m.title + "\n\n" +
		pad + m.progress.View() + pad + fmtDuration(m.timeRemaining) + "\n\n" +
		pad + helpStyle("Press any key to quit")
}

func tickCmd() tea.Cmd {
	return tea.Tick(500*time.Millisecond, func(t time.Time) tea.Msg {
		return tickMsg(t)
	})
}
