package models

import (
	"fmt"
	"log"
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
	Name      string
	DeletedAt time.Time
}

// Create a new project instance.
// DeletedAt defaults to the zero value for time.Time.
func NewProject(id uint, name string) *Project {
	return &Project{Name: name, DeletedAt: time.Time{}}
}

// Implement list.Item for Bubbletea TUI
// TODO: change this
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

func (g *GormProjectRepository) GetOrCreateProjectByID(projectID int) Project {
	proj := g.getProjectByID(projectID)
	if proj.ID == notFound {
		return g.CreateProject("")
	}
	return proj
}

func (g *GormProjectRepository) getProjectByID(projectID int) Project {
	var project Project
	if err := g.DB.Where("id = ?", projectID).Find(&project).Error; err != nil {
		log.Fatalf("Unable to get project by ID: %q", err)
	}
	return project
}

func (g *GormProjectRepository) PrintProjects() {
	projects := g.GetAllProjects()
	for _, project := range projects {
		fmt.Printf(Format, project.ID, project.Name)
	}
}

func (g *GormProjectRepository) GetAllProjects() []Project {
	var projects []Project
	if err := g.DB.Find(&projects).Error; err != nil {
		log.Fatalf("Projects not found: %q", err)
	}
	return projects
}

func (g *GormProjectRepository) CreateProject(name string) Project {
	if name == "" {
		name = newProjectPrompt()
	}
	proj := Project{Name: name}
	if err := g.DB.Create(&proj).Error; err != nil {
		log.Fatalf("Unable to create project: %q", err)
	}
	return proj
}

// TODO: check for cascade delete functionality for GORM
func (g *GormProjectRepository) DeleteProject(pe *ProjectWithEntries, er EntryRepository) {
	er.DeleteEntries(pe)
	if err := g.DB.Delete(&Project{}, pe.Project.ID).Error; err != nil {
		log.Fatalf("Unable to delete project: %q", err)
	}
}

// TODO: make pe's Project a *Project instead to simplify?
func (g *GormProjectRepository) RenameProject(pe *ProjectWithEntries) {
	name := newProjectPrompt()
	var project Project
	if err := g.DB.Where("id = ?", pe.Project.ID).First(&project).Error; err != nil {
		log.Fatalf("Unable to rename project: %q", err)
	}
	project.Name = name
	pe.Project.Name = name
	if err := g.DB.Save(&project).Error; err != nil {
		log.Fatalf("Unable to save project: %q", err)
	}
}

// TODO: move helper?
// helpers

func newProjectPrompt() string {
	var name string
	fmt.Println("what would you like to name your project?")
	fmt.Scanf("%s", &name)
	return name
}
