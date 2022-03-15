package frontend

import (
	"fmt"
	"os"

	"github.com/bashbunni/project-management/entry"
	"github.com/bashbunni/project-management/frontend/entryui"
	"github.com/bashbunni/project-management/frontend/projectui"
	"github.com/bashbunni/project-management/project"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	p    *tea.Program
	cmd  tea.Cmd
	cmds []tea.Cmd
)

type sessionState int

const (
	projectView sessionState = iota
	entryView
)

// implements tea.Model (Init, Update, View)
type mainModel struct {
	state   sessionState
	project projectui.Model
	entry   entryui.Model
	pr      *project.GormRepository
	er      *entry.GormRepository
	mode    string
}

// StartTea the entry point for the UI. Initializes the model.
func StartTea(pr project.GormRepository, er entry.GormRepository) {
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

	m := mainModel{}
	m.project = *projectui.New(input, &pr, &er, "projects")
	p = tea.NewProgram(m)
	p.EnterAltScreen()
	if err := p.Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func (m mainModel) Init() tea.Cmd {
	m.state = projectView
	return nil
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch m.state {
	case projectView:
		// update View
		newModel, newCmd := m.project.Update(msg)
		projectModel, ok := newModel.(projectui.Model)
		if !ok {
			panic("could not perform assertion on project model")
		}
		m.project = projectModel
		cmd = newCmd
	case entryView:
		// init entry view
		// TODO: add getActiveProjectID
		newModel, newCmd := m.entry.Update(msg)
		entryModel, ok := newModel.(entryui.Model)
		if !ok {
			panic("could not perform assertion on entry model")
		}
		m.entry = entryModel
		cmd = newCmd
	}
	return m, tea.Batch(cmds...)
}

func (m mainModel) View() string {
	return m.project.View()
}
