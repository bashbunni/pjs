package main

import (
	"fmt"
	"log"

	"github.com/bashbunni/project-management/entry"
	"github.com/bashbunni/project-management/project"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// projectPrompt: input validation to create new projects or edit existing
func projectPrompt(pr project.Repository) project.Project {
	var input int
	pr.PrintProjects()
	fmt.Println("Project ID: ")
	fmt.Scanf("%d", &input)
	// read in input + assign to project
	fmt.Printf("selection is %d \n", input)
	newproject, err := pr.CreateProject("")
	if err != nil {
		log.Fatal(err)
	}
	return newproject
}

func openSqlite() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("new.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("unable to open database: %v", err)
	}
	db.AutoMigrate(&entry.Entry{}, &project.Project{})
	return db
}

func main() {
	db := openSqlite()
	controlSubcommands(db)
}
