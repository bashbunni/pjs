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

// DeleteEntry: delete an entry by id
func DeleteEntry(pKey int, db *gorm.DB) {
	fmt.Println(pKey)
	db.Delete(&Entry{}, pKey)
}

// GetEntriesByDate: return all entries in a date range
func GetEntriesByDate(start time.Time, end time.Time, db *gorm.DB) []Entry {
	var entries []Entry
	db.Where("created_at >= ? and created_at <= ?", start, end).Find(&entries)
	return entries
}

// CreateEntry: write and save entry
func CreateEntry(pKey int, db *gorm.DB) {
	message := utils.CaptureInputFromFile()
	// convert []byte to string can be done vvv
	myproject := GetOrCreateProject(pKey, db)
	db.Create(&Entry{Message: string(message[:]), ProjectId: myproject.ID})
	fmt.Println(string(message[:]) + " was successfully written to " + myproject.Name)
}

func formattedOutputFromEntries(entries []Entry) string {
	var output string
	for _, entry := range entries {
		output += fmt.Sprintf("\n# %s\n#### %s\n\n%s\n", entry.CreatedAt.Format("Monday, Jan 2, 2006"), entry.CreatedAt.Format(time.Kitchen), entry.Message)
	}
	return output
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
	_, err = file.WriteString(output)

	if err != nil {
		return err
	}
	return nil
}

// OutputPDF writes an array of Entry to a PDF by piping their pretty-printed versions to pandoc
// TODO take ProjectWithEntries and title the PDF with the project name
// TODO group multiple entries on the same date under the same date header
// TODO better errors when output markdown is invalid LaTeX
func OutputPDF(entries []Entry) error {
	pandoc := exec.Command("pandoc", "-f", "markdown", "-o", "output.pdf")
	pandocStdin, pipeErr := pandoc.StdinPipe()
	if pipeErr != nil { return pipeErr }
	go func() {
		defer pandocStdin.Close()
		_, _ = pandocStdin.Write([]byte(formattedOutputFromEntries(entries)))
	}()
	_, runPandocErr := pandoc.CombinedOutput()
	if runPandocErr != nil { return runPandocErr }
	return nil
}
