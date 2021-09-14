package models

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/bashbunni/project-management/utils"
	"gorm.io/gorm"
)

type Entry struct {
	gorm.Model
	ProjectId uint
	Project   Project
	Message   string
}

const divider = "_______________________________________"

func DeleteEntry(pe *ProjectWithEntries, db *gorm.DB) {
	db.Delete(&Entry{}, pe.Project.ID)
	pe.UpdateEntries(db)
}

func DeleteEntries(pe *ProjectWithEntries, db *gorm.DB) {
	db.Where("project_id = ?", pe.Project.ID).Delete(&Entry{})
}

func GetEntriesByProject(projectID uint, db *gorm.DB) []Entry {
	var entries []Entry
	db.Where("project_id = ?", projectID).Find(&entries)
	return entries
}

func CreateEntry(pe *ProjectWithEntries, db *gorm.DB) {
	message := utils.CaptureInputFromFile()
	db.Create(&Entry{Message: string(message[:]), ProjectId: pe.Project.ID})
	pe.UpdateEntries(db)
	fmt.Println(string(message[:]) + " was successfully written to " + pe.Project.Name)
}

// outputs
func formattedOutputFromEntries(entries []Entry) []byte {
	var output string
	for _, entry := range entries {
		output += fmt.Sprintf("ID: %d\nCreated: %s\nMessage:\n %s\n %s\n", entry.ID, entry.CreatedAt.Format("2006-01-02"), entry.Message, divider)
	}
	return []byte(output)
}

func OutputEntriesToMarkdown(entries []Entry) error {
	file, err := os.OpenFile("./output.md", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return err
	}
	defer file.Close() // want defer as close to acquisition of resources as possible
	output := formattedOutputFromEntries(entries)
	_, err = file.Write(output)

	if err != nil {
		return err
	}
	return nil
}

func OutputEntriesToPDF(entries []Entry) error {
	output := formattedOutputFromEntries(entries)              // []byte
	pandoc := exec.Command("pandoc", "-s", "-o", "output.pdf") // c is going to run pandoc, so I'm assigning the pipe to c
	wc, wcerr := pandoc.StdinPipe()                            // io.WriteCloser, err
	if wcerr != nil {
		return wcerr
	}
	goerr := make(chan error)
	done := make(chan bool)
	go func() {
		defer wc.Close()
		_, err := wc.Write(output)
		goerr <- err
		close(goerr)
		close(done)
	}()
	if err := <-goerr; err != nil {
		return err
	}
	err := pandoc.Run()
	if err != nil {
		return err
	}
	return nil
}
