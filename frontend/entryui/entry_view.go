package entryui

import (
	"github.com/bashbunni/project-management/entry"
	"github.com/bashbunni/project-management/frontend/constants"
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

type Model struct {
	state           string
	viewport        viewport.Model
	er              *entry.GormRepository
	activeProjectID uint
	cmds            []tea.Cmd
	p               *tea.Program
}

func New(er *entry.GormRepository, activeProjectID uint, p *tea.Program) *Model {
	return &Model{}
}

func (m Model) Init() tea.Cmd {
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
		return err
	}
	str, err := glamour.Render(content, "dark")
	if err != nil {
		return err
	}
	m.viewport.SetContent(str)
	return nil
}

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
			m.state = "projects"
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
