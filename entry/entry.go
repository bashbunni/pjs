package entry

import (
	"gorm.io/gorm"
)

// Entry the entry model
type Entry struct {
	gorm.Model
	ProjectID uint `gorm:"foreignKey:Project"`
	Message   string
}

// Repository the CRUD functionality for entries
type Repository interface {
	DeleteEntryByID(entryID uint) error
	DeleteEntries(projectID uint) error
	GetEntriesByProjectID(projectID uint) ([]Entry, error)
	CreateEntry(message []byte, projectID uint) error
}

// GormRepository holds the gorm DB and is a EntryRepository
type GormRepository struct {
	DB *gorm.DB
}

// DeleteEntryByID delete an entry by its ID
func (g *GormRepository) DeleteEntryByID(entryID uint) error {
	result := g.DB.Delete(&Entry{}, entryID)
	return result.Error
}

// TODO: unused
// DeleteEntries delete all entries for a given project
func (g *GormRepository) DeleteEntries(projectID uint) error {
	result := g.DB.Where("project_id = ?", projectID).Delete(&Entry{})
	return result.Error
}

// GetEntriesByProjectID get all entries for a given project
func (g *GormRepository) GetEntriesByProjectID(projectID uint) ([]Entry, error) {
	var Entries []Entry
	result := g.DB.Where("project_id = ?", projectID).Find(&Entries)
	return Entries, result.Error
}

// CreateEntry create a new entry in the database
func (g *GormRepository) CreateEntry(message []byte, projectID uint) error {
	entry := Entry{Message: string(message[:]), ProjectID: projectID}
	result := g.DB.Create(&entry)
	return result.Error
}
