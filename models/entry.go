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
	DeleteEntryByID(entryID uint, pe *ProjectWithEntries) error
	DeleteEntries(pe *ProjectWithEntries) error
	GetEntriesByProjectID(projectID uint) ([]Entry, error)
	CreateEntry(message []byte, pe *ProjectWithEntries) error
}

type GormEntryRepository struct {
	DB *gorm.DB
}

func (g *GormEntryRepository) DeleteEntryByID(entryID uint, pe *ProjectWithEntries) error {
	result := g.DB.Delete(&Entry{}, entryID)
	resultErr := result.Error
	if resultErr != nil {
		return errors.Wrap(resultErr, utils.CannotDeleteEntry)
	}
	err := pe.UpdateEntries(g)
	if err != nil {
		return err
	}
	return nil
}

func (g *GormEntryRepository) DeleteEntries(pe *ProjectWithEntries) error {
	result := g.DB.Where("project_id = ?", pe.Project.ID).Delete(&Entry{})
	return result.Error
}

func (g *GormEntryRepository) GetEntriesByProjectID(projectID uint) ([]Entry, error) {
	var Entries []Entry
	result := g.DB.Where("project_id = ?", projectID).Find(&Entries)
	return Entries, result.Error
}

func (g *GormEntryRepository) CreateEntry(message []byte, pe *ProjectWithEntries) error {
	entry := Entry{Message: string(message[:]), ProjectID: pe.Project.ID}
	result := g.DB.Create(&entry)
	pe.UpdateEntries(g)
	return result.Error
}
