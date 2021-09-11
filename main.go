package main

import (
	"fmt"
	"os"

	"github.com/bashbunni/project-management/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

/*
TODO: instead of asking which project each time, the user will need to choose a project to work in
why? minimize queries on db -> query for entries once and again when a new entry is added
*/

func handleSubcommands(db *gorm.DB) {
	if len(os.Args) < 2 {
		fmt.Println("expected entry, output, or project subcommands")
		os.Exit(1)
	}

	var entries []models.Entry
	db.Find(&entries) // contains all data from table

	switch os.Args[1] {
	case "entry":
		entryCommands.Parse(os.Args[2:])
		handleEntryCommand(entries, db)
	case "output":
		outputCommands.Parse(os.Args[2:])
		handleOutputCommand(entries)
	case "project":
		projectCommands.Parse(os.Args[2:])
		handleProjectCommand(db)
	}
}

func handleEntryCommand(entries []models.Entry, db *gorm.DB) {
	if *createEntry {
		models.CreateEntry(db)
	}
	if *deleteEntry {
		models.DeleteEntry(db)
	}
}

func handleOutputCommand(entries []models.Entry) {
	if *markdown {
		models.OutputMarkdown(entries)
	}
	if *pdf {
		models.OutputPDF(entries)
	}
}

func handleProjectCommand(db *gorm.DB) {
	if *listAllProjects {
		models.PrintProjects(db)
	}
	if *deleteProject {
		models.DeleteProject(db)
	}
	if *editProject {
		models.RenameProject(db)
	}
}

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
	handleFlags(db)
}
