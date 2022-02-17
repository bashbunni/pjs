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

/* tasks
TODO: add CRUD for entries from entry view
*/

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
		return docStyle.Render(m.viewport.View())
	default:
		if m.input.Focused() {
			return docStyle.Render(m.list.View() + "\n" + m.input.View())
		}
		return docStyle.Render(m.list.View() + "\n")
	}
}

// entries
func (m model) handleEntriesList(msg tea.Msg, cmds []tea.Cmd, cmd tea.Cmd) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width
	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
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

func initEntries(m model) (viewport.Model, error) {
	vp := viewport.New(78, 20)
	vp.Style = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		PaddingRight(2)
	renderer, err := glamour.NewTermRenderer(glamour.WithStylePath("notty"))
	if err != nil {
		return vp, err
	}
	// TODO: pass project ID of chosen project to this function
	content, err := getEntryMessagesByProjectIDAsSingleString(1, m.er)
	str, err := renderer.Render(content)
	if err != nil {
		return vp, err
	}
	vp.SetContent(str)
	return vp, nil
}
