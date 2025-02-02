package files

import (
	"errors"
	"strings"
	"time"

	"github.com/charmbracelet/bubbles/filepicker"
	tea "github.com/charmbracelet/bubbletea"
)

type FileModel struct {
	filepicker   filepicker.Model
	SelectedFile string
	Quitting     bool
	err          error
}

type clearErrorMsg struct{}

func clearErrorAfter(t time.Duration) tea.Cmd {
	return tea.Tick(t, func(_ time.Time) tea.Msg {
		return clearErrorMsg{}
	})
}

func (m FileModel) Init() tea.Cmd {
	return m.filepicker.Init()
}

func (m FileModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			m.Quitting = true
			return m, nil
		}
	case clearErrorMsg:
		m.err = nil
	}
	var cmd tea.Cmd
	m.filepicker, cmd = m.filepicker.Update(msg)

	// Did the user select a file?
	if didSelect, path := m.filepicker.DidSelectFile(msg); didSelect {
		// Get the path of the selected file.
		m.SelectedFile = path
	}

	// Did the user select a disabled file?
	// This is only necessary to display an error to the user.
	if didSelect, path := m.filepicker.DidSelectDisabledFile(msg); didSelect {
		// Let's clear the selectedFile and display an error.
		m.err = errors.New(path + " is not valid.")
		m.SelectedFile = ""
		return m, tea.Batch(cmd, clearErrorAfter(2*time.Second))
	}

	return m, cmd
}

func (m FileModel) View() string {
	if m.Quitting {
		return ""
	}
	var s strings.Builder
	s.WriteString("\n  ")
	if m.err != nil {
		s.WriteString(m.filepicker.Styles.DisabledFile.Render(m.err.Error()))
	} else if m.SelectedFile == "" {
		s.WriteString("Pick a file:")
	} else {
		s.WriteString("Selected file: " + m.filepicker.Styles.Selected.Render(m.SelectedFile))
	}
	s.WriteString("\n\n" + m.filepicker.View() + "\n")
	return s.String()
}

func InitFilePicker(fileType string, startingDir string, dirAllowed bool) FileModel {
	fp := filepicker.New()
	fp.AllowedTypes = []string{fileType}
	fp.DirAllowed = dirAllowed
	fp.CurrentDirectory = startingDir
	fp.ShowPermissions = false
	fp.Height = 10
	fp.ShowSize = false
	return FileModel{
		filepicker: fp,
	}
}

/*
func main() {
	m := InitFilePicker(".mp3", "/home", false)
	tm, _ := tea.NewProgram(&m).Run()
	mm := tm.(FileModel)
	fmt.Println("\n  You selected: " + m.filepicker.Styles.Selected.Render(mm.selectedFile) + "\n")
}*/

/*
tm, _ := tea.NewProgram(&m).Run()
mm := tm.(FileModel)
fmt.Println("\n  You selected: " + m.filepicker.Styles.Selected.Render(mm.selectedFile) + "\n")
*/
