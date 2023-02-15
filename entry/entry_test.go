package entry

import (
	"fmt"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Setup(t *testing.T) (*gorm.DB, error) {
	t.Helper() // allows me to log Gorm errors later
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		return db, fmt.Errorf("unable to open in-memory SQLite DB: %w", err)
	}
	db.AutoMigrate(&Entry{})
	t.Cleanup(func() {
		db.Migrator().DropTable(&Entry{})
	})
	return db, nil
}

// DeleteEntryByID
func TestDeleteEntryForEmptyDB(t *testing.T) {
	db, err := Setup(t)
	if err != nil {
		t.Fatal(err)
	}
	er := GormRepository{DB: db}

	er.DeleteEntryByID(1)
	if err := db.Unscoped().Where("ID = 1").First(&Entry{}).Error; err == nil {
		t.Error("expected error")
	}
}

func TestDeleteEntryWithTwoEntries(t *testing.T) {
	db, err := Setup(t)
	if err != nil {
		t.Fatal(err)
	}
	er := GormRepository{DB: db}

	er.CreateEntry([]byte("hello world"), 1)
	er.CreateEntry([]byte("I am just a world"), 1)

	er.DeleteEntryByID(1)
	if err := db.Unscoped().Where("ID = 1").First(&Entry{}).Error; err != nil {
		t.Error("expected no error")
	}
}

// DeleteEntries
func TestDeleteEntriesForEmptyDB(t *testing.T) {
	db, err := Setup(t)
	if err != nil {
		t.Fatal(err)
	}

	er := GormRepository{DB: db}

	er.DeleteEntries(1)
	if err := db.Unscoped().Where("ID = 1").First(&Entry{}).Error; err == nil {
		t.Error("expected error")
	}
}

func TestDeleteEntriesWithTwoEntries(t *testing.T) {
	db, err := Setup(t)
	if err != nil {
		t.Fatal(err)
	}

	er := GormRepository{DB: db}

	er.CreateEntry([]byte("hello world"), 1)
	er.CreateEntry([]byte("I am just a world"), 1)

	er.DeleteEntries(1)
	if err := db.Unscoped().Where("ID = 1").First(&Entry{}).Error; err != nil {
		t.Error("expected no error")
	}
}

// GetEntriesByProjectID
func TestGetEntriesByProjectIDForEmptyDB(t *testing.T) {
	db, err := Setup(t)
	if err != nil {
		t.Fatal(err)
	}

	er := GormRepository{DB: db}

	got, _ := er.GetEntriesByProjectID(1)
	if len(got) != 0 {
		t.Error("expected an empty list of entries")
	}
}

func TestGetEntriesByProjectIDWithTwoEntries(t *testing.T) {
	db, err := Setup(t)
	if err != nil {
		t.Fatal(err)
	}

	er := GormRepository{DB: db}

	er.CreateEntry([]byte("hello world"), 1)
	er.CreateEntry([]byte("I am just a world"), 1)

	got, _ := er.GetEntriesByProjectID(1)
	if len(got) == 0 {
		t.Error("expected a list with entries")
	}
}

// CreateEntry -> covered in previous tests
