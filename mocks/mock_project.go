package mocks

import (
	"fmt"

	"github.com/bashbunni/project-management/models"
)

type MockProjectRepository struct {
	Projects map[uint]*models.Project
}

func (m MockProjectRepository) GetProjectByID(projectID uint) models.Project {
	// make getProjectByID return 0 if not found
	if project, ok := m.Projects[projectID]; ok {
		return *project 
	}
	return models.Project{}
}

func (m MockProjectRepository) PrintProjects() {
	if m.HasProjects() {
		projects := m.GetAllProjects() // err is nil for mock
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

func (m MockProjectRepository) HasProjects() bool {
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
