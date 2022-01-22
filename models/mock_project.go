package models

import (
	"fmt"
)

type MockProjectRepository struct {
	Projects map[uint]*Project
}

func (m MockProjectRepository) GetProjectByID(projectID uint) Project {
	// make getProjectByID return 0 if not found
	if project, ok := m.Projects[projectID]; ok {
		return *project 
	}
	return Project{}
}

func (m MockProjectRepository) PrintProjects() {
	if m.HasProjects() {
		projects := m.GetAllProjects() // err is nil for mock
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

func (m MockProjectRepository) HasProjects() bool {
	if len(m.Projects) > 0 {
		return true
	}
	return false
}

func (m MockProjectRepository) CreateProject(name string) Project {
	if name == "" {
		name = NewProjectPrompt()
	}
	proj := &Project{Name: name}
	m.Projects[proj.ID] = proj
	return *proj
}

func (m MockProjectRepository) DeleteProject(pe *ProjectWithEntries, er EntryRepository) {
	// what if projectID does not exist?
	if _, ok := m.Projects[pe.Project.ID]; ok {
		delete(m.Projects, pe.Project.ID)
	}
	er.DeleteEntries(pe)
}

func (m MockProjectRepository) RenameProject(pe *ProjectWithEntries) {
	name := NewProjectPrompt()
	if project, ok := m.Projects[pe.Project.ID]; ok {
		project.Name = name
	}
}
