package main

import (
	"fmt"
	"io/fs"
	"log"
	"os"

	"github.com/bashbunni/pjs/project"
)

// TODO: Defaults to $HOME/.pjs, can be changed by an env variable.
// TODO: have subdirectories named by project
// TODO: files named by date

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
func getProjects(home string) (projects []string, err error) {
	var de []fs.DirEntry
	de, err = os.ReadDir(home)
	if err != nil {
		return projects, err
	}
	for _, name := range de {
		projects = append(projects, name.Name())
	}
	return projects, err
}

func createProject(name string) error {
	if err := os.Mkdir(name, 0o755); err != nil {
		return fmt.Errorf("unable to create new project: %w", err)
	}
	return nil
}

func main() {
	var projects []string
	var home string

	// init home
	home, err := defaultHome()
	if err != nil {
		log.Fatal(err)
	}
	projects, err = getProjects(home)
	if err != nil {
		log.Fatal(err)
	}
	if len(projects) < 1 {
		name := project.NewProjectPrompt()
		if err := createProject(fmt.Sprintf("%s/%s", home, name)); err != nil {
			log.Fatal(err)
		}
	}
	projects, err = getProjects(home)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(home)
	fmt.Printf("%+v\n", projects)
}
