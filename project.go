package main

import (
	"fmt"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Project struct {
	gorm.Model
	Id      uint32
	Message string
}

func (p Project) getMsg() string {
	return p.Message
}

func (p Project) getId() uint32 {
	return p.Id
}

func (p Project) latest() {
	fmt.Printf("%d : %s\n", p.getId(), p.getMsg())
}

func (p *Project) add(message string, db *gorm.DB) {
	db.Create(&Project{Message: message})
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
	db.AutoMigrate(&Project{})

	// other things
	var project Project

	project.add(message, db)

	var projects []Project
	db.Find(&projects)
	db.First(&project)
	project.latest()

	for _, proj := range projects {
		proj.latest()
	}
}

// https://gorm.io/docs/#Quick-Start
