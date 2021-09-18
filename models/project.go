package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/bashbunni/project-management/utils"
	"gorm.io/gorm"
)

/*
have a struct for each source of data: e.g. DB, []Project
*/

const notFound uint = 0

// Entity
type Project struct {
	gorm.Model
	Name string
}

// Interface
type ProjectRepository interface {
	GetOrCreateProjectByID(projectID int) Project
	PrintProjects()
	hasProjects() bool
	getProjectByID(projectId int, db *gorm.DB) Project
	GetAllProjects() []Project
	CreateProject(name string) Project
}

// Mock Implementation
type MockProjectRepository struct {
	Projects []Project
}

func (m MockProjectRepository) GetOrCreateProjectByID(projectID uint) Project {
	proj := m.getProjectByID(projectID)
	if proj.ID == notFound {
		return m.CreateProject("")
	}
	return proj
}

func (m MockProjectRepository) getProjectByID(projectID uint) Project {
	for _, p := range m.Projects {
		if projectID == p.ID {
			return p
		}
	}
	return Project{}
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
	return m.Projects
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
	proj := Project{Name: name}
	m.Projects = append(m.Projects, proj)
	return proj
}

// Gorm implementation

type GormProjectRepository struct {
	db *gorm.DB
}

func (g GormProjectRepository) GetOrCreateProjectByID(projectID int) Project {
	proj := g.getProjectByID(projectID)
	if proj.ID == notFound {
		return g.CreateProject("")
	}
	return proj
}

func (g GormProjectRepository) getProjectByID(projectId int) Project {
	var project Project
	g.db.Where("id = ?", projectId).Find(&project)
	return project
}

func (g GormProjectRepository) PrintProjects() {
	if g.hasProjects() {
		projects := g.GetAllProjects()
		for _, project := range projects {
			fmt.Printf(Format, project.ID, project.Name)
		}
	} else {
		fmt.Printf("There are no projects available")
	}
}

func (g GormProjectRepository) GetAllProjects() []Project {
	var projects []Project
	if g.hasProjects() {
		g.db.Find(&projects)
	}
	return projects
}

func (g GormProjectRepository) hasProjects() bool {
	var projects []Project
	if err := g.db.Find(&projects).Error; err != nil {
		return false
	}
	return true
}

func (g GormProjectRepository) CreateProject(name string) Project {
	if name == "" {
		name = newProjectPrompt()
	}
	proj := Project{Name: name}
	g.db.Create(&proj)
	return proj
}

// TODO: check for cascade delete functionality for GORM
func DeleteProject(pe *ProjectWithEntries, db *gorm.DB) {
	// what if projectID does not exist?
	DeleteEntries(pe, db)
	db.Delete(&Project{}, pe.Project.ID)
}

// TODO: make pe's Project a *Project instead to simplify?
func RenameProject(pe *ProjectWithEntries, db *gorm.DB) {
	name := newProjectPrompt()
	var project Project
	db.Where("id = ?", pe.Project.ID).First(&project)
	project.Name = name
	pe.Project.Name = name
	db.Save(&project)
}

func newProjectPrompt() string {
	var name string
	fmt.Println("what would you like to name your project?")
	fmt.Scanf("%s", &name)
	return name
}
