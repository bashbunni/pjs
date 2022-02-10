package frontend

import (
	"fmt"
	"log"
	"os"

	"github.com/bashbunni/project-management/models"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

// TODO: track selected project
// TODO: render list of entries

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
	input textinput.Model
	active models.Project
	pr *models.GormProjectRepository
	er *models.GormEntryRepository
	keymap keymap
	editing bool
	err error
}

type keymap struct {
	create key.Binding
	enter key.Binding
	rename key.Binding
	delete key.Binding
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case updateEntryListMsg:
		entries := make([]list.Item, 0, len(msg.entries))
		for _, e := range msg.entries {
			entries = append(entries, item{
				title: fmt.Sprintf("%d", e.ID),
			})
		}
 	case createProjectListMsg:
 		projects, err := m.pr.GetAllProjects()
 		m.projects = list.NewModel(projectsToItems(projects), list.NewDefaultDelegate(), 0, 0)
 		if err != nil {
			m.err = err
			// TODO: have this display in status-bar in View
 		}
	case tea.KeyMsg:
		if !m.input.Focused() { 
			switch {
				case key.Matches(msg, m.keymap.create):
					m.input.Focus()
					cmds = append(cmds, textinput.Blink)
				case msg.String() == "ctrl+c":
					return m, tea.Quit
				case key.Matches(msg, m.keymap.enter):
					// TODO: update list
					return m, nil
				case key.Matches(msg, m.keymap.rename):
					// TODO: rename project
					return m, nil
				case key.Matches(msg, m.keymap.delete):
					// TODO: delete project
					return m, nil
			}
		}
		if m.input.Focused() {
			if key.Matches(msg, m.keymap.enter) {
				createProjectCmd(m.input.Value(), m.pr)
				m.input.Blur()
			}
		}

	case tea.WindowSizeMsg:
		top, right, bottom, left := docStyle.GetMargin()
		m.projects.SetSize(msg.Width-left-right, msg.Height-top-bottom)
	}

	m.projects, cmd = m.projects.Update(msg)
	cmds = append(cmds, cmd)
	m.input, cmd = m.input.Update(msg)
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	return docStyle.Render(m.projects.View() + "\n" + m.input.View())
}

// functions

// Initial model (AKA first View)
func ChooseProject(pr models.GormProjectRepository, er models.GormEntryRepository) {
	projects, err := pr.GetAllProjects()
	if err != nil {
		log.Fatal(err)
	}
	items := projectsToItems(projects)
	m := model{projects: list.NewModel(items, list.NewDefaultDelegate(), 0, 0), pr: &pr, er: &er, keymap: 
	keymap{
		create: key.NewBinding(
			key.WithKeys("c"),
			key.WithHelp("c", "create"),
		),
		enter: key.NewBinding(
			key.WithKeys("enter"),
			key.WithHelp("enter", "select"),
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

	tea.LogToFile("debug.log", "debug")
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
