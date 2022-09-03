package entryui

import (
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
	errMsg         struct{ error } // TODO: have this implement Error()
	UpdatedEntries []entry.Entry
	UpdateMe       struct{}
	BackMsg        bool
)

type editorFinishedMsg struct {
	err  error
	file *os.File
}

var cmd tea.Cmd

// Model implements tea.Model
type Model struct {
	viewport        viewport.Model
	er              *entry.GormRepository
	activeProjectID uint
	p               *tea.Program
	error           string
	windowSize      tea.WindowSizeMsg
	paginator       paginator.Model
	entries         []entry.Entry
}

// Init run any intial IO on program start
func (m Model) Init() tea.Cmd {
	return nil
}

// New initialize the entryui model for your program
func New(er *entry.GormRepository, activeProjectID uint, p *tea.Program, windowSize tea.WindowSizeMsg) *Model {
	m := Model{er: er, activeProjectID: activeProjectID, windowSize: windowSize}
	m.p = p
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

func (m *Model) setupEntries() tea.Msg {
	var err error
	var entries []entry.Entry
	if entries, err = m.er.GetEntriesByProjectID(m.activeProjectID); err != nil {
		log.Fatalf("failed to get entries: %v", err)
	}
	entries = entry.ReverseEntries(entries)
	return UpdatedEntries(entries)
}

func (m *Model) setViewportContent() {
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
func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd
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

func (m Model) helpView() string {
	// TODO: use the keymaps to populate the help string
	return constants.HelpStyle("\n ↑/↓: navigate  • esc: back • c: create entry • d: delete entry • q: quit\n")
}

// View return the text UI to be output to the terminal
func (m Model) View() string {
	formatted := lipgloss.JoinVertical(lipgloss.Left, "\n", m.viewport.View(), m.helpView(), constants.ErrStyle(m.error), m.paginator.View())
	return constants.DocStyle.Render(formatted)
}

/* helpers */

func calculateHeight(height int) int {
	return height - height/7
}
