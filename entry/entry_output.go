package entry

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/bashbunni/project-management/database/models"
)

const divider = "---"

// FormattedOutputFromEntries format all entries as a single string in reverse chronological order
func FormattedOutputFromEntries(Entries []models.Entry) []byte {
	var output string
	for i := len(Entries) - 1; i >= 0; i-- {
		output += fmt.Sprintf("ID: %d\nCreated: %s\nMessage:\n\n %s\n %s\n", Entries[i].ID, Entries[i].CreatedAt.Format("2006-01-02"), Entries[i].Message, divider)
	}
	return []byte(output)
}

// FormatEntry return the entry details as a formatted string
func FormatEntry(entry models.Entry) string {
	return fmt.Sprintf("ID: %d\nCreated: %s\nMessage:\n\n %s\n %s\n", entry.ID, entry.CreatedAt.Format("2006-01-02"), entry.Message, divider)
}

// ReverseList reverse the provided list
func ReverseList(list []models.Entry) []models.Entry {
	var output []models.Entry
	for i := len(list) - 1; i >= 0; i-- {
		output = append(output, list[i])
	}
	return output
}

// OutputEntriesToMarkdown create an output file that contains the given entries in a formatted string
func OutputEntriesToMarkdown(entries []models.Entry) error {
	file, err := os.OpenFile("./output.md", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0o644)
	if err != nil {
		return fmt.Errorf("Cannot create file: %v", err)
	}
	defer func() {
		err = file.Close() // want defer as close to acquisition of resources as possible
	}()
	output := FormattedOutputFromEntries(entries)
	_, err = file.Write(output)
	if err != nil {
		return fmt.Errorf("Cannot save file: %v", err)
	}
	return err
}

// OutputEntriesToPDF create a PDF from the given entries in their string format
func OutputEntriesToPDF(entries []models.Entry) error {
	output := FormattedOutputFromEntries(entries)              // []byte
	pandoc := exec.Command("pandoc", "-s", "-o", "output.pdf") // c is going to run pandoc, so I'm assigning the pipe to c
	wc, wcerr := pandoc.StdinPipe()                            // io.WriteCloser, err
	if wcerr != nil {
		return fmt.Errorf("Cannot stdin to pandoc: %v", wcerr)
	}
	goerr := make(chan error)
	done := make(chan bool)
	go func() {
		var err error
		defer func() {
			err = wc.Close()
		}()
		_, err = wc.Write(output)
		goerr <- err
		close(goerr)
		close(done)
	}()
	if err := <-goerr; err != nil {
		return fmt.Errorf("Cannot write file to pandoc: %v", err)
	}
	err := pandoc.Run()
	if err != nil {
		return fmt.Errorf("Cannot run pandoc: %v", err)
	}
	return nil
}
