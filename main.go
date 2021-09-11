package main

import (
	"fmt"

	"github.com/bashbunni/project-management/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// projectPrompt: input validation to create new projects or edit existing
func projectPrompt(db *gorm.DB) models.Project {
	var input int
	models.PrintProjects(db)
	fmt.Println("Project ID: ")
	fmt.Scanf("%d", &input)
	// read in input + assign to project
	fmt.Printf("selection is %d \n", input)
	return models.CreateProject("", db)
}

func OpenSqlite() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{
		PrepareStmt: true, // caches queries for faster calls
	})
	if err != nil {
		panic("failed to connect database")
	}
	return db
}

func main() {
	// setup
	db := OpenSqlite()
	// migrate the schema
	db.AutoMigrate(&models.Entry{}, &models.Project{})
	controlSubcommands(db)
}
