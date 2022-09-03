package entryui

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/bashbunni/project-management/utils"
	tea "github.com/charmbracelet/bubbletea"
)

const defaultEditor = "vim"

func openEditorCmd() tea.Cmd {
	file, err := os.CreateTemp(os.TempDir(), "")
	if err != nil {
		return func() tea.Msg {
			return errMsg{error: err}
		}
	}
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = defaultEditor
	}
	c := exec.Command(editor, file.Name())
	return tea.ExecProcess(c, func(err error) tea.Msg {
		return editorFinishedMsg{err, file}
	})
}

func (m Model) createEntryCmd(file *os.File) tea.Cmd {
	return func() tea.Msg {
		input, err := utils.ReadFile(file)
		if err != nil {
			return errMsg{fmt.Errorf("cannot read file in createEntryCmd: %v", err)}
			// TODO: why is this giving me an error when input != ""
		}
		if err := m.er.CreateEntry(input, m.activeProjectID); err != nil {
			return errMsg{fmt.Errorf("cannot create entry: %v", err)}
		}
		if err := os.Remove(file.Name()); err != nil {
			return errMsg{fmt.Errorf("cannot remove file: %v", err)}
		}
		if closeErr := file.Close(); closeErr != nil {
			return errMsg{fmt.Errorf("unable to close file: %v", err)}
		}
		return m.setupEntries()
	}
}
