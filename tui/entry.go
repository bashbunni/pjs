package tui

import (
	"fmt"
	"log"
	"os"

	"github.com/bashbunni/project-management/entry"
	"github.com/bashbunni/project-management/tui/constants"
	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/paginator"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

type (
	errMsg         struct{ error }
	UpdatedEntries []entry.Model
	UpdateMe       struct{}
	BackMsg        bool
)

type editorFinishedMsg struct {
	err  error
	file *os.File
}

var cmd tea.Cmd

// Entry implements tea.Model
type Entry struct {
	viewport        viewport.Model
	activeProjectID uint
	error           string
	windowSize      tea.WindowSizeMsg
	paginator       paginator.Model
	entries         []entry.Model
}

// Init run any intial IO on program start
func (m Entry) Init() tea.Cmd {
	return nil
}

// InitEntry initialize the entryui model for your program
func InitEntry(er *entry.GormRepository, activeProjectID uint, p *tea.Program, windowSize tea.WindowSizeMsg) *Entry {
	m := Entry{activeProjectID: activeProjectID, windowSize: windowSize}
	m.viewport = viewport.New(windowSize.Width, calculateHeight(windowSize.Height))
	m.viewport.Style = lipgloss.NewStyle().
		Align(lipgloss.Bottom)

	// init paginator
	m.paginator = paginator.New()
	m.paginator.Type = paginator.Dots
	m.paginator.ActiveDot = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "235", Dark: "252"}).Render("•")
	m.paginator.InactiveDot = lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{Light: "250", Dark: "238"}).Render("•")

	m.entries = m.setupEntries().(UpdatedEntries)
	m.paginator.SetTotalPages(len(m.entries))
	// set content
	m.setViewportContent()
	return &m
}

func (m *Entry) setupEntries() tea.Msg {
	var err error
	var entries []entry.Model
	if entries, err = constants.Er.GetEntriesByProjectID(m.activeProjectID); err != nil {
		return errMsg{fmt.Errorf("Cannot find project: %v", err)}
	}
	entries = entry.ReverseList(entries)
	return UpdatedEntries(entries)
}

func (m *Entry) setViewportContent() {
	var content string
	if len(m.entries) == 0 {
		content = "There are no entries for this project :)"
	} else {
		content = entry.FormatEntry(m.entries[m.paginator.Page])
	}
	str, err := glamour.Render(content, "dark")
	if err != nil {
		m.error = "could not render content with glamour"
	}
	m.viewport.SetContent(str)
}

// Update handle IO and commands
func (m Entry) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case BackMsg:
		return InitProject(), nil
	case tea.WindowSizeMsg:
		m.viewport.Width = msg.Width
		m.viewport.Height = calculateHeight(msg.Height)
	case errMsg:
		m.error = msg.Error()
	case editorFinishedMsg:
		if msg.err != nil {
			return m, tea.Quit
		}
		cmds = append(cmds, m.createEntryCmd(msg.file))
	case UpdatedEntries:
		log.Println("created new entry")
		log.Println(msg)
		m.entries = msg
		m.paginator.SetTotalPages(len(m.entries))
		m.setViewportContent()
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, constants.Keymap.Create):
			return m, openEditorCmd()
		case key.Matches(msg, constants.Keymap.Back):
			return m, func() tea.Msg {
				return BackMsg(true)
			}
		case key.Matches(msg, constants.Keymap.Quit):
			return m, tea.Quit
		}
	}

	m.paginator, cmd = m.paginator.Update(msg)
	cmds = append(cmds, cmd)
	m.setViewportContent() // refresh the content on every Update call
	return m, tea.Batch(cmds...)
}

func (m Entry) helpView() string {
	// TODO: use the keymaps to populate the help string
	return constants.HelpStyle("\n ↑/↓: navigate  • esc: back • c: create entry • d: delete entry • q: quit\n")
}

// View return the text UI to be output to the terminal
func (m Entry) View() string {
	formatted := lipgloss.JoinVertical(lipgloss.Left, "\n", m.viewport.View(), m.helpView(), constants.ErrStyle(m.error), m.paginator.View())
	return constants.DocStyle.Render(formatted)
}

/* helpers */

func calculateHeight(height int) int {
	return height - height/7
}
