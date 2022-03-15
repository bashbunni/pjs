package entryui

import (
	"log"

	"github.com/bashbunni/project-management/entry"
	"github.com/bashbunni/project-management/utils"
	tea "github.com/charmbracelet/bubbletea"
)

func (m Model) createEntryCmd(activeProject uint, er *entry.GormRepository) tea.Cmd {
	return func() tea.Msg {
		m.p.ReleaseTerminal()
		err := er.CreateEntry(utils.CaptureInputFromFile(), activeProject)
		if err != nil {
			log.Print(err)
			return errMsg{err}
		}
		m.p.RestoreTerminal()
		return updateEntryListMsg{}
	}
}
