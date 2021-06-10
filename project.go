package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

const Markdown = "markdown"
const Csv = "csv"
const format = "%d : %s\n"

/*
TODO:
- open editor of choice to type message
- create new project
- choose project prompmpmppmpt
*/

type Entry struct {
	gorm.Model
	ProjectId uint
	Project   Project
	Message   string
}

type Project struct {
	gorm.Model
	Name string
}

func (e Entry) getMsg() string {
	return e.Message
}

func (e Entry) getId() uint {
	return e.ID
}

func printAll(p Project, db *gorm.DB) {
	// should take in an array of entries
	var entries []Entry
	db.Where("project_id = ?", p.ID).Find(&entries) // note to self: queries should be snakecase
	for _, e := range entries {
		fmt.Printf(format, e.getId(), e.getMsg())
	}
}

func (p *Project) saveNewEntry(message string, db *gorm.DB) {
	db.Create(&Entry{Message: message, ProjectId: p.ID})
}

func saveNewProject(name string, db *gorm.DB) Project {
	proj := Project{Name: name}
	db.Create(&proj)
	return proj
}

func printProjects(db *gorm.DB) {
	var projects []Project
	db.Find(&projects) // note to self: queries should be snakecase
	fmt.Printf(projects)
	/*
		for _, p := range projects {
			fmt.Printf(format, p.id, p.name)
		}
	*/
}

// TODO: test these functions

func OpenFileInEditor(filename string) (err error) {
	editor := os.Getenv("EDITOR")
	// should always have a default, right?

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

// TODO: figure out what this shit is
func CaptureInputFromEditor() ([]byte, error) {
	file, err := ioutil.TempFile(os.TempDir(), "*")
	if err != nil {
		return []byte{}, err
	}
	filename := file.Name()

	defer os.Remove(filename)

	if err = file.Close(); err != nil {
		return []byte{}, err
	}

	if err = OpenFileInEditor(filename); err != nil {
		return []byte{}, err
	}

	bytes, err := ioutil.ReadFile(filename)
	if err != nil {
		return []byte{}, err
	}

	return bytes, err
}

func main() {
	// setup
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	args := os.Args[:]

	if len(args) <= 1 {
		fmt.Println("Please add a message to commit")
		os.Exit(1)
	}

	/*
	   TODO:
	   - create temp file
	   - read in temp file
	*/

	reader := bufio.NewReader(os.Stdin)
	fmt.Println("What project would you like to choose?")
	printProjects(db)
	chosenone, _ := reader.ReadString('\n')
	// read in input + assign to project

	// migrate the schema
	db.AutoMigrate(&Entry{}, &Project{})

	// other things
	/*
		var project Project
		project = saveNewProject("bread's toaster", db)
		project.saveNewEntry(message, db)
	*/
	var entries []Entry
	db.Find(&entries) // contains all data from table
	db.First(&entries)

	printAll(project, db)
}

// https://gorm.io/docs/#Quick-Start
