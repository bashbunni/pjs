package entry

import (
	"log"
	"testing"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Setup(t *testing.T) *gorm.DB {
	t.Helper() // allows me to log Gorm errors later
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		log.Fatalf("unable to open in-memory SQLite DB: %v", err)
	}
	db.AutoMigrate(&Model{})
	t.Cleanup(func() {
		db.Migrator().DropTable(&Model{})
	})
	return db
}

// DeleteEntryByID
func TestDeleteEntryForEmptyDB(t *testing.T) {
	db := Setup(t)
	er := GormRepository{DB: db}

	er.DeleteEntryByID(1)
	if err := db.Unscoped().Where("ID = 1").First(&Model{}).Error; err == nil {
		t.Error("expected error")
	}
}

func TestDeleteEntryWithTwoEntries(t *testing.T) {
	db := Setup(t)
	er := GormRepository{DB: db}

	er.CreateEntry([]byte("hello world"), 1)
	er.CreateEntry([]byte("I am just a world"), 1)

	er.DeleteEntryByID(1)
	if err := db.Unscoped().Where("ID = 1").First(&Model{}).Error; err != nil {
		t.Error("expected no error")
	}
}

// DeleteEntries
func TestDeleteEntriesForEmptyDB(t *testing.T) {
	db := Setup(t)
	er := GormRepository{DB: db}

	er.DeleteEntries(1)
	if err := db.Unscoped().Where("ID = 1").First(&Model{}).Error; err == nil {
		t.Error("expected error")
	}
}

func TestDeleteEntriesWithTwoEntries(t *testing.T) {
	db := Setup(t)
	er := GormRepository{DB: db}

	er.CreateEntry([]byte("hello world"), 1)
	er.CreateEntry([]byte("I am just a world"), 1)

	er.DeleteEntries(1)
	if err := db.Unscoped().Where("ID = 1").First(&Model{}).Error; err != nil {
		t.Error("expected no error")
	}
}

// GetEntriesByProjectID
func TestGetEntriesByProjectIDForEmptyDB(t *testing.T) {
	db := Setup(t)
	er := GormRepository{DB: db}

	got, _ := er.GetEntriesByProjectID(1)
	if len(got) != 0 {
		t.Error("expected an empty list of entries")
	}
}

func TestGetEntriesByProjectIDWithTwoEntries(t *testing.T) {
	db := Setup(t)
	er := GormRepository{DB: db}

	er.CreateEntry([]byte("hello world"), 1)
	er.CreateEntry([]byte("I am just a world"), 1)

	got, _ := er.GetEntriesByProjectID(1)
	if len(got) == 0 {
		t.Error("expected a list with entries")
	}
}

// CreateEntry -> covered in previous tests
