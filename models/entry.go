package models

import (
	"time"

	"gorm.io/gorm"
)

type Entry struct {
	gorm.Model
	ID        uint
	ProjectID uint // TODO: get rid of duplicate data
	Project   Project
	Message   string
	DeletedAt time.Time
}

type EntryRepository interface {
	DeleteEntryByID(entryID uint, pe *ProjectWithEntries)
	DeleteEntries(pe *ProjectWithEntries)
	GetEntriesByProjectID(projectID uint) ([]Entry, error)
	CreateEntry(message []byte, pe *ProjectWithEntries) error
}

type GormEntryRepository struct {
	DB *gorm.DB
}

func (g GormEntryRepository) DeleteEntryByID(entryID uint, pe *ProjectWithEntries) {
	g.DB.Delete(&Entry{}, entryID)
	pe.UpdateEntries(g)
}

func (g GormEntryRepository) DeleteEntries(pe *ProjectWithEntries) {
	g.DB.Where("project_id = ?", pe.Project.ID).Delete(&Entry{})
}

func (g GormEntryRepository) GetEntriesByProjectID(projectID uint) ([]Entry, error) {
	var Entries []Entry
	result := g.DB.Where("project_id = ?", projectID).Find(&Entries)
	return Entries, result.Error
}

func (g GormEntryRepository) CreateEntry(message []byte, pe *ProjectWithEntries) error {
	entry := Entry{Message: string(message[:]), ProjectID: pe.Project.ID}
	result := g.DB.Create(&entry)
	pe.UpdateEntries(g)
	return result.Error
}
