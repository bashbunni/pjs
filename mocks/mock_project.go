package mocks

import (
	"errors"
	"fmt"
	"time"

	"github.com/bashbunni/project-management/utils"
)

type MockProjectRepository struct {
	Projects map[uint]*Project
}

func (m MockProjectRepository) GetOrCreateProjectByID(projectID uint) Project {
	proj, err := m.getProjectByID(projectID)
	// make getProjectByID return 0 if not found
	if errors.Is(err, utils.ErrProjectNotFound) {
		return m.CreateProject("")
	}
	return proj
}

func (m MockProjectRepository) getProjectByID(projectID uint) (Project, error) {
	// make getProjectByID return 0 if not found
	if project, ok := m.Projects[projectID]; ok {
		return *project, nil
	}
	return Project{}, utils.ErrProjectNotFound
}

func (m MockProjectRepository) PrintProjects() {
	if m.hasProjects() {
		projects := m.GetAllProjects()
		for _, project := range projects {
			fmt.Printf(Format, project.ID, project.Name)
		}
	} else {
		fmt.Printf("There are no projects available")
	}
}

func (m MockProjectRepository) GetAllProjects() []Project {
	var projects []Project
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

func (m MockProjectRepository) CreateProject(name string) Project {
	if name == "" {
		name = newProjectPrompt()
	}
	proj := &Project{ID: uint(len(m.Projects) + 1), Name: name}
	m.Projects[proj.ID] = proj
	return *proj
}

func (m MockProjectRepository) DeleteProject(pe *ProjectWithEntries, er EntryRepository) {
	// what if projectID does not exist?
	if project, ok := m.Projects[pe.Project.ID]; ok {
		project.DeletedAt = time.Now()
	}
	er.DeleteEntries(pe)
}

func (m MockProjectRepository) RenameProject(pe *ProjectWithEntries) {
	name := newProjectPrompt()
	if project, ok := m.Projects[pe.Project.ID]; ok {
		project.Name = name
	}
}
