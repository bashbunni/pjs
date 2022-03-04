package project

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
	db.AutoMigrate(&Project{})
	t.Cleanup(func() {
		db.Migrator().DropTable(&Project{})
	})
	return db
}

// TestCreateProject

func TestCreateProjectForEmptyDB(t *testing.T) {
	db := Setup(t)
	pr := GormRepository{DB: db}

	pr.CreateProject("hello")
	pr.CreateProject("world")

	got, _ := pr.GetAllProjects()
	want := []Project{{Name: "hello"}, {Name: "world"}}
	for i := range want {
		if got[i].Name != want[i].Name {
			t.Errorf("got %s want %s", got[i].Name, want[i].Name)
		}
	}
}

// TestHasProjects

func TestHasNoProjectsForEmptyDB(t *testing.T) {
	db := Setup(t)
	pr := GormRepository{DB: db}

	got := pr.HasProjects()
	want := false
	if got != want {
		t.Errorf("got %t want %t", got, want)
	}
}

func TestHasTwoProjects(t *testing.T) {
	db := Setup(t)
	pr := GormRepository{DB: db}

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
	pr := GormRepository{DB: db}

	got, _ := pr.GetAllProjects()
	if len(got) != 0 {
		t.Error("did not get an empty project list")
	}
}

func TestGetTwoProjects(t *testing.T) {
	db := Setup(t)
	pr := GormRepository{DB: db}

	pr.CreateProject("hello")
	pr.CreateProject("world")

	got, _ := pr.GetAllProjects()
	want := []Project{{Name: "hello"}, {Name: "world"}}
	for i := range want {
		if got[i].Name != want[i].Name {
			t.Errorf("got %s want %s", got[i].Name, want[i].Name)
		}
	}
}

// TestGetProjectByID

func TestGetProjectFromEmptyDB(t *testing.T) {
	db := Setup(t)
	pr := GormRepository{DB: db}

	_, err := pr.GetProjectByID(1)
	if err == nil {
		t.Error("expected an error")
	}
}

func TestGetProjectFromNonEmptyDB(t *testing.T) {
	db := Setup(t)
	pr := GormRepository{DB: db}

	pr.CreateProject("hello")
	pr.CreateProject("world")

	got, err := pr.GetProjectByID(1)
	want := Project{Name: "hello"}
	if err != nil || reflect.DeepEqual(got, want) {
		t.Errorf("got %s want %s. err == %v", got.Name, want.Name, err)
	}
}
