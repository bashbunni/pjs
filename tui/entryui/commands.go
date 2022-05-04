package entryui

import (
	"os"
	"os/exec"
	"github.com/pkg/errors"
	"github.com/bashbunni/project-management/utils"
	tea "github.com/charmbracelet/bubbletea"
)

func openEditorCmd() tea.Cmd {
	file, err := os.CreateTemp(os.TempDir(), "")
	if err != nil {
		return func() tea.Msg {
			return errMsg{errors.Wrap(err, "cannot create temp file")}
		}
	}
	filename := file.Name()
	c := exec.Command(os.Getenv("EDITOR"), filename)
	return tea.Exec(tea.WrapExecCommand(c), func(err error) tea.Msg {
		return editorFinishedMsg{err, file}
	})
}

func (m Model) updateEntriesCmd() tea.Msg {
	m.setViewportContent()
	return updatedMsg{}
}

func (m Model) createEntryCmd(file *os.File) tea.Cmd {
	return func() tea.Msg {
	defer file.Close()
	input, err := utils.ReadFile(file.Name())
	if err != nil {
	   return errMsg{errors.Wrap(err, "cannot read file in createEntryCmd")}
	}
	if m.er.CreateEntry(input, m.activeProjectID); err != nil {
	   return errMsg{errors.Wrap(err, "cannot create entry")}
	}
	if err := os.Remove(file.Name()); err != nil {
	   return errMsg{errors.Wrap(err, "cannot remove file")}
	}
	return updateEntryListMsg{input}
	}
}
