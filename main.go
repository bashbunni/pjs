package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
)

// TODO: Defaults to $HOME/.pjs, can be changed by an env variable.
// TODO: have subdirectories named by project
// TODO: files named by date
// TODO: add flag for opening a specific project without opening list

func checkHome(home string) error {
	var mkDirErr error
	if _, err := os.Stat(home); err != nil {
		mkDirErr = os.Mkdir(home, 0o755)
	}
	return mkDirErr
}

func defaultHome() (home string, err error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		err = fmt.Errorf("No home directory found: %w", err)
		return home, err
	}

	home = fmt.Sprintf("%s/.pjs", homeDir)
	err = checkHome(home)
	return home, err
}

// getProjects: get names of all directories in $HOME/.pjs
func getProjects() (projects []Project, err error) {
	// TODO: handle error
	home, _ := defaultHome()
	var de []fs.DirEntry
	de, err = os.ReadDir(home)
	if err != nil {
		return projects, err
	}

	for _, name := range de {
		projects = append(projects, Project(name.Name()))
	}
	return projects, err
}

func main() {
	var projects []Project
	var home string

	// init home
	home, err := defaultHome()
	if err != nil {
		log.Fatal(err)
	}

	projects, err = getProjects()
	if err != nil {
		log.Fatal(err)
	}

	if len(projects) < 1 {
		name := NewProjectPrompt()
		if err := write(fmt.Sprintf("%s/%s", home, name)); err != nil {
			log.Fatal(err)
		}
	}

	projects, err = getProjects()
	if err != nil {
		log.Fatal(err)
	}
	StartTea()
}

func StartTea() {
	if f, err := tea.LogToFile("debug.log", "debug"); err != nil {
		fmt.Println("couldn't open file for logging")
		os.Exit(1)
	} else {
		defer func() {
			err = f.Close()
			if err != nil {
				log.Fatal(err)
			}
		}()
	}

	m := InitModel()
	p := tea.NewProgram(m, tea.WithAltScreen())
	if err := p.Start(); err != nil {
		log.Fatal(err)
	}
}

// Repository CRUD operations
type Repository interface {
	Delete()
	Rename()
}
