package models

import (
	"fmt"
	"os"
	"os/exec"
	"time"

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

// DeleteEntry deletes an entry by id
func DeleteEntry(pKey int, db *gorm.DB) {
	fmt.Println(pKey)
	db.Delete(&Entry{}, pKey)
}

// GetEntriesByDate returns all entries in a date range
func GetEntriesByDate(start time.Time, end time.Time, db *gorm.DB) []Entry {
	var entries []Entry
	db.Where("created_at >= ? and created_at <= ?", start, end).Find(&entries)
	return entries
}

// CreateEntry interactively captures entry message from a temp file and saves it in the project identified by pKey
func CreateEntry(pKey int, db *gorm.DB) {
	message := utils.CaptureInputFromFile()
	// convert []byte to string can be done vvv
	myproject := GetOrCreateProject(pKey, db)
	db.Create(&Entry{Message: string(message[:]), ProjectId: myproject.ID})
	fmt.Println(string(message[:]) + " was successfully written to " + myproject.Name)
}

func formattedOutputFromEntries(entries []Entry) []byte {
	var output string
	for _, entry := range entries {
		output += fmt.Sprintf("ID: %d\nCreated: %s\nMessage:\n %s\n %s\n", entry.ID, entry.CreatedAt.Format("2006-01-02"), entry.Message, divider)
	}
	return []byte(output)
}

func OutputMarkdownByDateRange(start time.Time, end time.Time, db *gorm.DB) {
	entries := GetEntriesByDate(start, end, db)
	OutputMarkdown(entries)
}

func OutputMarkdown(entries []Entry) error {
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

func OutputPDF(entries []Entry) error {
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
