package entry

import (
	"github.com/bashbunni/project-management/database/dbconn"
	"github.com/bashbunni/project-management/database/models"
	"gorm.io/gorm"
)

// Repository the CRUD functionality for entries
type Repository interface {
	DeleteEntryByID(entryID uint) error
	DeleteEntries(projectID uint) error
	GetEntriesByProjectID(projectID uint) ([]models.Entry, error)
	CreateEntry(message []byte, projectID uint) error
}

// GormRepository holds the gorm DB and is a EntryRepository
type GormRepository struct {
	DB  *gorm.DB
	WDB dbconn.GormWrapper
}

// DeleteEntryByID delete an entry by its ID
func (g *GormRepository) DeleteEntryByID(entryID uint) error {
	result := g.WDB.Delete(&models.Entry{}, entryID)
	return result.Error()
}

// TODO: unused
// DeleteEntries delete all entries for a given project
func (g *GormRepository) DeleteEntries(projectID uint) error {
	result := g.WDB.Where("project_id = ?", projectID).Delete(&models.Entry{})
	return result.Error()
}

// GetEntriesByProjectID get all entries for a given project
func (g *GormRepository) GetEntriesByProjectID(projectID uint) ([]models.Entry, error) {
	var Entries []models.Entry
	result := g.WDB.Where("project_id = ?", projectID).Find(&Entries)
	return Entries, result.Error()
}

// CreateEntry create a new entry in the database
func (g *GormRepository) CreateEntry(message []byte, projectID uint) error {
	entry := models.Entry{Message: string(message[:]), ProjectID: projectID}
	result := g.WDB.Create(&entry)
	return result.Error()
}
