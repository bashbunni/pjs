package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/bashbunni/project-management/models"
	"gorm.io/gorm"
)

func controlSubcommands(db *gorm.DB) *models.ProjectWithEntries {
	if !hasSubcommands() {
		log.Fatal("no subcommands given")
	}
	project := parseProjectID(os.Args[1], db)
	pe := models.CreateProjectWithEntries(project, db)
	switch os.Args[2] {
	case "entry":
		entryCommands.Parse(os.Args[3:])
		controlEntryCommand(pe, db)
	case "output":
		outputCommands.Parse(os.Args[3:])
		controlOutputCommand(pe.GetEntries())
	case "project":
		projectCommands.Parse(os.Args[3:])
		controlProjectCommand(pe, db)
	default:
		fmt.Println("entry, output, or project subcommand not given")
		os.Exit(2)
	}
	return pe
}

func hasSubcommands() bool {
	if len(os.Args) < 2 {
		fmt.Println("expected entry, output, or project subcommands after project ID")
		return false
	}
	return true
}

func controlEntryCommand(pe *models.ProjectWithEntries, db *gorm.DB) {
	if *createEntry {
		models.CreateEntry(pe, db)
	}
	if *deleteEntry {
		models.DeleteEntry(pe, db)
	}
}

func controlOutputCommand(entries []models.Entry) {
	if *markdown {
		err := models.OutputEntriesToMarkdown(entries)
		if err != nil {
			log.Fatal(err)
		}
	}
	if *pdf {
		err := models.OutputEntriesToPDF(entries)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func controlProjectCommand(pe *models.ProjectWithEntries, db *gorm.DB) {
	if *listAllProjects {
		models.PrintProjects(db)
	}
	if *deleteProject {
		models.DeleteProject(pe, db)
		os.Exit(0)
	}
	if *editProject {
		models.RenameProject(pe, db)
	}
}

func parseProjectID(input string, db *gorm.DB) models.Project {
	projectID, err := strconv.Atoi(input)
	if err != nil {
		log.Fatal(err)
	}
	return models.GetOrCreateProjectByID(projectID, db)
}



