package frontend

import (
	"log"

	"github.com/bashbunni/project-management/models"
	tea "github.com/charmbracelet/bubbletea"
)

func updateEntryListCmd(activeProject int, er *models.GormEntryRepository) tea.Cmd {
	return func() tea.Msg {
		entries, err := er.GetEntriesByProjectID(uint(activeProject+1))
		log.Println(len(entries))
		if err != nil {
			return errMsg{err}
		}
		return updateEntryListMsg{entries}
	}
}

// open a dialogue to enter project name

func createProjectCmd(name string) {
/* TODO:
- hit c -> creates a new empty project
- hit enter -> saves the project
	- write to DB
- close typing dialogue and refresh list of projects
	- show a new view
*/
}

// TODO: implement
func createEntryCmd(activeProject int, er *models.GormEntryRepository) tea.Cmd {
	return func() tea.Msg {
		err := er.CreateEntry([]byte("hello"), uint(activeProject+1))
		if err != nil {
			return errMsg{err}
		}
		return updateEntryListCmd(activeProject, er)
	}
}

