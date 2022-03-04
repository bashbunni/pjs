package frontend

import (
	"log"

	"github.com/bashbunni/project-management/entry"
	"github.com/bashbunni/project-management/project"
	"github.com/bashbunni/project-management/utils"
	tea "github.com/charmbracelet/bubbletea"
)

func createProjectCmd(name string, pr *project.GormRepository) tea.Cmd {
	return func() tea.Msg {
		_, err := pr.CreateProject(name)
		if err != nil {
			return errMsg{err}
		}
		return updateProjectListMsg{}
	}
}

func renameProjectCmd(id uint, pr *project.GormRepository, name string) tea.Cmd {
	return func() tea.Msg {
		pr.RenameProject(id, name)
		return renameProjectMsg{}
	}
}

func deleteProjectCmd(id uint, pr *project.GormRepository) tea.Cmd {
	return func() tea.Msg {
		err := pr.DeleteProject(id)
		if err != nil {
			return errMsg{err}
		}
		return updateProjectListMsg{}
	}
}

func (m model) createEntryCmd(activeProject uint, er *entry.GormRepository) tea.Cmd {
	return func() tea.Msg {
		p.ReleaseTerminal()
		err := er.CreateEntry(utils.CaptureInputFromFile(), activeProject)
		if err != nil {
			log.Print(err)
			return errMsg{err}
		}
		p.RestoreTerminal()
		return updateEntryListMsg{}
	}
}
