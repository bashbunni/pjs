package projectui

import (
	"github.com/bashbunni/project-management/project"
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
		projects, err := pr.GetAllProjects()
		if err != nil {
			return errMsg{err}
		}
		items := projectsToItems(projects)

		return renameProjectMsg(items)
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

func selectProjectCmd(ActiveProjectID uint) tea.Cmd {
	return func() tea.Msg {
		return SelectMsg{ActiveProjectID: ActiveProjectID}
	}
}
