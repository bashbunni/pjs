package main

import (
	"time"

	"github.com/bashbunni/project-management/models"
	"github.com/bashbunni/project-management/utils"
)

// Mock Implementation
type MockEntryRepository struct {
	Entries map[uint]*models.Entry
}

func (m MockEntryRepository) DeleteEntryByID(entryID uint, pe *models.ProjectWithEntries) {
	// entryID starts at 1, so we subtract 1 the index
	//	SOFT DELETE
	m.Entries[entryID-1].DeletedAt = time.Now()
}

func (m MockEntryRepository) DeleteEntries(pe *models.ProjectWithEntries) {
	m.Entries = make(map[uint]*models.Entry)
	pe.UpdateEntries(m)
}

func (m MockEntryRepository) GetEntriesByProjectID(projectID uint) []models.Entry {
	var entries []models.Entry
	// db IDs start at 1 not 0 therefore also go to one above length of entries map
	for i := 1; i <= len(m.Entries); i++ {
		if m.Entries[uint(i)].ProjectID == projectID {
			entries = append(entries, *m.Entries[uint(i)])
		}
	}
	return entries
}

func (m MockEntryRepository) CreateEntry(pe *models.ProjectWithEntries) {
	message := utils.CaptureInputFromFile()
	entry := &models.Entry{ID: uint(len(m.Entries) + 1), Message: string(message[:]), ProjectID: pe.Project.ID}
	m.storeEntry(entry, pe)
}

func (m MockEntryRepository) storeEntry(entry *models.Entry, pe *models.ProjectWithEntries) {
	m.Entries[entry.ID] = entry
	pe.UpdateEntries(m)
	// TODO: add some kind of validation confirmation for users
}
