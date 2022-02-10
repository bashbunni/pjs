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

func createProjectCmd(name string, pr *models.GormProjectRepository) tea.Cmd {
	return func() tea.Msg {
		project, err := pr.CreateProject(name)
		if err != nil {
			return errMsg{err}
		}
		return createProjectListMsg{project}
	}
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

