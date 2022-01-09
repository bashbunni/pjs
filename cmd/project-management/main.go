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

func OpenSqlite() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("./test.db"), &gorm.Config{})
	if err != nil {
		log.Fatal(err)
	}
	return db
}

func main() {
	db := OpenSqlite()
	db.AutoMigrate(&models.Entry{}, &models.Project{})
	fmt.Println("entered main")
	gp := models.GormProjectRepository{DB: db}
	fmt.Println(gp.GetAllProjects())
	gp.CreateProject("chat")
	controlSubcommands(db)
}
