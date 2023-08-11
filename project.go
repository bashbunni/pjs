package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

// TODO: add a note about manually deleting projects

// Project is a project you'd like to track notes for
type Project string

func NewProject(name string) (Project, error) {
	p := Project(name)
	err := write(p.Path())
	return p, err
}

// Path: returns the project path
func (p Project) Path() string {
	pwd, _ := defaultHome()
	return filepath.Join(pwd, string(p))
}

func write(path string) error {
	if err := os.Mkdir(path, 0o755); err != nil {
		return fmt.Errorf("unable to create new project: %w", err)
	}

	if err := os.Mkdir(fmt.Sprintf("%s/.archived", path), 0o755); err != nil {
		return fmt.Errorf("unable to create archived dir: %w", err)
	}

	return nil
}

// Delete: archives the project directory
func (p Project) Delete() error {
	path, _ := defaultHome()
	return os.Rename(p.Path(), fmt.Sprintf("%s/.archived/%s", path, string(p)))
}

// Rename: rename project
func (p Project) Rename(name string) error {
	path, _ := defaultHome()
	return os.Rename(p.Path(), fmt.Sprintf("%s/%s", path, name))
}

// NewProjectPrompt create a new project from user input to console
func NewProjectPrompt() string {
	var name string
	fmt.Println("what would you like to name your project?")
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Scan()
	name = scanner.Text()
	return name
}

/* implementing list.Item */
func (p Project) Title() string       { return string(p) }
func (p Project) Description() string { return "" }
func (p Project) FilterValue() string { return string(p) }
