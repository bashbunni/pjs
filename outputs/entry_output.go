package outputs

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/bashbunni/project-management/models"
	"github.com/bashbunni/project-management/utils"
	"github.com/pkg/errors"
)

// TODO: make this not trash
const divider = "_______________________________________"

// helpers
func formattedOutputFromEntries(Entries []models.Entry) []byte {
	var output string
	for _, entry := range Entries {
		output += fmt.Sprintf("ID: %d\nCreated: %s\nMessage:\n %s\n %s\n", entry.ID, entry.CreatedAt.Format("2006-01-02"), entry.Message, divider)
	}
	return []byte(output)
}

// globals
func OutputEntriesToMarkdown(entries []models.Entry) error {
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

func OutputEntriesToPDF(entries []models.Entry) error {
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
func main() {
	fmt.Println("vim-go")
}
