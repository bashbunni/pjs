package models

import (
	"log"
	"reflect"
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
	db.AutoMigrate(&Entry{}, &Project{})
	t.Cleanup(func() {
		db.Migrator().DropTable(&Entry{}, &Project{})
	})
	return db
}


// TestCreateProject

func TestCreateProjectForEmptyDB(t *testing.T) {
	db := Setup(t)
	pr := GormProjectRepository{DB: db}

	pr.CreateProject("hello")
	pr.CreateProject("world")

	got := pr.GetAllProjects()
	want := []Project{{Name: "hello"}, {Name: "world"}}
	if reflect.DeepEqual(got, want) {
		t.Error("did not get correct project list")
	}
}

// TestHasProjects

func TestHasNoProjectsForEmptyDB(t *testing.T) {
	db := Setup(t)
	pr := GormProjectRepository{DB: db}

	got := pr.HasProjects()
	want := false
	if got != want {
		t.Errorf("got %t want %t", got, want)
	}
}

func TestHasTwoProjects(t *testing.T) {
	db := Setup(t)
	pr := GormProjectRepository{DB: db}

	pr.CreateProject("hello")
	pr.CreateProject("world")

	got := pr.HasProjects()
	want := true
	if got != want {
		t.Errorf("got %t want %t", got, want)
	}
}

// TestGetAllProjects

func TestGetProjectsFromEmptyDB(t *testing.T) {
	db := Setup(t)
	pr := GormProjectRepository{DB: db}

	got := pr.GetAllProjects()
	want := []Project{}
	if reflect.DeepEqual(got, want) {
		t.Error("did not get an empty project list")
	}
}

func TestGetTwoProjects(t *testing.T) {
	db := Setup(t)
	pr := GormProjectRepository{DB: db}

	pr.CreateProject("hello")
	pr.CreateProject("world")

	got := pr.GetAllProjects()
	want := []Project{{Name: "hello"}, {Name: "world"}}
	if reflect.DeepEqual(got, want) {
		t.Error("did not get correct project list")
	}
}

// TestGetProjectByID

func TestGetProjectFromEmptyDB(t *testing.T) {
	db := Setup(t)
	pr := GormProjectRepository{DB: db}

	pr.GetProjectByID(1)
}

func TestGetProjectFromNonEmptyDB(t *testing.T) {
	db := Setup(t)
	pr := GormProjectRepository{DB: db}

	pr.CreateProject("hello")
	pr.CreateProject("world")

	pr.GetProjectByID(1)
}

