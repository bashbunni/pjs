package main

import (
	"fmt"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Project struct {
	gorm.Model
	Id      string
	Message string
}

func (p Project) getMsg() string {
	return p.Message
}

func (p Project) getId() string {
	return p.Id
}

func (p Project) latest() {
	fmt.Printf("%s : %s\n", p.Id, p.Message)
}

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&Project{})
	//	db.Create(&Project{Id: "12345", Message: "hello"})

	var project Project
	var projects []Project
	db.Find(&projects)
	db.First(&project)
	project.latest()

	for _, proj := range projects {
		proj.latest()
	}
}

// https://gorm.io/docs/#Quick-Start
