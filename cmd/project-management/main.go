package main

import (
	"fmt"
	"log"

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

func OpenSqlite() (*gorm.DB, error) {
	return gorm.Open(sqlite.Open("new.db"), &gorm.Config{})
}

func main() {
	db, err := OpenSqlite()
	if err != nil {
		log.Fatalf("unable to open database: %v", err)
	}
	db.AutoMigrate(&models.Entry{}, &models.Project{})
	fmt.Println("entered main")
	gp := models.GormProjectRepository{DB: db}
	fmt.Println(gp.GetAllProjects())
	var projects []models.Project
	db.Raw("SELECT * FROM projects").Scan(&projects)
	fmt.Println(projects)
	controlSubcommands(db)
}
