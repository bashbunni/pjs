package main

import (
	"fmt"
	"log"

	"github.com/bashbunni/project-management/frontend"
	"github.com/bashbunni/project-management/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// projectPrompt: input validation to create new projects or edit existing
func projectPrompt(pr models.ProjectRepository) models.Project {
	var input int
	pr.PrintProjects()
	fmt.Println("Project ID: ")
	fmt.Scanf("%d", &input)
	// read in input + assign to project
	fmt.Printf("selection is %d \n", input)
	return pr.CreateProject("")
}

func OpenSqlite() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{
		PrepareStmt: true, // caches queries for faster calls
	})
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func main() {
	// setup
	db := OpenSqlite()
	// migrate the schema
	db.AutoMigrate(&models.Entry{}, &models.Project{})
	fmt.Println("entered main")
	frontend.Menu()
	controlSubcommands(db)
}
