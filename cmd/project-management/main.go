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
	newproject, err := pr.CreateProject("")
	if err != nil {
		log.Fatal(err)
	}
	return newproject
}

func OpenSqlite() *gorm.DB {
	db, err :=  gorm.Open(sqlite.Open("new.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("unable to open database: %v", err)
	}
	db.AutoMigrate(&models.Entry{}, &models.Project{})
	return db
}

func main() {
	db := OpenSqlite()
	fmt.Println("entered main")
	gp := models.GormProjectRepository{DB: db}
	fmt.Println(gp.GetAllProjects())
	var projects []models.Project
	db.Raw("SELECT * FROM projects").Scan(&projects)
	fmt.Println(projects)
	controlSubcommands(db)
}
