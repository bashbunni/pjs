package models

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
func CreateProjectWithEntries(project Project, er EntryRepository) *ProjectWithEntries {
	entries := er.GetEntriesByProjectID(project.ID)
	return &ProjectWithEntries{project, entries}
}

func (pe *ProjectWithEntries) UpdateEntries(er EntryRepository) {
	pe.Entries = er.GetEntriesByProjectID(pe.Project.ID)
}
