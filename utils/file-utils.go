package utils

import (
	"fmt"
	"log"
	"os"

	"github.com/pkg/errors"
)

// CreateTempFile creates a temporary file to be opened in the editor
func CreateTempFile() *os.File {
	file, err := os.CreateTemp(os.TempDir(), "*")
	if err != nil {
		log.Fatalf("Unable to create new file: %v\n", err)
	}
	return file
}

// ReadFile returns the contents of the temp file as a string of bytes
func ReadFile(file *os.File) ([]byte, error) {
	bytes, err := os.ReadFile(file.Name())
	if err != nil {
		return []byte(""), errors.Wrap(err, fmt.Sprintf("Unable to read temp file: %s\n", file.Name()))
	}
	return bytes, nil
}
