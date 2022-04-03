package frontend

import (
	"fmt"
	"os"

	"github.com/bashbunni/project-management/entry"
	"github.com/bashbunni/project-management/frontend/entryui"
	"github.com/bashbunni/project-management/frontend/projectui"
	"github.com/bashbunni/project-management/project"
	tea "github.com/charmbracelet/bubbletea"
)

var (
	p *tea.Program
)

type sessionState int

const (
	projectView sessionState = iota
	entryView
)

// implements tea.Model (Init, Update, View)
type mainModel struct {
	state           sessionState
	project         tea.Model
	entry           tea.Model
	pr              *project.GormRepository
	er              *entry.GormRepository
	activeProjectID uint
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

	m := New(&pr, &er)
	p = tea.NewProgram(m)
	p.EnterAltScreen()
	if err := p.Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

func New(pr *project.GormRepository, er *entry.GormRepository) mainModel {
	return mainModel{
		state:   projectView,
		project: projectui.New(pr, er),
		pr:      pr,
		er:      er,
	}
}

func (m mainModel) Init() tea.Cmd {
	return nil
}

func (m mainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case entryui.BackMsg:
		m.state = projectView
	case projectui.SelectMsg:
		m.activeProjectID = msg.ActiveProjectID
		m.state = entryView
	}

	switch m.state {
	case projectView:
		newProject, newCmd := m.project.Update(msg)
		projectModel, ok := newProject.(projectui.Model)
		if !ok {
			panic("could not perform assertion on projectui model")
		}
		m.project = projectModel
		cmd = newCmd
	case entryView:
		m.entry = *entryui.New(m.er, m.activeProjectID, p)
		newEntry, newCmd := m.entry.Update(msg)
		entryModel, ok := newEntry.(entryui.Model)
		if !ok {
			panic("could not perform assertion on entryui model")
		}
		m.entry = entryModel
		cmd = newCmd
	}
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

func (m mainModel) View() string {
	switch m.state {
	case entryView:
		return m.entry.View()
	default:
		return m.project.View()
	}
}
