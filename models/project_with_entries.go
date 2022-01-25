package models

import (
	"github.com/bashbunni/project-management/utils"
	"github.com/pkg/errors"
)

type ProjectWithEntries struct {
	Project Project
	Entries []Entry
}

// TODO: does GetOrCreateProject provide sufficient validation?

// getters + setters
func (pe ProjectWithEntries) GetEntries() []Entry {
	return pe.Entries
}

// functions
func CreateProjectWithEntries(project Project, er EntryRepository) (*ProjectWithEntries, error) {
	entries, err := er.GetEntriesByProjectID(project.ID)
	if err != nil {
		return &ProjectWithEntries{}, errors.Wrap(err, utils.CannotCreateProjectWithEntries)
	}
	return &ProjectWithEntries{project, entries}, nil
}

