package main

import (
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
const defaultInput = 1

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
	if hasProjects(db) {
		projects := getAllProjects(db)
		for _, p := range projects {
			fmt.Printf(format, p.ID, p.Name)
		}
	} else {
		fmt.Printf("There are no projects available")
	}
}

// error handling in case no projects are found
func hasProjects(db *gorm.DB) bool {
	var projects []Project
	if err := db.Find(&projects).Error; err != nil {
		return false
	}
	return true
}

func countProjects(db *gorm.DB) int {
	var projects []Project
	db.Find(&projects) // note to self: queries should be snakecase
	return len(projects)
}

// TODO: refactor; do we need project? can we just to result?
func getProject(projId int, db *gorm.DB) Project {
	var project Project
	db.Where("id = ?", projId).Find(&project)
	return project
}

// TODO: check if this works
func getProjectByName(projName string, db *gorm.DB) (Project, error) {
	var project Project
	if err := db.Where("name = ?", projName).Find(&project).Error; err != nil {
		return project, fmt.Errorf("Error: Project %s not found", projName)
	}
	db.Where("name = ?", projName).Find(&project)
	return project, nil
}

func getAllProjects(db *gorm.DB) []Project {
	var projects []Project
	if hasProjects(db) {
		db.Find(&projects)
	}
	return projects
}

// open a new file in nvim or default editor; helper function
func OpenFileInEditor(filename string) (err error) {
	editor := os.Getenv("EDITOR")
	// should always have a default, right?
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

// create temp file, edit it, delete it
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

// INPUT VALIDATION
func projectPrompt(db *gorm.DB) Project {
	var input int
	printProjects(db)
	fmt.Println("Project ID: ")
	fmt.Scanf("%d", &input)
	// read in input + assign to project
	fmt.Printf("selection is %d \n", input)
	proj := getProject(input, db)
	if proj.ID == 0 {
		var name string
		fmt.Println("what would you like to name your new project?")
		fmt.Scanf("%s", &name)
		printProjects(db)
		return saveNewProject(name, db)
	}
	return proj
}

func main() {
	// setup
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&Entry{}, &Project{})
	message, err := CaptureInputFromEditor()
	// convert []byte to string can be done vvv
	fmt.Println(string(message[:]))

	myproject := projectPrompt(db)

	// create new entry with the message string
	myproject.saveNewEntry(string(message[:]), db)
	/*
		retrieve the project to add new entry
		see existing entries for that project
	*/

	// migrate the schema
	var entries []Entry
	db.Find(&entries) // contains all data from table
	db.First(&entries)
}

// https://gorm.io/docs/#Quick-Start
