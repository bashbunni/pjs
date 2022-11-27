package main

import (
	"log"

	dbx "github.com/bashbunni/project-management/database"
	"github.com/bashbunni/project-management/database/models"
	"github.com/bashbunni/project-management/database/repos"
	"github.com/bashbunni/project-management/entry"
	"github.com/bashbunni/project-management/project"
	"github.com/bashbunni/project-management/tui"
	"github.com/pkg/errors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func openSqlite() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("new.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("unable to open database: %v", err)
	}
	err = db.AutoMigrate(&entry.Entry{}, &project.Project{})
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func main() {
	if err := dbx.Setup(); err != nil {
		log.Fatalf("unable to create db: %s\n", err)
	}

	db, err := dbx.Connect()
	if err != nil {
		log.Fatalf("unable to connect to db: %s\n", err)
	}

	pr := repos.NewProjectRepo(db)
	er := entry.GormRepository{WDB: db}
	projects, err := pr.GetAllProjects()
	if err != nil {
		log.Fatal(err)
	}
	if len(projects) < 1 {
		name := project.NewProjectPrompt()
		if err := pr.CreateProject(&models.Project{Name: name}); err != nil {
			log.Fatal(errors.Wrap(err, "error creating project"))
		}
	} else {
		tui.StartTea(pr, er)
	}
}
