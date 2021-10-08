package models

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/bashbunni/project-management/utils"
	"github.com/pkg/errors"
	"gorm.io/gorm"
)

// Entity
type Entry struct {
	gorm.Model
	ID        uint
	ProjectID uint // TODO: get rid of duplicate data
	Project   Project
	Message   string
	DeletedAt time.Time
}

// Interface
type EntryRepository interface {
	DeleteEntryByID(entryID uint, pe *ProjectWithEntries)
	DeleteEntries(pe *ProjectWithEntries)
	GetEntriesByProjectID(projectID uint) []Entry
	CreateEntry(pe *ProjectWithEntries)
}

// TODO: make this not trash
const divider = "_______________________________________"

// Gorm implementation

type GormEntryRepository struct {
	DB *gorm.DB
}

func (g GormEntryRepository) DeleteEntryByID(entryID uint, pe *ProjectWithEntries) {
	g.DB.Delete(&Entry{}, entryID)
	pe.UpdateEntries(g)
}

func (g GormEntryRepository) DeleteEntries(pe *ProjectWithEntries) {
	g.DB.Where("project_id = ?", pe.Project.ID).Delete(&Entry{})
}

func (g GormEntryRepository) GetEntriesByProjectID(projectID uint) []Entry {
	var Entries []Entry
	g.DB.Where("project_id = ?", projectID).Find(&Entries)
	return Entries
}

func (g GormEntryRepository) CreateEntry(pe *ProjectWithEntries) {
	message := utils.CaptureInputFromFile()
	g.DB.Create(&Entry{Message: string(message[:]), ProjectID: pe.Project.ID})
	pe.UpdateEntries(g)

	fmt.Println(string(message[:]) + " was successfully written to " + pe.Project.Name)
}

// outputs
func formattedOutputFromEntries(Entries []Entry) []byte {
	var output string
	for _, entry := range Entries {
		output += fmt.Sprintf("ID: %d\nCreated: %s\nMessage:\n %s\n %s\n", entry.ID, entry.CreatedAt.Format("2006-01-02"), entry.Message, divider)
	}
	return []byte(output)
}

func OutputEntriesToMarkdown(entries []Entry) error {
	file, err := os.OpenFile("./output.md", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return utils.ErrCannotCreateFile
	}
	defer file.Close() // want defer as close to acquisition of resources as possible
	output := formattedOutputFromEntries(entries)
	_, err = file.Write(output)
	if err != nil {
		return utils.ErrCannotSaveFile
	}
	return nil
}

func OutputEntriesToPDF(entries []Entry) error {
	output := formattedOutputFromEntries(entries)              // []byte
	pandoc := exec.Command("pandoc", "-s", "-o", "output.pdf") // c is going to run pandoc, so I'm assigning the pipe to c
	wc, wcerr := pandoc.StdinPipe()                            // io.WriteCloser, err
	if wcerr != nil {
		return utils.ErrPandoc
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
		return errors.Wrap(err, utils.CannotWriteToFilePandoc)
	}
	err := pandoc.Run()
	if err != nil {
		return errors.Wrap(err, utils.CannotRunPandoc)
	}
	return nil
}
