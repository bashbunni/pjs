package frontend

import (
	"fmt"
	"os"

	"github.com/bashbunni/project-management/models"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

// implements tea.Model (Init, Update, View)
type model struct {
	projects list.Model
	keymap keymap
}

type keymap struct {
	create key.Binding
	rename key.Binding
	delete key.Binding
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "ctrl+c" {
			return m, tea.Quit
		}
	case tea.WindowSizeMsg:
		top, right, bottom, left := docStyle.GetMargin()
		m.projects.SetSize(msg.Width-left-right, msg.Height-top-bottom)
	}

	var cmd tea.Cmd
	m.projects, cmd = m.projects.Update(msg)
	return m, cmd
}

func (m model) View() string {
	return docStyle.Render(m.projects.View())
}

// functions

func ChooseProject(projects []models.Project) {
	items := projectsToItems(projects)
	items = append(items, item{title: "Create Project", desc: "create and open a new project"})
	m := model{projects: list.NewModel(items, list.NewDefaultDelegate(), 0, 0), keymap: 
	keymap{
		create: key.NewBinding(
			key.WithKeys("c"),
			key.WithHelp("c", "create"),
		),
		rename: key.NewBinding(
			key.WithKeys("r"),
			key.WithHelp("r", "rename"),
		),	
		delete: key.NewBinding(
			key.WithKeys("d"),
			key.WithHelp("d", "delete"),
		),
	},
}
	m.projects.Title = "Projects"
	m.projects.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			m.keymap.create,
			m.keymap.rename,
			m.keymap.delete,
		}
	}

	p := tea.NewProgram(m)
	p.EnterAltScreen()

	if err := p.Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

// convert []model.Project to []list.Item
func projectsToItems(projects []models.Project) []list.Item {
	items := make([]list.Item, len(projects))
	for i, proj := range projects {
		items[i] = list.Item(proj)
	}
	return items
}
