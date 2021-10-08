package models

import (
	"fmt"
	"time"

	"github.com/bashbunni/project-management/utils"
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
	GetEntriesByProjectID(projectID uint) []Entry
	CreateEntry(pe *ProjectWithEntries)
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

func (g GormEntryRepository) GetEntriesByProjectID(projectID uint) []Entry {
	var Entries []Entry
	g.DB.Where("project_id = ?", projectID).Find(&Entries)
	return Entries
}

func (g GormEntryRepository) CreateEntry(pe *ProjectWithEntries) {
	message := utils.CaptureInputFromFile()
	g.DB.Create(&Entry{Message: string(message[:]), ProjectID: pe.Project.ID})
	pe.UpdateEntries(g)

	fmt.Println(string(message[:]) + " was successfully written to " + pe.Project.Name)
}
