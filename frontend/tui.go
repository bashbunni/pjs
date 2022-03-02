package frontend

import (
	"fmt"
	"log"
	"os"

	"github.com/bashbunni/project-management/models"
	"github.com/bashbunni/project-management/outputs"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

var docStyle = lipgloss.NewStyle().Margin(1, 2)
var helpStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("241")).Render

var (
	p        *tea.Program
	renderer *glamour.TermRenderer
)

// implements tea.Model (Init, Update, View)
type model struct {
	state    string
	viewport viewport.Model
	list     list.Model
	input    textinput.Model
	pr       *models.GormProjectRepository
	er       *models.GormEntryRepository
	keymap   keymap
	mode     string
}

type keymap struct {
	create key.Binding
	enter  key.Binding
	rename key.Binding
	delete key.Binding
	back   key.Binding
}

// The entry point for the UI. Initializes the model.
func StartTea(pr models.GormProjectRepository, er models.GormEntryRepository) {
	if os.Getenv("HELP_DEBUG") != "" {
		if f, err := tea.LogToFile("debug.log", "help"); err != nil {
			fmt.Println("Couldn't open a file for logging:", err)
			os.Exit(1)
		} else {
			defer f.Close()
		}
	}

	input := textinput.New()
	input.Prompt = "$ "
	input.Placeholder = "Project name..."
	input.CharLimit = 250
	input.Width = 50

	var err error
	renderer, err = glamour.NewTermRenderer(
		glamour.WithStylePath("dracula"),
	)
	if err != nil {
		panic(err)
	}

	projects, err := pr.GetAllProjects()
	if err != nil {
		log.Fatal(err)
	}
	items := projectsToItems(projects)
	m := initProjectView(items, input, &pr, &er)
	p = tea.NewProgram(m)
	p.EnterAltScreen()
	if err := p.Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		// update entry view size
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - 8

		// update list size
		top, right, bottom, left := docStyle.GetMargin()
		m.list.SetSize(msg.Width-left-right, msg.Height-top-bottom-1)
	}

	switch m.state {
	case "viewProjectList":
		return m.handleProjectList(msg, cmds, cmd)
	case "viewEntries":
		return m.handleEntriesList(msg, cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	switch m.state {
	case "viewEntries":
		return docStyle.Render(m.viewport.View() + m.helpView())
	default:
		if m.input.Focused() {
			return docStyle.Render(m.list.View() + "\n" + m.input.View())
		}
		return docStyle.Render(m.list.View() + "\n")
	}
}

func getEntryMessagesByProjectIDAsSingleString(id uint, er *models.GormEntryRepository) (string, error) {
	entries, err := er.GetEntriesByProjectID(id)
	if err != nil {
		return "", err
	}
	return string(outputs.FormattedOutputFromEntries(entries)), nil
}
