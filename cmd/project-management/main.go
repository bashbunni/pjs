package main

import (
	"fmt"
	"log"

	"github.com/bashbunni/project-management/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// TODO: fix db opening -> no projects showing up

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
	db, err := gorm.Open(sqlite.Open("./test.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func main() {
	db := OpenSqlite()
	fmt.Println(db)
	db.AutoMigrate(&models.Entry{}, &models.Project{})
	fmt.Println("entered main")
	pr := models.GormProjectRepository{DB: db}
	projects, err := pr.GetAllProjects()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(projects)
	controlSubcommands(db)
}
