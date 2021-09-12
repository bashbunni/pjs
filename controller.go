package main

import (
	"fmt"
	"os"
	"strconv"

	"github.com/bashbunni/project-management/models"
	"gorm.io/gorm"
)

func hasSubcommands() bool {
	if len(os.Args) < 3 {
		fmt.Println("expected entry, output, or project subcommands after project ID")
		return false
	}
	return true
}

func parseProjectID(input string, db *gorm.DB) models.Project {
	projectID, err := strconv.Atoi(input)
	if err != nil {
		fmt.Errorf("unable to convert projectID to int")
	}
	return models.GetOrCreateProject(projectID, db)
}

func controlEntryCommand(entries []models.Entry, db *gorm.DB) {
	if *createEntry {
		models.CreateEntry(db)
	}
	if *deleteEntry {
		models.DeleteEntry(db)
	}
}

func controlOutputCommand(entries []models.Entry) {
	if *markdown {
		models.OutputMarkdown(entries)
	}
	if *pdf {
		models.OutputPDF(entries)
	}
}

func controlProjectCommand(db *gorm.DB) {
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

func controlSubcommands(db *gorm.DB) *models.ProjectWithEntries {
	if !hasSubcommands() {
		os.Exit(1)
	}

	project := parseProjectID(os.Args[1], db)
	thisProject := models.CreateProjectWithEntries(project, db)

	switch os.Args[3] {
	case "entry":
		entryCommands.Parse(os.Args[2:])
		controlEntryCommand(thisProject.GetEntries(), db)
	case "output":
		outputCommands.Parse(os.Args[2:])
		controlOutputCommand(thisProject.GetEntries())
	case "project":
		projectCommands.Parse(os.Args[2:])
		controlProjectCommand(db)
	}
	return thisProject
}
