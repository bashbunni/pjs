package utils

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"

	"github.com/pkg/errors"
)

// CaptureInputFromFile capture user input from within their text editor
func CaptureInputFromFile() []byte {
	file := createFile()
	filename := file.Name()
	defer os.Remove(filename)
	if err := file.Close(); err != nil {
		log.Fatalf("Unable to close temp file: %v\n", err)
	}
	if err := openFileInEditor(filename); err != nil {
		log.Fatalf("Unable to open editor: %v\n", err)
	}
	return readFile(filename)
}

func openFileInEditor(filename string) (err error) {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "nvim"
	}
	exe, err := exec.LookPath(editor)
	if err != nil {
		return errors.Wrap(err, "cannot open editor")
	}
	cmd := exec.Command(exe, filename)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func createFile() *os.File {
	file, err := ioutil.TempFile(os.TempDir(), "*")
	if err != nil {
		log.Fatalf("Unable to create new file: %v\n", err)
	}
	return file
}

func readFile(filename string) []byte {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Unable to read temp file: %v\n", err)
		// TODO: do better error handling
	}
	return bytes
}
