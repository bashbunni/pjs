package main

import (
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
)

// TODO: fix GormRepository vs Repository
type (
	SyncProjects      struct{}
	editorFinishedMsg struct {
		err      error
		contents []byte
	}
	errMsg error
	mode   int
)

const (
	nav mode = iota
	edit
	create
)

var WindowSize struct {
	Height int
	Width  int
}

// Model the entryui model definition
type Model struct {
	mode     mode
	projects []Project
	list     list.Model
	input    textinput.Model
	quitting bool
}

// InitProject initialize the projectui model for your program
func InitModel() tea.Model {
	input := textinput.New()
	input.Prompt = "$ "
	input.Placeholder = "Project name..."
	input.CharLimit = 250
	input.Width = 50

	items, _ := newList()
	m := Model{mode: nav, list: list.NewModel(items, list.NewDefaultDelegate(), 8, 8), input: input}
	if WindowSize.Height != 0 {
		top, right, bottom, left := DocStyle.GetMargin()
		m.list.SetSize(WindowSize.Width-left-right, WindowSize.Height-top-bottom-1)
	}
	m.list.Title = "projects"
	m.list.AdditionalShortHelpKeys = func() []key.Binding {
		return []key.Binding{
			Keymap.Create,
			Keymap.Rename,
			Keymap.Delete,
			Keymap.Back,
		}
	}
	return m
}

func newList() ([]list.Item, error) {
	projects, err := getProjects()
	return projectsToItems(projects), err
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
		WindowSize.Width = msg.Width
		WindowSize.Height = msg.Height
		top, right, bottom, left := DocStyle.GetMargin()
		m.list.SetSize(msg.Width-left-right, msg.Height-top-bottom-1)
	case SyncProjects:
		items, _ := newList()
		m.list.SetItems(items)
	case tea.KeyMsg:
		if m.input.Focused() {
			if key.Matches(msg, Keymap.Enter) {
				if m.mode == create {
					cmds = append(cmds, createProjectCmd(m.input.Value()))
				}
				if m.mode == edit {
					cmds = append(cmds, renameProjectCmd(
						Project(m.list.SelectedItem().FilterValue()),
						m.input.Value()))
				}
				m.input.SetValue("")
				m.mode = nav
				m.input.Blur()
			}
			if key.Matches(msg, Keymap.Back) {
				m.input.SetValue("")
				m.mode = nav
				m.input.Blur()
			}
			// only log keypresses for the input field when it's focused
			m.input, cmd = m.input.Update(msg)
			cmds = append(cmds, cmd)
		} else {
			switch {
			case key.Matches(msg, Keymap.Create):
				m.mode = create
				m.input.Focus()
				cmd = textinput.Blink
			case key.Matches(msg, Keymap.Quit):
				m.quitting = true
				return m, tea.Quit
			case key.Matches(msg, Keymap.Enter):
				// TODO: show entries
			case key.Matches(msg, Keymap.Rename):
				m.mode = edit
				m.input.Focus()
				cmd = textinput.Blink
			case key.Matches(msg, Keymap.Delete):
				cmd = deleteProjectCmd(Project(m.list.SelectedItem().FilterValue()))
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
		return DocStyle.Render(m.list.View() + "\n" + m.input.View())
	}
	return DocStyle.Render(m.list.View() + "\n")
}

// TODO: use generics
// projectsToItems convert []Project to []list.Item
func projectsToItems(projects []Project) []list.Item {
	items := make([]list.Item, len(projects))
	for i, proj := range projects {
		items[i] = list.Item(proj)
	}
	return items
}

/* commands */

func createProjectCmd(name string) tea.Cmd {
	return func() tea.Msg {
		if _, err := NewProject(name); err != nil {
			return errMsg(err)
		}
		return SyncProjects{}
	}
}

func renameProjectCmd(p Project, name string) tea.Cmd {
	return func() tea.Msg {
		if err := p.Rename(name); err != nil {
			return errMsg(err)
		}
		return SyncProjects{}
	}
}

func deleteProjectCmd(p Project) tea.Cmd {
	return func() tea.Msg {
		if err := p.Delete(); err != nil {
			return errMsg(err)
		}
		return SyncProjects{}
	}
}
