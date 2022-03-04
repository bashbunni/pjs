package project

import (
	"fmt"
	"log"

	"github.com/pkg/errors"
	"gorm.io/gorm"
)

const (
	notFound uint   = 0
	format   string = "%d : %s\n"
)

// Project the project holds entries
type Project struct {
	gorm.Model
	Name string
}

// NewProject create a new project instance.
// DeletedAt defaults to the zero value for time.Time.
func NewProject(id uint, name string) *Project {
	return &Project{Name: name}
}

// Implement list.Item for Bubbletea TUI

// Title the project title to display in a list
func (p Project) Title() string { return p.Name }

// Description the project description to display in a list
func (p Project) Description() string { return fmt.Sprintf("%d", p.ID) }

// FilterValue choose what field to use for filtering in a Bubbletea list component
func (p Project) FilterValue() string { return p.Name }

// Repository CRUD operations for Projects
type Repository interface {
	PrintProjects()
	HasProjects() bool
	GetProjectByID(projectID uint) (Project, error)
	GetAllProjects() ([]Project, error)
	CreateProject(name string) (Project, error)
	DeleteProject(projectID uint) error
	RenameProject(projectID uint) error
}

// GormRepository holds the gorm DB and is a ProjectRepository
type GormRepository struct {
	DB *gorm.DB
}

// GetProjectByID get a project by ID
func (g *GormRepository) GetProjectByID(projectID uint) (Project, error) {
	var project Project
	if err := g.DB.Where("id = ?", projectID).First(&project).Error; err != nil {
		return project, errors.Wrap(err, cannotFindProject)
	}
	return project, nil
}

// PrintProjects print all projects to the console
func (g *GormRepository) PrintProjects() {
	projects, err := g.GetAllProjects()
	if err != nil {
		log.Fatal(err)
	}
	for _, project := range projects {
		fmt.Printf(format, project.ID, project.Name)
	}
}

// GetAllProjects retrieve all projects from the database
func (g *GormRepository) GetAllProjects() ([]Project, error) {
	var projects []Project
	if err := g.DB.Find(&projects).Error; err != nil {
		return projects, errors.Wrap(err, emptyTable)
	}
	return projects, nil
}

// HasProjects see if a database has any projects
func (g *GormRepository) HasProjects() bool {
	if projects, _ := g.GetAllProjects(); len(projects) == 0 {
		return false
	}
	return true
}

// CreateProject add a new project to the database
func (g *GormRepository) CreateProject(name string) (Project, error) {
	proj := Project{Name: name}
	if err := g.DB.Create(&proj).Error; err != nil {
		return proj, errors.Wrap(err, cannotCreateProject)
	}
	return proj, nil
}

// DeleteProject delete a project by ID
func (g *GormRepository) DeleteProject(projectID uint) error {
	if err := g.DB.Delete(&Project{}, projectID).Error; err != nil {
		return errors.Wrap(err, cannotDeleteProject)
	}
	return nil
}

// RenameProject rename an existing project
func (g *GormRepository) RenameProject(id uint, name string) {
	var newProject Project
	if err := g.DB.Where("id = ?", id).First(&newProject).Error; err != nil {
		log.Fatalf("Unable to rename project: %q", err)
	}
	newProject.Name = name
	if err := g.DB.Save(&newProject).Error; err != nil {
		log.Fatalf("Unable to save project: %q", err)
	}
}

// NewProjectPrompt create a new project from user input to console
func NewProjectPrompt() string {
	var name string
	fmt.Println("what would you like to name your project?")
	fmt.Scanf("%s", &name)
	return name
}
