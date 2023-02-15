package main

import (
	"fmt"
	"log"

	"github.com/bashbunni/pjs/entry"
	"github.com/bashbunni/pjs/project"
	"github.com/bashbunni/pjs/tui"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func openSqlite() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("new.db"), &gorm.Config{})
	if err != nil {
		return db, fmt.Errorf("unable to open database: %w", err)
	}
	err = db.AutoMigrate(&entry.Entry{}, &project.Project{})
	if err != nil {
		return db, fmt.Errorf("unable to migrate database: %w", err)
	}
	return db, nil
}

func main() {
	db, err := openSqlite()
	if err != nil {
		log.Fatal(err)
	}
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
			log.Fatalf("error creating project: %v", err)
		}
	} else {
		tui.StartTea(pr, er)
	}
}
