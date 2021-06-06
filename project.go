package main

import (
	"fmt"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const Markdown = "markdown"
const Csv = "csv"
const format = "%d : %s\n"

type Entry struct {
	gorm.Model
	ProjectId uint
	Project   Project
	Message   string
}

type Project struct {
	gorm.Model
	Name string
}

func (e Entry) getMsg() string {
	return e.Message
}

func (e Entry) getId() uint {
	return e.ID
}

func printAll(p Project, db *gorm.DB) {
	// should take in an array of entries
	var entries []Entry
	db.Where("project_id = ?", p.ID).Find(&entries) // note to self: queries should be snakecase
	for _, e := range entries {
		fmt.Printf(format, e.getId(), e.getMsg())
	}
}

func (p *Project) saveNewEntry(message string, db *gorm.DB) {
	db.Create(&Entry{Message: message, ProjectId: p.ID})
}

func saveNewProject(name string, db *gorm.DB) Project {
	proj := Project{Name: name}
	db.Create(&proj)
	return proj
}

func main() {
	// setup
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	args := os.Args[:]

	if len(args) <= 1 {
		fmt.Println("Please add a message to commit")
		os.Exit(1)
	}

	message := os.Args[1]

	// migrate the schema
	db.AutoMigrate(&Entry{}, &Project{})

	// other things
	var project Project
	project = saveNewProject("bread's toaster", db)
	project.saveNewEntry(message, db)

	var entries []Entry
	db.Find(&entries) // contains all data from table
	db.First(&entries)

	printAll(project, db)
}

// https://gorm.io/docs/#Quick-Start
