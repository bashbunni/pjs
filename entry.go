package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/paginator"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/glamour"
	"github.com/charmbracelet/lipgloss"
)

const defaultEditor = "vi"

type Entry struct {
	path      string
	viewport  viewport.Model
	paginator paginator.Model
	entries   []string
}

func InitEntry(path string) Entry {
	vp := viewport.New(WindowSize.Width, WindowSize.Height)
	e := getEntries(path)
	p := paginator.New()
	p.SetTotalPages(len(e))
	entry := Entry{
		path,
		vp,
		p,
		e,
	}
	entry.setViewportContent()
	return entry
}

func openEditorCmd(path string) tea.Cmd {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = defaultEditor
	}
	c := exec.Command(editor, path)
	return tea.ExecProcess(c, func(err error) tea.Msg {
		// TODO: return the file contents to update viewport content
		contents, _ := os.ReadFile(path)
		return editorFinishedMsg{err, contents}
	})
}

// NewFilePath creates a markdown file to be opened in the editor
func NewFilePath(path string) (filepath string) {
	today := time.Now().Format("2006-01-02")
	filepath = fmt.Sprintf("%s/%s.md", path, today)
	return filepath
}

// ReadFile returns the contents of the file as a string
func ReadFile(path string) (string, error) {
	out, err := os.ReadFile(path)
	if err != nil {
		return "", fmt.Errorf("%w: unable to read file: %s", err, path)
	}
	return string(out), nil
}

func getEntries(path string) []string {
	var entries []string
	de, err := os.ReadDir(path)
	if err != nil {
		fmt.Errorf("unable to read dir: %w", err)
	}

	for _, entry := range de {
		if !entry.IsDir() {
			entries = append(entries, entry.Name())
		}
	}
	return entries
}

/* tea model interface */

// TODO: get num entries for paginator
// TODO: load entries as needed

// Init get first entry
func (m Entry) Init() tea.Cmd {
	return nil
}

func (m Entry) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	// TODO: have main model handle resizes and quits
	var cmd tea.Cmd
	var cmds []tea.Cmd
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		WindowSize.Width = msg.Width
		WindowSize.Height = msg.Height
	case tea.KeyMsg:
		switch {
		case key.Matches(msg, Keymap.Quit):
			return m, tea.Quit
		case key.Matches(msg, Keymap.Back):
			// TODO: don't re-init list each time
			return InitModel(), nil

		case key.Matches(msg, Keymap.Create):
			cmds = append(cmds, openEditorCmd(NewFilePath(m.path)))
		case key.Matches(msg, Keymap.Open):
			cmds = append(cmds, openEditorCmd(m.currentFile()))
		}
	}
	m.viewport, cmd = m.viewport.Update(msg)
	cmds = append(cmds, cmd)
	m.paginator, cmd = m.paginator.Update(msg)
	cmds = append(cmds, cmd)
	m.setViewportContent() // refresh the content on every Update call
	return m, tea.Batch(cmds...)
}

func (m Entry) View() string {
	return lipgloss.JoinVertical(lipgloss.Left, m.viewport.View(), m.paginator.View())
}

func (m *Entry) setViewportContent() {
	var content string
	if len(m.entries) == 0 {
		content = "There are no entries for this project :)"
	} else {
		content, _ = ReadFile(m.currentFile())
	}
	str, _ := glamour.Render(content, "dark")
	m.viewport.SetContent(str)
}

func (m *Entry) currentFile() string {
	return fmt.Sprintf("%s/%s.md", m.path, m.entries[m.paginator.Page])
}
