package utils

import (
	"io/ioutil"
	"log"
	"os"
	"os/exec"
)

// OpenFileInEditor: a new file in nvim or default editor; helper function
func OpenFileInEditor(filename string) (err error) {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "nvim"
	}
	exe, err := exec.LookPath(editor)
	if err != nil {
		return err
	}
	cmd := exec.Command(exe, filename)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

func CreateFile() *os.File {
	file, err := ioutil.TempFile(os.TempDir(), "*")
	if err != nil {
		log.Fatalf("Unable to create new file: %v\n", err)
	}
	return file
}

func ReadFile(filename string) []byte {
	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Fatalf("Unable to read temp file: %v\n", err)
	}
	return bytes
}

// CaptureInputFromFile: temp file, edit it, delete it
func CaptureInputFromFile() []byte {
	file := CreateFile()
	filename := file.Name()
	defer os.Remove(filename)
	if err := file.Close(); err != nil {
		log.Fatalf("Unable to close temp file: %v\n", err)
	}
	if err := OpenFileInEditor(filename); err != nil {
		log.Fatalf("Unable to open editor: %v\n", err)
	}
	return ReadFile(filename)
}
