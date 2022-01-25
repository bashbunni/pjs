package models

import (
	"github.com/bashbunni/project-management/utils"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

type Entry struct {
	gorm.Model
	ProjectID uint `gorm:"foreignKey:Project"`
	Project   Project
	Message   string
}

type EntryRepository interface {
	DeleteEntryByID(entryID uint) error
	DeleteEntries(projectID uint) error
	GetEntriesByProjectID(projectID uint) ([]Entry, error)
	CreateEntry(message []byte, projectID uint) error
}

type GormEntryRepository struct {
	DB *gorm.DB
}

func (g *GormEntryRepository) DeleteEntryByID(entryID uint) error {
	result := g.DB.Delete(&Entry{}, entryID)
	resultErr := result.Error
	if resultErr != nil {
		return errors.Wrap(resultErr, utils.CannotDeleteEntry)
	}
	return nil
}

func (g *GormEntryRepository) DeleteEntries(projectID uint) error {
	result := g.DB.Where("project_id = ?", projectID).Delete(&Entry{})
	return result.Error
}

func (g *GormEntryRepository) GetEntriesByProjectID(projectID uint) ([]Entry, error) {
	var Entries []Entry
	result := g.DB.Where("project_id = ?", projectID).Find(&Entries)
	return Entries, result.Error
}

func (g *GormEntryRepository) CreateEntry(message []byte, projectID uint) error {
	entry := Entry{Message: string(message[:]), ProjectID: projectID}
	result := g.DB.Create(&entry)
	return result.Error
}
