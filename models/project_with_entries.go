package models

import "gorm.io/gorm"

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

func CreateProjectWithEntries(project Project, db *gorm.DB) *ProjectWithEntries {
	var entries []Entry
	db.Where("project_id = ?", project.ID).Find(&entries)
	return &ProjectWithEntries{project, entries}
}

func (pe *ProjectWithEntries) UpdateEntries(db *gorm.DB) {
	pe.Entries = GetEntriesByProject(pe.Project.ID, db)
}
