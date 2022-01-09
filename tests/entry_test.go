package tests

/*
TODO:
Write tests for:
type EntryRepository interface {
	DeleteEntryByID(entryID uint, pe *ProjectWithEntries) error
	DeleteEntries(pe *ProjectWithEntries) error
	GetEntriesByProjectID(projectID uint) ([]Entry, error)
	CreateEntry(message []byte, pe *ProjectWithEntries) error
}
*/

mockProject := Project{1, "first project", nil}

mockData := models.MockEntryRepository{
	map[uint]*models.Entry{
		{1: &models.Entry{1, "I'm an entry"},  ProjectWithEntries{mockProject, []Entry{
	{},
		}}}}
	}
}
