package frontend

import (
	"fmt"
	"log"
	"os"

	"github.com/bashbunni/project-management/models"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

// TODO: render list of entries -> need to update m.list with entries instead of projects and change title to proj name
// TODO: stop cursor from moving in edit mode -> override list.Model

var docStyle = lipgloss.NewStyle().Margin(1, 2)

type item struct {
	title, desc string
}

func (i item) Title() string       { return i.title }
func (i item) Description() string { return i.desc }
func (i item) FilterValue() string { return i.title }

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
	err      error // TODO: does this get used
}

type keymap struct {
	create key.Binding
	enter  key.Binding
	rename key.Binding
	delete key.Binding
	back   key.Binding
}

func (m model) Init() tea.Cmd {
	return nil
}

// organize Update, View by state
// are we checking the messages for the state?
// I'm looking to do different things with each message depending on the state
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	var cmd tea.Cmd

	switch m.state {
	case "viewProjectList":
		return m.handleProjectList(msg, cmds, cmd)
	case "viewEntries":
		return m.handleEntriesList(msg, cmds, cmd)
	}
	return m, tea.Batch(cmds...)
}

func (m model) View() string {
	if m.input.Focused() {
		return docStyle.Render(m.list.View() + "\n" + m.input.View())
	}
	return docStyle.Render(m.list.View() + "\n")
}

// functions

// state functions
func StartTea(pr models.GormProjectRepository, er models.GormEntryRepository) {
	input := textinput.New()
	input.Prompt = "$ "
	input.Placeholder = "Project name..."
	input.CharLimit = 250
	input.Width = 50

	projects, err := pr.GetAllProjects()
	if err != nil {
		log.Fatal(err)
	}
	items := projectsToItems(projects)
	m := initProjectView(items, input, &pr, &er)

	tea.LogToFile("debug.log", "debug")
	p := tea.NewProgram(m)
	p.EnterAltScreen()

	if err := p.Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func (m model) handleEntriesList(msg tea.Msg, cmds []tea.Cmd, cmd tea.Cmd) (model, tea.Cmd) {
	m.initEntries()
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width
		// check for my cmds being done running
		// check for keypresses?
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		default:
			m.viewport, cmd = m.viewport.Update(msg)
		}
	}
	return m, tea.Batch(cmds...)
}

func (m model) initEntries() error {
	content := `# Hi
	I'm back
	## titles
	- and stuff`
	// change state to viewEntries
	m.state = "viewEntries"
	// init viewport (per glamour example)
	vp := viewport.New(78, 20)
	vp.Style = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		PaddingRight(2)
	renderer, err := glamour.NewTermRenderer(glamour.WithStylePath("notty"))
	if err != nil {
		return err
	}
	str, err := renderer.Render(content)
	if err != nil {
		return err
	}
	vp.SetContent(str)
	m.viewport = vp
	return nil
}
