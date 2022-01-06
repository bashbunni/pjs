package models

import (
	"fmt"
	"time"

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

// Create a new project instance.
// DeletedAt defaults to the zero value for time.Time.
func NewProject(id uint, name string) *Project {
	return &Project{ID: id, Name: name, DeletedAt: time.Time{}}
}

// Implement list.Item for Bubbletea TUI
func (p Project) Title() string       { return p.Name }
func (p Project) Description() string { return fmt.Sprintf("%d", p.ID) }
func (p Project) FilterValue() string { return p.Name }

// Interface
type ProjectRepository interface {
	GetOrCreateProjectByID(projectID int) Project
	PrintProjects()
	hasProjects() bool
	getProjectByID(projectId int) Project
	GetAllProjects() ([]Project, error)
	CreateProject(name string) Project
	DeleteProject(pe *ProjectWithEntries, er EntryRepository)
	RenameProject(pe *ProjectWithEntries)
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
	projects, err := g.GetAllProjects()
	if err != nil {
		fmt.Println("No projects found")
		return
	}
	for _, project := range projects {
		fmt.Printf(Format, project.ID, project.Name)
	}
}

func (g GormProjectRepository) GetAllProjects() ([]Project, error) {
	var projects []Project
	result := g.DB.Find(&projects)
	// TODO:
	// errors.Is(result.Error, gorm.ErrRecordNotFound)
	return projects, result.Error
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
