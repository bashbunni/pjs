package tui

import (
	"log"

	"github.com/bashbunni/project-management/project"
	"github.com/bashbunni/project-management/tui/constants"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// TODO: fix GormRepository vs Repository
type (
	updateProjectListMsg struct{}
	renameProjectMsg     []list.Item
)

// SelectMsg the message to change the view to the selected entry
type SelectMsg struct {
	ActiveProjectID uint
}

type mode int

const (
	nav mode = iota
	edit
	create
)

// Model the entryui model definition
type Model struct {
	mode     mode
	list     list.Model
	input    textinput.Model
	quitting bool
}

// InitProject initialize the projectui model for your program
func InitProject() tea.Model {
	input := textinput.New()
	input.Prompt = "$ "
	input.Placeholder = "Project name..."
	input.CharLimit = 250
	input.Width = 50

	items := newProjectList(constants.Pr)
	m := Model{mode: nav, list: list.NewModel(items, list.NewDefaultDelegate(), 8, 8), input: input}
	if constants.WindowSize.Height != 0 {
		top, right, bottom, left := constants.DocStyle.GetMargin()
		m.list.SetSize(constants.WindowSize.Width-left-right, constants.WindowSize.Height-top-bottom-1)
	}
	m.list.Title = "projects"
	m.list.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			constants.Keymap.Create,
			constants.Keymap.Rename,
			constants.Keymap.Delete,
			constants.Keymap.Back,
		}
	}
	return m
}

func newProjectList(pr *project.GormRepository) []list.Item {
	projects, err := pr.GetAllProjects()
	if err != nil {
		log.Fatal(err)
	}
	return projectsToItems(projects)
}

// Init run any intial IO on program start
func (m Model) Init() tea.Cmd {
	return nil
}

// Update handle IO and commands
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		constants.WindowSize = msg
		top, right, bottom, left := constants.DocStyle.GetMargin()
		m.list.SetSize(msg.Width-left-right, msg.Height-top-bottom-1)
	case updateProjectListMsg:
		projects, err := constants.Pr.GetAllProjects()
		if err != nil {
			log.Fatal(err)
		}
		items := projectsToItems(projects)
		m.list.SetItems(items)
		m.mode = nav
	case renameProjectMsg:
		m.list.SetItems(msg)
		m.mode = nav
	case tea.KeyMsg:
		if m.input.Focused() {
			if key.Matches(msg, constants.Keymap.Enter) {
				if m.mode == create {
					cmds = append(cmds, createProjectCmd(m.input.Value(), constants.Pr))
				}
				if m.mode == edit {
					cmds = append(cmds, renameProjectCmd(m.getActiveProjectID(), constants.Pr, m.input.Value()))
				}
				m.input.SetValue("")
				m.mode = nav
				m.input.Blur()
			}
			if key.Matches(msg, constants.Keymap.Back) {
				m.input.SetValue("")
				m.mode = nav
				m.input.Blur()
			}
			// only log keypresses for the input field when it's focused
			m.input, cmd = m.input.Update(msg)
			cmds = append(cmds, cmd)
		} else {
			switch {
			case key.Matches(msg, constants.Keymap.Create):
				m.mode = create
				m.input.Focus()
				cmd = textinput.Blink
			case key.Matches(msg, constants.Keymap.Quit):
				m.quitting = true
				return m, tea.Quit
			case key.Matches(msg, constants.Keymap.Enter):
				activeProject := m.list.SelectedItem().(project.Project)
				entry := InitEntry(constants.Er, activeProject.ID, constants.P)
				return entry.Update(constants.WindowSize)
			case key.Matches(msg, constants.Keymap.Rename):
				m.mode = edit
				m.input.Focus()
				cmd = textinput.Blink
			case key.Matches(msg, constants.Keymap.Delete):
				items := m.list.Items()
				if len(items) > 0 {
					cmd = deleteProjectCmd(m.getActiveProjectID(), constants.Pr)
				}
			default:
				m.list, cmd = m.list.Update(msg)
			}
			cmds = append(cmds, cmd)
		}
	}
	return m, tea.Batch(cmds...)
}

// View return the text UI to be output to the terminal
func (m Model) View() string {
	if m.quitting {
		return ""
	}
	if m.input.Focused() {
		return constants.DocStyle.Render(m.list.View() + "\n" + m.input.View())
	}
	return constants.DocStyle.Render(m.list.View() + "\n")
}

// TODO: use generics
// projectsToItems convert []model.Project to []list.Item
func projectsToItems(projects []project.Project) []list.Item {
	items := make([]list.Item, len(projects))
	for i, proj := range projects {
		items[i] = list.Item(proj)
	}
	return items
}

func (m Model) getActiveProjectID() uint {
	items := m.list.Items()
	activeItem := items[m.list.Index()]
	return activeItem.(project.Project).ID
}
