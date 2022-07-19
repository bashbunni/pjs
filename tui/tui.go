package tui

import (
	"fmt"
	"log"
	"os"

	"github.com/bashbunni/project-management/entry"
	"github.com/bashbunni/project-management/project"
	"github.com/bashbunni/project-management/tui/entryui"
	"github.com/bashbunni/project-management/tui/projectui"
	tea "github.com/charmbracelet/bubbletea"
)

var p *tea.Program

type sessionState int

const (
	projectView sessionState = iota
	entryView
)

// MainModel the main model of the program; holds other models and bubbles
type MainModel struct {
	state           sessionState
	project         tea.Model
	entry           tea.Model
	pr              *project.GormRepository
	er              *entry.GormRepository
	activeProjectID uint
	windowSize      tea.WindowSizeMsg
}

// StartTea the entry point for the UI. Initializes the model.
func StartTea(pr project.GormRepository, er entry.GormRepository) {
	if f, err := tea.LogToFile("debug.log", "help"); err != nil {
		fmt.Println("Couldn't open a file for logging:", err)
		os.Exit(1)
	} else {
		defer func() {
			err = f.Close()
			if err != nil {
				log.Fatal(err)
			}
		}()
	}

	m := New(&pr, &er)
	p = tea.NewProgram(m)
	p.EnterAltScreen()
	if err := p.Start(); err != nil {
		fmt.Println("Error running program:", err)
		os.Exit(1)
	}
}

// New initialize the main model for your program
func New(pr *project.GormRepository, er *entry.GormRepository) MainModel {
	return MainModel{
		state:   projectView,
		project: projectui.New(pr, er),
		pr:      pr,
		er:      er,
	}
}

// Init run any intial IO on program start
func (m MainModel) Init() tea.Cmd {
	return nil
}

// Update handle IO and commands
func (m MainModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.windowSize = msg // pass this along to the entry view so it uses the full window size when it's initialized
	case entryui.BackMsg:
		m.state = projectView
	case projectui.SelectMsg:
		m.activeProjectID = msg.ActiveProjectID
		m.entry = entryui.New(m.er, m.activeProjectID, p, m.windowSize)
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
		newEntry, newCmd := m.entry.Update(msg)
		entryModel, ok := newEntry.(entryui.Model)
		if !ok {
			panic("could not perform assertion on entryui model")
		}
		return entryModel, newCmd
	}
	cmds = append(cmds, cmd)
	return m, tea.Batch(cmds...)
}

// View return the text UI to be output to the terminal
func (m MainModel) View() string {
	switch m.state {
	case entryView:
		return m.entry.View()
	default:
		return m.project.View()
	}
}
