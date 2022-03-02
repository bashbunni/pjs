package frontend

import (
	"github.com/charmbracelet/bubbles/key"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
)

func (m model) handleEntriesList(msg tea.Msg, cmds []tea.Cmd, cmd tea.Cmd) (model, tea.Cmd) {
	switch msg := msg.(type) {
	case updateEntryListMsg:
		// update vp.SetContent
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, m.keymap.create):
			cmds = append(cmds, m.createEntryCmd(m.getActiveProjectID(), m.er))
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

func (m *model) initEntries() error {
	//	vp := viewport.New(78, 20)
	m.viewport.Style = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		PaddingRight(2)
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
