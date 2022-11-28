package repos

import (
	"fmt"
	"log"

	"github.com/bashbunni/project-management/database/dbconn"
	"github.com/bashbunni/project-management/database/models"
)

const (
	format string = "%d : %s\n"
)

type ProjectRepository interface {
	PrintProjects()
	HasProjects() bool
	GetProjectByID(id uint) (models.Project, error)
	GetAllProjects() ([]models.Project, error)
	CreateProject(p *models.Project) error
	DeleteProject(id uint) error
	RenameProject(id uint, name string)
}

type projectRepo struct {
	dbConn dbconn.GormWrapper
}

func NewProjectRepo(db dbconn.GormWrapper) ProjectRepository {
	return projectRepo{dbConn: db}
}

func (r projectRepo) PrintProjects() {
	projects, err := r.GetAllProjects()
	if err != nil {
		log.Fatal(err)
	}
	for _, project := range projects {
		fmt.Printf(format, project.ID, project.Name)
	}
}

func (r projectRepo) HasProjects() bool {
	if projects, _ := r.GetAllProjects(); len(projects) == 0 {
		return false
	}
	return true
}
func (r projectRepo) GetProjectByID(id uint) (models.Project, error) {
	p := models.Project{}
	if err := r.dbConn.Where("id = ?", id).First(&p).Error(); err != nil {
		return p, fmt.Errorf("failed to find project of ID %d: %w", id, err)
	}
	return p, nil
}

func (r projectRepo) GetAllProjects() ([]models.Project, error) {
	projects := []models.Project{}
	if err := r.dbConn.Find(&projects).Error(); err != nil {
		return projects, fmt.Errorf("no projects found: %w", err)
	}
	return projects, nil
}

// note(tauraamui): GORM auto updates a model's ID field
// with whatever on create, so we really just want to update
// the given model, not return a full copy of it
func (r projectRepo) CreateProject(p *models.Project) error {
	if err := r.dbConn.Create(p).Error(); err != nil {
		return fmt.Errorf("failed to create project: %w", err)
	}
	return nil
}

func (r projectRepo) DeleteProject(id uint) error {
	if err := r.dbConn.Delete(&models.Project{}, id).Error(); err != nil {
		return fmt.Errorf("failed to delete project: %w", err)
	}
	return nil
}

func (r projectRepo) RenameProject(id uint, name string) {
	var newProject models.Project
	var err error

	if selerr := r.dbConn.Where("id = ?", id).First(&newProject).Error(); err != nil {
		err = fmt.Errorf("failed to find project of ID %d: %w", id, selerr)
	}

	newProject.Name = name
	if saveerr := r.dbConn.Save(&newProject).Error(); saveerr != nil {
		saveerr = fmt.Errorf("failed to save name change: %w", saveerr)

		if err == nil {
			err = saveerr
		} else {
			err = fmt.Errorf("%v: %w", saveerr, err)
		}
	}

	log.Fatalf("failed to rename project: %q", err)
}
