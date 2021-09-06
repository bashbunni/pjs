package models

import "gorm.io/gorm"

type ProjectWithEntries struct {
	Project Project
	Entries []Entry
}

func (pe *ProjectWithEntries) CreateProjectWithEntries(projectID int, db *gorm.DB) {
	db.Where("id = ?", projectID).First(pe.Project)
	db.Where("project_id = ?", pe.Project.ID).Find(pe.Entries)
}

func (pe *ProjectWithEntries) UpdateEntries(db *gorm.DB) {
	pe.Entries = GetEntriesByProject(pe.Project.ID, db)
}
