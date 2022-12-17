package entry

import (
	"log"
	"testing"

	"github.com/bashbunni/project-management/database/dbconn"
	"github.com/bashbunni/project-management/database/models"
	"github.com/bashbunni/project-management/database/repos"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func Setup(t *testing.T) dbconn.GormWrapper {
	t.Helper() // allows me to log Gorm errors later
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{})
	if err != nil {
		log.Fatalf("unable to open in-memory SQLite DB: %v", err)
	}
	db.AutoMigrate(&models.Entry{})
	t.Cleanup(func() {
		db.Migrator().DropTable(&models.Entry{})
	})
	return dbconn.Wrap(db)
}

// DeleteEntryByID
func TestDeleteEntryForEmptyDB(t *testing.T) {
	db := Setup(t)
	er := repos.NewEntryRepo(db)

	er.DeleteEntryByID(1)
	if err := db.Unscoped().Where("ID = 1").First(&models.Entry{}).Error(); err == nil {
		t.Error("expected error")
	}
}

func TestDeleteEntryWithTwoEntries(t *testing.T) {
	db := Setup(t)
	er := repos.NewEntryRepo(db)

	er.CreateEntryInProjectOfID([]byte("hello world"), 1)
	er.CreateEntryInProjectOfID([]byte("I am just a world"), 1)

	er.DeleteEntryByID(1)
	if err := db.Unscoped().Where("ID = 1").First(&models.Entry{}).Error(); err != nil {
		t.Error("expected no error")
	}
}

// DeleteEntries
func TestDeleteEntriesForEmptyDB(t *testing.T) {
	db := Setup(t)
	er := repos.NewEntryRepo(db)

	er.DeleteEntriesInProjectOfID(1)
	if err := db.Unscoped().Where("ID = 1").First(&models.Entry{}).Error(); err == nil {
		t.Error("expected error")
	}
}

func TestDeleteEntriesWithTwoEntries(t *testing.T) {
	db := Setup(t)
	er := repos.NewEntryRepo(db)

	er.CreateEntryInProjectOfID([]byte("hello world"), 1)
	er.CreateEntryInProjectOfID([]byte("I am just a world"), 1)

	er.DeleteEntriesInProjectOfID(1)
	if err := db.Unscoped().Where("ID = 1").First(&models.Entry{}).Error(); err != nil {
		t.Error("expected no error")
	}
}

// GetEntriesByProjectID
func TestGetEntriesByProjectIDForEmptyDB(t *testing.T) {
	db := Setup(t)
	er := repos.NewEntryRepo(db)

	got, _ := er.GetEntriesByProjectID(1)
	if len(got) != 0 {
		t.Error("expected an empty list of entries")
	}
}

func TestGetEntriesByProjectIDWithTwoEntries(t *testing.T) {
	db := Setup(t)
	er := repos.NewEntryRepo(db)

	er.CreateEntryInProjectOfID([]byte("hello world"), 1)
	er.CreateEntryInProjectOfID([]byte("I am just a world"), 1)

	got, _ := er.GetEntriesByProjectID(1)
	if len(got) == 0 {
		t.Error("expected a list with entries")
	}
}

// CreateEntry -> covered in previous tests
