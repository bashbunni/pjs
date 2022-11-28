package project

import (
	"log"
	"reflect"
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
	db.AutoMigrate(&models.Project{})
	t.Cleanup(func() {
		db.Migrator().DropTable(&models.Project{})
	})
	return dbconn.Wrap(db)
}

// TestCreateProject

func TestCreateProjectForEmptyDB(t *testing.T) {
	db := Setup(t)
	pr := repos.NewProjectRepo(db)

	pr.CreateProject(&models.Project{Name: "hello"})
	pr.CreateProject(&models.Project{Name: "world"})

	got, _ := pr.GetAllProjects()
	want := []models.Project{{Name: "hello"}, {Name: "world"}}
	for i := range want {
		if got[i].Name != want[i].Name {
			t.Errorf("got %s want %s", got[i].Name, want[i].Name)
		}
	}
}

// TestHasProjects

func TestHasNoProjectsForEmptyDB(t *testing.T) {
	db := Setup(t)
	pr := repos.NewProjectRepo(db)

	got := pr.HasProjects()
	want := false
	if got != want {
		t.Errorf("got %t want %t", got, want)
	}
}

func TestHasTwoProjects(t *testing.T) {
	db := Setup(t)
	pr := repos.NewProjectRepo(db)

	pr.CreateProject(&models.Project{Name: "hello"})
	pr.CreateProject(&models.Project{Name: "world"})

	got := pr.HasProjects()
	want := true
	if got != want {
		t.Errorf("got %t want %t", got, want)
	}
}

// TestGetAllProjects

func TestGetProjectsFromEmptyDB(t *testing.T) {
	db := Setup(t)
	pr := repos.NewProjectRepo(db)

	got, _ := pr.GetAllProjects()
	if len(got) != 0 {
		t.Error("did not get an empty project list")
	}
}

func TestGetTwoProjects(t *testing.T) {
	db := Setup(t)
	pr := repos.NewProjectRepo(db)

	pr.CreateProject(&models.Project{Name: "hello"})
	pr.CreateProject(&models.Project{Name: "world"})

	got, _ := pr.GetAllProjects()
	want := []models.Project{{Name: "hello"}, {Name: "world"}}
	for i := range want {
		if got[i].Name != want[i].Name {
			t.Errorf("got %s want %s", got[i].Name, want[i].Name)
		}
	}
}

// TestGetProjectByID

func TestGetProjectFromEmptyDB(t *testing.T) {
	db := Setup(t)
	pr := repos.NewProjectRepo(db)

	_, err := pr.GetProjectByID(1)
	if err == nil {
		t.Error("expected an error")
	}
}

func TestGetProjectFromNonEmptyDB(t *testing.T) {
	db := Setup(t)
	pr := repos.NewProjectRepo(db)

	pr.CreateProject(&models.Project{Name: "hello"})
	pr.CreateProject(&models.Project{Name: "world"})

	got, err := pr.GetProjectByID(1)
	want := models.Project{Name: "hello"}
	if err != nil || reflect.DeepEqual(got, want) {
		t.Errorf("got %s want %s. err == %v", got.Name, want.Name, err)
	}
}
