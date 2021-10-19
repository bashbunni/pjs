package mocks

import (
	"errors"
	"time"

	"github.com/bashbunni/project-management/models"
)

// Mock Implementation
type MockEntryRepository struct {
	Entries map[uint]*models.Entry
}

func (m MockEntryRepository) DeleteEntryByID(entryID uint, pe *models.ProjectWithEntries) error {
	// entryID starts at 1, so we subtract 1 the index
	//	SOFT DELETE
	m.Entries[entryID-1].DeletedAt = time.Now()
	if m.Entries[entryID-1].DeletedAt == nil {
		// TODO: finish this function
		return errors.New("unable to delete entry")
	}
	return nil
}

func (m MockEntryRepository) DeleteEntries(pe *models.ProjectWithEntries) error {
	m.Entries = make(map[uint]*models.Entry)
	err := pe.UpdateEntries(m)
	if err != nil {
		return err
	}
	return nil
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

func (m MockEntryRepository) CreateEntry(message []byte, pe *models.ProjectWithEntries) error {
	entry := &models.Entry{ID: uint(len(m.Entries) + 1), Message: string(message[:]), ProjectID: pe.Project.ID}
	err := m.storeEntry(entry, pe)
	return err
}

func (m MockEntryRepository) storeEntry(entry *models.Entry, pe *models.ProjectWithEntries) error {
	m.Entries[entry.ID] = entry
	err := pe.UpdateEntries(m)
	return err
}
