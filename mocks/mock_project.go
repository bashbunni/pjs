package mocks

import (
	"errors"
	"fmt"

	"github.com/bashbunni/project-management/models"
	"github.com/bashbunni/project-management/utils"
)

type MockProjectRepository struct {
	Projects map[uint]*models.Project
}

func (m MockProjectRepository) GetOrCreateProjectByID(projectID uint) models.Project {
	proj, err := m.getProjectByID(projectID)
	// make getProjectByID return 0 if not found
	if errors.Is(err, utils.ErrProjectNotFound) {
		return m.CreateProject("")
	}
	return proj
}

func (m MockProjectRepository) getProjectByID(projectID uint) (models.Project, error) {
	// make getProjectByID return 0 if not found
	if project, ok := m.Projects[projectID]; ok {
		return *project, nil
	}
	return models.Project{}, utils.ErrProjectNotFound
}

func (m MockProjectRepository) PrintProjects() {
	if m.hasProjects() {
		projects := m.GetAllProjects()
		for _, project := range projects {
			fmt.Printf(models.Format, project.ID, project.Name)
		}
	} else {
		fmt.Printf("There are no projects available")
	}
}

func (m MockProjectRepository) GetAllProjects() []models.Project {
	var projects []models.Project
	for _, project := range m.Projects {
		projects = append(projects, *project)
	}
	return projects
}

func (m MockProjectRepository) hasProjects() bool {
	if len(m.Projects) > 0 {
		return true
	}
	return false
}

func (m MockProjectRepository) CreateProject(name string) models.Project {
	if name == "" {
		name = models.NewProjectPrompt()
	}
	proj := &models.Project{Name: name}
	m.Projects[proj.ID] = proj
	return *proj
}

func (m MockProjectRepository) DeleteProject(pe *models.ProjectWithEntries, er models.EntryRepository) {
	// what if projectID does not exist?
	if _, ok := m.Projects[pe.Project.ID]; ok {
		delete(m.Projects, pe.Project.ID)
	}
	er.DeleteEntries(pe)
}

func (m MockProjectRepository) RenameProject(pe *models.ProjectWithEntries) {
	name := models.NewProjectPrompt()
	if project, ok := m.Projects[pe.Project.ID]; ok {
		project.Name = name
	}
}
