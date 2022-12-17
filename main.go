package main

import (
	"fmt"
	"log"

	"github.com/bashbunni/project-management/database/dbconn"
	"github.com/bashbunni/project-management/database/models"
	"github.com/bashbunni/project-management/database/repos"
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
	err = db.AutoMigrate(&models.Entry{}, &models.Project{})
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func openWrappedDB() (dbconn.GormWrapper, error) {
	db, err := gorm.Open(sqlite.Open("new.db"), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to open db: %w", err)
	}

	wdb := dbconn.Wrap(db)
	if err := models.AutoMigrate(wdb); err != nil {
		return nil, err
	}
	return wdb, nil
}

func main() {
	// db, err := openWrappedDB()
	// if err != nil {
	// 	log.Fatalf("unable to connect to db: %s\n", err)
	// }

	db := dbconn.Wrap(openSqlite())

	pr := repos.NewProjectRepo(db)
	er := repos.NewEntryRepo(db)
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
