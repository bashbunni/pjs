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
	ID        uint
	Name      string
	DeletedAt time.Time
}

// Interface
type ProjectRepository interface {
	GetOrCreateProjectByID(projectID int) Project
	PrintProjects()
	hasProjects() bool
	getProjectByID(projectId int) Project
	GetAllProjects() []Project
	CreateProject(name string) Project
	DeleteProject(pe *ProjectWithEntries, er EntryRepository)
	RenameProject(pe *ProjectWithEntries)
}

// Mock Implementation
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

// Gorm implementation
type GormProjectRepository struct {
	DB *gorm.DB
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
	g.DB.Where("id = ?", projectId).Find(&project)
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
		g.DB.Find(&projects)
	}
	return projects
}

func (g GormProjectRepository) hasProjects() bool {
	var projects []Project
	if err := g.DB.Find(&projects).Error; err != nil {
		return false
	}
	return true
}

func (g GormProjectRepository) CreateProject(name string) Project {
	if name == "" {
		name = newProjectPrompt()
	}
	proj := Project{Name: name}
	g.DB.Create(&proj)
	return proj
}

// TODO: check for cascade delete functionality for GORM
func (g GormProjectRepository) DeleteProject(pe *ProjectWithEntries, er EntryRepository) {
	// what if projectID does not exist?
	er.DeleteEntries(pe)
	g.DB.Delete(&Project{}, pe.Project.ID)
}

// TODO: make pe's Project a *Project instead to simplify?
func (g GormProjectRepository) RenameProject(pe *ProjectWithEntries) {
	name := newProjectPrompt()
	var project Project
	g.DB.Where("id = ?", pe.Project.ID).First(&project)
	project.Name = name
	pe.Project.Name = name
	g.DB.Save(&project)
}

func newProjectPrompt() string {
	var name string
	fmt.Println("what would you like to name your project?")
	fmt.Scanf("%s", &name)
	return name
}
