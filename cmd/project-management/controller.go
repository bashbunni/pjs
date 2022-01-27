package main

import (
	"log"

	"github.com/bashbunni/project-management/frontend"
	"github.com/bashbunni/project-management/models"
	"gorm.io/gorm"
)

func controlSubcommands(db *gorm.DB) {
	pr := models.GormProjectRepository{DB: db}
	projects, err := pr.GetAllProjects()
	if err != nil {
		log.Fatal(err)
	}
	if len(projects) < 1 {
		name := models.NewProjectPrompt()
		pr.CreateProject(name)
	} else {
		frontend.ChooseProject(pr, models.GormEntryRepository{DB: db})
	}
}

 //func hasSubcommands() bool {
 //	if len(os.Args) < 2 {
 //		fmt.Println("expected entry, output, or project subcommands after project ID")
 //		return false
 //	}
 //	return true
 //}
 //
 //func controlEntryCommand(pe *models.ProjectWithEntries, er models.EntryRepository) {
 //	if *createEntry {
 //		message := utils.CaptureInputFromFile()
 //		er.CreateEntry(message, pe)
 //	}
 //	if *deleteEntry != 0 {
 //		er.DeleteEntryByID(*deleteEntry, pe)
 //	}
 //}
 //
 //func controlOutputCommand(entries []models.Entry) {
 //	if *markdown {
 //		err := outputs.OutputEntriesToMarkdown(entries)
 //		if err != nil {
 //			log.Fatal(err)
 //		}
 //	}
 //	if *pdf {
 //		err := outputs.OutputEntriesToPDF(entries)
 //		if err != nil {
 //			log.Fatal(err)
 //		}
 //	}
 //}
 //
 //func controlProjectCommand(pe *models.ProjectWithEntries, pr models.ProjectRepository, er models.EntryRepository) {
 //	if *listAllProjects {
 //		pr.PrintProjects()
 //	}
 //	if *deleteProject {
 //		pr.DeleteProject(pe, er)
 //		os.Exit(0)
 //	}
 //	if *editProject {
 //		pr.RenameProject(pe)
 //	}
 //}
