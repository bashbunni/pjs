package main

import (
	"log"

	"github.com/bashbunni/project-management/entry"
	"github.com/bashbunni/project-management/project"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func openSqlite() *gorm.DB {
	db, err := gorm.Open(sqlite.Open("new.db"), &gorm.Config{})
	if err != nil {
		log.Fatalf("unable to open database: %v", err)
	}
	db.AutoMigrate(&entry.Entry{}, &project.Project{})
	return db
}

func main() {
	db := openSqlite()
	controlSubcommands(db)
}
