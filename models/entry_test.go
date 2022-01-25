package models

import "testing"

// DeleteEntryByID
func TestDeleteEntryForEmptyDB(t *testing.T) {
	db := Setup(t)
	er := GormEntryRepository{DB: db}

	er.DeleteEntryByID(1)
	if err := db.Unscoped().Where("ID = 1").First(&Entry{}).Error; err == nil {
		t.Error("expected error")
	}
}

func TestDeleteEntryWithTwoEntries(t *testing.T) {
	db := Setup(t)
	pr := GormProjectRepository{DB: db}
	project, _ := pr.CreateProject("project1")
	er := GormEntryRepository{DB: db}

	er.CreateEntry([]byte("hello world"), project.ID)
	er.CreateEntry([]byte("I am just a world"), project.ID)

	er.DeleteEntryByID(1)
	if err := db.Unscoped().Where("ID = 1").First(&Entry{}).Error; err != nil {
		t.Error("expected no error")
	}
}

// DeleteEntries
func TestDeleteEntriesForEmptyDB(t *testing.T) {
	db := Setup(t)
	er := GormEntryRepository{DB: db}

	er.DeleteEntries(1)
	if err := db.Unscoped().Where("ID = 1").First(&Entry{}).Error; err == nil {
		t.Error("expected error")
	}
}

func TestDeleteEntriesWithTwoEntries(t *testing.T) {
	db := Setup(t)
	pr := GormProjectRepository{DB: db}
	project, _ := pr.CreateProject("project1")
	er := GormEntryRepository{DB: db}

	er.CreateEntry([]byte("hello world"), project.ID)
	er.CreateEntry([]byte("I am just a world"), project.ID)

	er.DeleteEntries(project.ID)
	if err := db.Unscoped().Where("ID = 1").First(&Entry{}).Error; err != nil {
		t.Error("expected no error")
	}
}

// GetEntriesByProjectID
func TestGetEntriesByProjectIDForEmptyDB(t *testing.T) {
	db := Setup(t)
	er := GormEntryRepository{DB: db}

	got, _ := er.GetEntriesByProjectID(1)
	if len(got) != 0 {
		t.Error("expected an empty list of entries")
	}
}

func TestGetEntriesByProjectIDWithTwoEntries(t *testing.T) {
	db := Setup(t)
	pr := GormProjectRepository{DB: db}
	project, _ := pr.CreateProject("project1")
	er := GormEntryRepository{DB: db}

	er.CreateEntry([]byte("hello world"), project.ID)
	er.CreateEntry([]byte("I am just a world"), project.ID)

	got, _ := er.GetEntriesByProjectID(1)
	if len(got) == 0 {
		t.Error("expected a list with entries")
	}
}

// CreateEntry -> covered in previous tests
