package frontend

import (
	"github.com/bashbunni/project-management/models"
	"github.com/bashbunni/project-management/utils"
	tea "github.com/charmbracelet/bubbletea"
)

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

func (m model) createEntryCmd(activeProject uint, er *models.GormEntryRepository) tea.Cmd {
	return func() tea.Msg {
		p.ReleaseTerminal()
		err := er.CreateEntry(utils.CaptureInputFromFile(), activeProject)
		if err != nil {
			return errMsg{err}
		}
		p.RestoreTerminal()
		return updateEntryListMsg{}
	}
}
