package main

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	tea "github.com/charmbracelet/bubbletea"
)

type Entry []byte

const defaultEditor = "vi"

func openEditorCmd(path string) tea.Cmd {
	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = defaultEditor
	}
	c := exec.Command(editor, path)
	return tea.ExecProcess(c, func(err error) tea.Msg {
		// TODO: return the file contents to update viewport content
		contents, _ := os.ReadFile(path)
		return editorFinishedMsg{err, contents}
	})
}

// CreateFile creates a markdown file to be opened in the editor
func CreateFile(path string) (*os.File, error) {
	today := time.Now().Format("2006-01-02")
	file, err := os.Create(fmt.Sprintf("%s/%s.md", path, today))
	if err != nil {
		return file, fmt.Errorf("unable to create new file: %w", err)
	}
	return file, nil
}

func entryExists(name string) bool {
	if _, err := os.Stat(name); err == nil {
		return true
	}
	return false
}

// ReadFile returns the contents of the temp file as a string of bytes
func ReadFile(file *os.File) ([]byte, error) {
	bytes, err := os.ReadFile(file.Name())
	if err != nil {
		return []byte(""), fmt.Errorf("%w: unable to read temp file: %s", err, file.Name())
	}
	return bytes, nil
}

func getEntries(path string) []Entry {
	var entries []Entry
	de, err := os.ReadDir(path)
	if err != nil {
		fmt.Errorf("unable to read dir: %w", err)
	}

	for _, entry := range de {
		if !entry.IsDir() {
			entries = append(entries, Entry(entry.Name()))
		}
	}
	return entries
}
