package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type Project struct {
	gorm.Model
	Id      string
	Message string
}

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&Project{})
	db.Create(&Project{Id: "12345", Message: "hello"})

	var project Project
	db.First(&project, 1)
	db.First(&project, "code =  ?", "D42")

	db.Model(&project).Update("Message", "Hello world")
	db.Model(&project).Updates(Project{Id: "12345", Message: "hello"})
}

// https://gorm.io/docs/#Quick-Start
