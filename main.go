package main

import (
	"flag"
	"fmt"
	"log"
	"time"

	"github.com/bashbunni/project-management/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// mainMenu: flag action handling
func handleFlags(db *gorm.DB) {
	flag.Parse()
	var entries []models.Entry
	db.Find(&entries) // contains all data from table
	if *cEntry != -1 {
		models.CreateEntry(*cEntry, db)
	}
	if *deleteEntry != -1 {
		models.DeleteEntry(*deleteEntry, db)
	}
	if *listProj {
		models.PrintProjects(db)
	}
	if *deleteProj != -1 {
		models.DeleteProject(*deleteProj, db)
	}
	if *editProj != -1 {
		models.RenameProject(*editProj, db)
	}
	if *markdown {
		models.OutputMarkdown(entries)
	}
	if *pdf {
		models.OutputPDF(entries)
	}
	if *start != "" {
		st, errst := time.Parse("2006-01-02", *start)
		if errst != nil {
			log.Fatal(errst)
		}
		en, erren := time.Parse("2006-01-02", *end)
		if erren != nil {
			log.Fatal(erren)
		}
		models.OutputMarkdownByDateRange(st, en, db)
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
