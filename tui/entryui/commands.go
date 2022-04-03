package entryui

import (
	"log"

	"github.com/bashbunni/project-management/entry"
	"github.com/bashbunni/project-management/utils"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) createEntryCmd(activeProject uint, er *entry.GormRepository) tea.Cmd {
	return func() tea.Msg {
		err := m.p.ReleaseTerminal()
		if err != nil {
			log.Print(err)
			return errMsg{err}
		}
		input, err := utils.CaptureInputFromFile()
		if err != nil {
			log.Print(err)
			return errMsg{err}
		}
		err = er.CreateEntry(input, activeProject)
		if err != nil {
			log.Print(err)
			return errMsg{err}
		}
		err = m.p.RestoreTerminal()
		if err != nil {
			log.Print(err)
			return errMsg{err}
		}
		return updateEntryListMsg{}
	}
}
