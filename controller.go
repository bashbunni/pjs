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
	pr := models.GormProjectRepository{DB: db}
	project := parseProjectID(os.Args[1], pr)
	er := models.GormEntryRepository{DB: db}
	pe := models.CreateProjectWithEntries(project, er)
	switch os.Args[2] {
	case "entry":
		entryCommands.Parse(os.Args[3:])
		controlEntryCommand(pe, er)
	case "output":
		outputCommands.Parse(os.Args[3:])
		controlOutputCommand(pe.GetEntries())
	case "project":
		projectCommands.Parse(os.Args[3:])
		controlProjectCommand(pe, pr, er)
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

func controlEntryCommand(pe *models.ProjectWithEntries, er models.EntryRepository) {
	if *createEntry {
		er.CreateEntry(pe)
	}
	if *deleteEntry != 0 {
		er.DeleteEntryByID(*deleteEntry, pe)
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

func controlProjectCommand(pe *models.ProjectWithEntries, pr models.ProjectRepository, er models.EntryRepository) {
	if *listAllProjects {
		pr.PrintProjects()
	}
	if *deleteProject {
		pr.DeleteProject(pe, er)
		os.Exit(0)
	}
	if *editProject {
		pr.RenameProject(pe)
	}
}

func parseProjectID(input string, pr models.ProjectRepository) models.Project {
	projectID, err := strconv.Atoi(input)
	if err != nil {
		log.Fatal(err)
	}
	return pr.GetOrCreateProjectByID(projectID)
}
