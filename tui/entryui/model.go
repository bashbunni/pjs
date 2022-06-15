package entryui

import (
	"log"

	"github.com/bashbunni/project-management/entry"
	"github.com/bashbunni/project-management/tui/constants"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

// TODO: move from ioutil to os/io -> e.g. os.CreateTemp

var cmd tea.Cmd

// BackMsg change state back to project view
type BackMsg bool

// Model entryui model
type Model struct {
	viewport        viewport.Model
	er              *entry.GormRepository
	activeProjectID uint
	p               *tea.Program
	error           string
	windowSize      tea.WindowSizeMsg
}

// Init run any intial IO on program start
func (m Model) Init() tea.Cmd {
	return nil
}

func calculateHeight(height int) int {
	return height - height/4
}

// New initialize the entryui model for your program
func New(er *entry.GormRepository, activeProjectID uint, p *tea.Program, windowSize tea.WindowSizeMsg) *Model {
	m := Model{er: er, activeProjectID: activeProjectID, windowSize: windowSize}
	m.p = p
	vp := viewport.New(windowSize.Width, calculateHeight(windowSize.Height))
	log.Printf("width of screen: %d", windowSize.Width)
	m.viewport = vp
	m.viewport.Style = lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("62")).
		Align(lipgloss.Bottom)
	m.setViewportContent()
	return &m
}

func (m *Model) setViewportContent() {
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
}

// Update handle IO and commands
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// TODO: fix viewport sizing with keypresses and on init
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width
		m.viewport.Height = calculateHeight(msg.Height)
	case errMsg:
		m.error = msg.Error()
	case editorFinishedMsg:
		if msg.err != nil {
			return m, tea.Quit
		}
		cmd = m.createEntryCmd(msg.file)
	case updateEntryListMsg:
		return m, m.updateEntriesCmd
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, constants.Keymap.Create):
			return m, openEditorCmd()
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
		}
	}
	return m, cmd
}

func (m Model) helpView() string {
	return constants.HelpStyle("\n ↑/↓: navigate  • esc: back • c: create entry • d: delete entry • q: quit\n")
}

func (m Model) errorView() string {
	return constants.ErrStyle(m.error)
}

// View return the text UI to be output to the terminal
func (m Model) View() string {
	formatted := lipgloss.JoinVertical(lipgloss.Left, m.viewport.View(), m.helpView(), m.errorView())
	return constants.DocStyle.Render(formatted)
}

func getEntryMessagesByProjectIDAsSingleString(id uint, er *entry.GormRepository) (string, error) {
	entries, err := er.GetEntriesByProjectID(id)
	if err != nil {
		return "", err
	}
	return string(entry.FormattedOutputFromEntries(entries)), nil
}
