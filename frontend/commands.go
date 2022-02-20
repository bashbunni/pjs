package frontend

import (
	"log"

	"github.com/bashbunni/project-management/models"
	"github.com/bashbunni/project-management/utils"
	tea "github.com/charmbracelet/bubbletea"
)

// projects

func createProjectCmd(name string, pr *models.GormProjectRepository) tea.Cmd {
	return func() tea.Msg {
		_, err := pr.CreateProject(name)
		if err != nil {
			return errMsg{err}
		}
		return updateProjectListMsg{}
	}
}

func renameProjectCmd(id uint, pr *models.GormProjectRepository, name string) tea.Cmd {
	return func() tea.Msg {
		pr.RenameProject(id, name)
		return renameProjectMsg{}
	}
}

func deleteProjectCmd(id uint, pr *models.GormProjectRepository) tea.Cmd {
	return func() tea.Msg {
		err := pr.DeleteProject(id)
		if err != nil {
			return errMsg{err}
		}
		return updateProjectListMsg{}
	}
}

// entries

// TODO: implement
func createEntryCmd(activeProject uint, er *models.GormEntryRepository) tea.Cmd {
	return func() tea.Msg {
		// TODO: doesn't work
		err := er.CreateEntry(utils.CaptureInputFromFile(), activeProject)
		if err != nil {
			return errMsg{err}
		}
		return nil
	}
}

func updateEntryListCmd(activeProject uint, er *models.GormEntryRepository) tea.Cmd {
	return func() tea.Msg {
		entries, err := er.GetEntriesByProjectID(activeProject)
		log.Println(len(entries))
		if err != nil {
			return errMsg{err}
		}
		return updateEntryListMsg{entries}
	}
}
