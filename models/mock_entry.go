package models

import (
	"errors"
)

// Mock Implementation
type MockEntryRepository struct {
	Entries map[uint]*Entry
}

func (m MockEntryRepository) DeleteEntryByID(entryID uint, pe *ProjectWithEntries) error {
	// entryID starts at 1, so we subtract 1 the index
	//	SOFT DELETE
	if _, ok := m.Entries[entryID-1]; ok {
		delete(m.Entries, entryID-1)
	}
	if _, ok := m.Entries[entryID-1]; ok {
		return errors.New("unable to delete entry")
	}
	return nil
}

func (m MockEntryRepository) DeleteEntries(pe *ProjectWithEntries) error {
	m.Entries = make(map[uint]*Entry)
	err := pe.UpdateEntries(m)
	if err != nil {
		return err
	}
	return nil
}

func (m MockEntryRepository) GetEntriesByProjectID(projectID uint) ([]Entry, error) {
	var entries []Entry
	// db IDs start at 1 not 0 therefore also go to one above length of entries map
	for i := 1; i <= len(m.Entries); i++ {
		if m.Entries[uint(i)].ProjectID == projectID {
			entries = append(entries, *m.Entries[uint(i)])
		}
	}
	return entries, nil
}

func (m MockEntryRepository) CreateEntry(message []byte, pe *ProjectWithEntries) error {
	entry := &Entry{Message: string(message[:]), ProjectID: pe.Project.ID}
	err := m.storeEntry(entry, pe)
	return err
}

func (m MockEntryRepository) storeEntry(entry *Entry, pe *ProjectWithEntries) error {
	m.Entries[entry.ID] = entry
	err := pe.UpdateEntries(m)
	return err
}
