package main

import (
	"log"

	"github.com/bashbunni/pjs/entry"
	"github.com/bashbunni/pjs/project"
	"github.com/bashbunni/pjs/tui"
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
	db := openSqlite()
	pr := project.GormRepository{DB: db}
	er := entry.GormRepository{DB: db}
	projects, err := pr.GetAllProjects()
	if err != nil {
		log.Fatal(err)
	}
	if len(projects) < 1 {
		name := project.NewProjectPrompt()
		_, err := pr.CreateProject(name)
		if err != nil {
			log.Fatal(errors.Wrap(err, "error creating project"))
		}
	} else {
		tui.StartTea(pr, er)
	}
}
