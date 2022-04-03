package entryui

import (
	"github.com/bashbunni/project-management/entry"
	"github.com/bashbunni/project-management/tui/constants"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

var (
	cmd  tea.Cmd
	cmds []tea.Cmd
)

// BackMsg change state back to project view
type BackMsg bool

// Model entryui model
type Model struct {
	viewport        viewport.Model
	er              *entry.GormRepository
	activeProjectID uint
	cmds            []tea.Cmd
	p               *tea.Program
	error           string
}

// Init run any intial IO on program start
func (m Model) Init() tea.Cmd {
	return nil
}

// New initialize the entryui model for your program
func New(er *entry.GormRepository, activeProjectID uint, p *tea.Program) *Model {
	m := Model{er: er, activeProjectID: activeProjectID}
	vp := viewport.New(8, 8)
	m.viewport = vp
	m.viewport.Style = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		PaddingRight(2)
	content, err := getEntryMessagesByProjectIDAsSingleString(m.activeProjectID, m.er)
	if content == "" {
		content = "There are no entries for this project :)"
	}
	if err != nil {
		m.error = "cannot get entry messages as single string"
	}
	str, err := glamour.Render(content, "dark")
	if err != nil {
		m.error = "could not render content with glamour"
	}
	m.viewport.SetContent(str)
	return &m
}

// Update handle IO and commands
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width
		m.viewport.Height = msg.Height - 8
	case updateEntryListMsg:
		// update vp.SetContent
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, constants.Keymap.Create):
			cmds = append(cmds, m.createEntryCmd(m.activeProjectID, m.er))
		case key.Matches(msg, constants.Keymap.Back):
			return m, func() tea.Msg {
				return BackMsg(true)
			}
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

func (m Model) helpView() string {
	return constants.HelpStyle("\n ↑/↓: navigate  • esc: back • c: create entry • d: delete entry • q: quit\n")
}

// View return the text UI to be output to the terminal
func (m Model) View() string {
	return constants.DocStyle.Render(m.viewport.View() + m.helpView())
}

func getEntryMessagesByProjectIDAsSingleString(id uint, er *entry.GormRepository) (string, error) {
	entries, err := er.GetEntriesByProjectID(id)
	if err != nil {
		return "", err
	}
	return string(entry.FormattedOutputFromEntries(entries)), nil
}
