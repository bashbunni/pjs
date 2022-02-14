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

// projects

func createProjectCmd(name string, pr *models.GormProjectRepository) tea.Cmd {
	return func() tea.Msg {
		project, err := pr.CreateProject(name)
		if err != nil {
			return errMsg{err}
		}
		return createProjectListMsg{project}
	}
}

func renameProjectCmd(id uint, pr *models.GormProjectRepository, name string) tea.Cmd {
	return func() tea.Msg {
		pr.RenameProject(id, name)
		return renameProjectMsg{}
	}
}

// entries

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

