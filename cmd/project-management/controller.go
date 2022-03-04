package main

import (
	"log"
	"os"

	"github.com/bashbunni/project-management/entry"
	"github.com/bashbunni/project-management/frontend"
	"github.com/bashbunni/project-management/project"
	"gorm.io/gorm"
)

func controlSubcommands(db *gorm.DB) {
	pr := project.GormRepository{DB: db}
	projects, err := pr.GetAllProjects()
	if err != nil {
		log.Fatal(err)
	}
	if len(projects) < 1 {
		name := project.NewProjectPrompt()
		pr.CreateProject(name)
	} else {
		frontend.StartTea(pr, entry.GormRepository{DB: db})
	}
}

func hasSubcommands() bool {
	if len(os.Args) < 2 {
		return false
	}
	return true
}
