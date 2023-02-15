package utils

import (
	"fmt"
	"os"
)

// CreateTempFile creates a temporary file to be opened in the editor
func CreateTempFile() (*os.File, error) {
	file, err := os.CreateTemp(os.TempDir(), "*")
	if err != nil {
		return file, fmt.Errorf("unable to create new file: %w", err)
	}
	return file, nil
}

// ReadFile returns the contents of the temp file as a string of bytes
func ReadFile(file *os.File) ([]byte, error) {
	bytes, err := os.ReadFile(file.Name())
	if err != nil {
		return []byte(""), fmt.Errorf("%w: unable to read temp file: %s", err, file.Name())
	}
	return bytes, nil
}
