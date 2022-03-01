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

// implements tea.Model (Init, Update, View)
type model struct {
	state           string
	viewport        viewport.Model
	list            list.Model
	input           textinput.Model
	pr              *models.GormProjectRepository
	er              *models.GormEntryRepository
	keymap          keymap
	mode            string
	activeProjectID uint
	ready           bool
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
	if os.Getenv("HELP_DEBUG") != "" {
		if f, err := tea.LogToFile("debug.log", "help"); err != nil {
			fmt.Println("Couldn't open a file for logging:", err)
			os.Exit(1)
		} else {
			defer f.Close()
		}
	}
	p := tea.NewProgram(m)
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

func (m model) handleEntriesList(msg tea.Msg, cmds []tea.Cmd, cmd tea.Cmd) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case updateEntryListMsg:
		// update vp.SetContent
	case tea.WindowSizeMsg:
		// TODO: why isn't this working?
		if !m.ready {
			m.viewport = viewport.New(msg.Width, msg.Height-2)
			m.initEntries()
			m.viewport.YPosition = 1
			m.ready = true
		}
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.create):
			cmds = append(cmds, createEntryCmd(m.activeProjectID, m.er))
		case key.Matches(msg, m.keymap.back):
			m.state = "viewProjectList"
		case msg.String() == "ctrl+c":
			return m, tea.Quit
		case msg.String() == "q":
			return m, tea.Quit
		default:
			m.viewport, cmd = m.viewport.Update(msg)
			cmds = append(cmds, cmd)
		}
	}
	return m, tea.Batch(cmds...)
}

func getEntryMessagesByProjectIDAsSingleString(id uint, er *models.GormEntryRepository) (string, error) {
	entries, err := er.GetEntriesByProjectID(id)
	if err != nil {
		return "", err
	}
	return string(outputs.FormattedOutputFromEntries(entries)), nil
}

func (m *model) initEntries() error {
	//	vp := viewport.New(78, 20)
	m.viewport.Style = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		PaddingRight(2)
	renderer, err := glamour.NewTermRenderer(
		glamour.WithStylePath("dracula"),
		glamour.WithAutoStyle(),
	)
	if err != nil {
		return err
	}
	content, err := getEntryMessagesByProjectIDAsSingleString(m.getActiveProjectID(), m.er)
	if err != nil {
		return err
	}
	str, err := renderer.Render(content)
	if err != nil {
		return err
	}
	m.viewport.SetContent(str)
	return nil
}

func (m model) helpView() string {
	return helpStyle("\n ↑/↓: navigate  • esc: back • c: create entry • d: delete entry • q: quit\n")
}
