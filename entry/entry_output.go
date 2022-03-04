package entry

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/pkg/errors"
)

const divider = "---"

// FormattedOutputFromEntries format all entries as a single string in reverse chronological order
func FormattedOutputFromEntries(Entries []Entry) []byte {
	var output string
	for i := len(Entries) - 1; i >= 0; i-- {
		output += fmt.Sprintf("ID: %d\nCreated: %s\nMessage:\n %s\n %s\n", Entries[i].ID, Entries[i].CreatedAt.Format("2006-01-02"), Entries[i].Message, divider)
	}
	return []byte(output)
}

// OutputEntriesToMarkdown create an output file that contains the given entries in a formatted string
func OutputEntriesToMarkdown(entries []Entry) error {
	file, err := os.OpenFile("./output.md", os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		return errors.Wrap(err, errCannotCreateFile)
	}
	defer file.Close() // want defer as close to acquisition of resources as possible
	output := FormattedOutputFromEntries(entries)
	_, err = file.Write(output)
	if err != nil {
		return errors.Wrap(err, errCannotSaveFile)
	}
	return nil
}

// OutputEntriesToPDF create a PDF from the given entries in their string format
func OutputEntriesToPDF(entries []Entry) error {
	output := FormattedOutputFromEntries(entries)              // []byte
	pandoc := exec.Command("pandoc", "-s", "-o", "output.pdf") // c is going to run pandoc, so I'm assigning the pipe to c
	wc, wcerr := pandoc.StdinPipe()                            // io.WriteCloser, err
	if wcerr != nil {
		return errors.Wrap(wcerr, errPandoc)
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
		return errors.Wrap(err, errCannotWriteToFilePandoc)
	}
	err := pandoc.Run()
	if err != nil {
		return errors.Wrap(err, errCannotRunPandoc)
	}
	return nil
}
