package frontend

import (
	"fmt"
	"log"
	"os"

	"github.com/bashbunni/project-management/models"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
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
	// TODO: check if this is already happening with list.Model
	active models.Project
	pr *models.GormProjectRepository
	er *models.GormEntryRepository
	keymap keymap
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
	switch msg := msg.(type) {
	case updateEntryListMsg:
		entries := make([]list.Item, 0, len(msg.entries))
		for _, e := range msg.entries {
			entries = append(entries, item{
				title: fmt.Sprintf("%d", e.ID),
			})
		}
	case tea.KeyMsg:
		switch {
			case key.Matches(msg, m.keymap.create):
			// TODO: create project
				return m, createEntryCmd(m.projects.Cursor(), m.er)
			case msg.String() == "ctrl+c":
				return m, tea.Quit
			case key.Matches(msg, m.keymap.enter):
				return m, updateEntryListCmd(m.projects.Cursor(), m.er)
			case key.Matches(msg, m.keymap.rename):
				// TODO: rename project
				return m, nil
			case key.Matches(msg, m.keymap.delete):
				// TODO: delete project
				return m, nil
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

// TODO: change this to EntryRepository
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

// TODO: are these necessary
type updateEntryListMsg struct {
	entries []models.Entry
}
type errMsg struct {error}

func updateEntryListCmd(activeProject int, er *models.GormEntryRepository) tea.Cmd {
	return func() tea.Msg {
		entries, err := er.GetEntriesByProjectID(uint(activeProject+1))
		log.Println(len(entries))
		if err != nil {
			return errMsg{err}
		}
		return updateEntryListMsg{entries}
	}
}


func createEntryCmd(activeProject int, er *models.GormEntryRepository) tea.Cmd {
	return func() tea.Msg {
		err := er.CreateEntry([]byte("hello"), uint(activeProject+1))
		if err != nil {
			return errMsg{err}
		}
		return updateEntryListCmd(activeProject, er)
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
