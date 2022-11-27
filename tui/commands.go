package tui

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/bashbunni/project-management/database/models"
	"github.com/bashbunni/project-management/database/repos"
	"github.com/bashbunni/project-management/tui/constants"
	"github.com/bashbunni/project-management/utils"
	tea "github.com/charmbracelet/bubbletea"
)

const defaultEditor = "vim"

/* PROJECTS */

func openEditorCmd() tea.Cmd {
	file, err := os.CreateTemp(os.TempDir(), "")
	if err != nil {
		return func() tea.Msg {
			return errMsg{error: err}
		}
	}
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = defaultEditor
	}
	c := exec.Command(editor, file.Name())
	return tea.ExecProcess(c, func(err error) tea.Msg {
		return editorFinishedMsg{err, file}
	})
}

func (m Entry) createEntryCmd(file *os.File) tea.Cmd {
	return func() tea.Msg {
		input, err := utils.ReadFile(file)
		if err != nil {
			return errMsg{fmt.Errorf("cannot read file in createEntryCmd: %v", err)}
		}
		if err := constants.Er.CreateEntry(input, m.activeProjectID); err != nil {
			return errMsg{fmt.Errorf("cannot create entry: %v", err)}
		}
		if err := os.Remove(file.Name()); err != nil {
			return errMsg{fmt.Errorf("cannot remove file: %v", err)}
		}
		if closeErr := file.Close(); closeErr != nil {
			return errMsg{fmt.Errorf("unable to close file: %v", err)}
		}
		return m.setupEntries()
	}
}

/* ENTRIES */

func createProjectCmd(name string, pr repos.ProjectRepository) tea.Cmd {
	return func() tea.Msg {
		if err := pr.CreateProject(&models.Project{Name: name}); err != nil {
			return errMsg{err}
		}
		return updateProjectListMsg{}
	}
}

func renameProjectCmd(id uint, pr repos.ProjectRepository, name string) tea.Cmd {
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

func deleteProjectCmd(id uint, pr repos.ProjectRepository) tea.Cmd {
	return func() tea.Msg {
		err := pr.DeleteProject(id)
		if err != nil {
			return errMsg{err}
		}
		return updateProjectListMsg{}
	}
}
