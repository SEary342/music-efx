package menu

import (
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
	List    list.Model
	Items   []MenuItem
	Choice  string
	Exiting bool
}

func New(mainMenuItems []MenuItem, title string, filter bool, status bool) MenuModel {
	var listItems []list.Item
	for _, mi := range mainMenuItems {
		listItems = append(listItems, item{title: mi.Title, desc: mi.Description})
	}
	l := list.New(listItems, list.NewDefaultDelegate(), 0, 0)
	l.Title = title
	l.SetFilteringEnabled(filter)
	l.SetShowStatusBar(status)

	return MenuModel{
		Title:  title,
		Filter: filter,
		Status: status,
		List:   l,
	}
}

func (m MenuModel) Init() tea.Cmd {
	return nil
}

func (m MenuModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch keypress := msg.String(); keypress {
		case "ctrl+c", "esc":
			m.Exiting = true
			return m, tea.Quit
		case "enter":
			if selectedItem, ok := m.List.SelectedItem().(item); ok {
				m.Choice = selectedItem.Title()
			}
			return m, nil
		}
	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.List.SetSize(msg.Width-h, msg.Height-v)
	}

	var cmd tea.Cmd
	m.List, cmd = m.List.Update(msg)
	return m, cmd
}

func (m MenuModel) View() string {
	return docStyle.Render(m.List.View())
}
