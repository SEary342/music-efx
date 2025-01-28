package menu

import (
	"fmt"
	"os"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type MenuItem struct {
	Title       string
	Description string
}

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

type MenuModel struct {
	Title   string
	Filter  bool
	Status  bool
	list    list.Model
	Items   []MenuItem
	Choice  string
	Exiting bool
}

func (m MenuModel) Init() tea.Cmd {
	return nil
}

func (m *MenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "esc":
			m.Exiting = true
			return m, tea.Quit
		case "enter":
			i, ok := m.list.SelectedItem().(item)

			if ok {
				m.Choice = i.Title()
			}
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func (m MenuModel) View() string {
	return docStyle.Render(m.list.View())
}

func Menu(m *MenuModel) {
	var listItems []list.Item
	for _, mi := range m.Items {
		listItems = append(listItems, item{title: mi.Title, desc: mi.Description})
	}
	l := list.New(listItems, list.NewDefaultDelegate(), 0, 0)
	l.Title = m.Title
	l.SetFilteringEnabled(m.Filter)
	l.SetShowStatusBar(m.Status)

	m.list = l

	p := tea.NewProgram(m, tea.WithAltScreen())

	if _, err := p.Run(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}
