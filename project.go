package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"

	"github.com/pkg/errors"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

/* constants */
const Markdown = "markdown"
const Csv = "csv"
const Format = "%d : %s\n"
const DefaultInput = 1

/* flags */
var (
	// stringvar := flag.String("optionname", "defaultvalue", "description of the flag")
	cEntry      = flag.Bool("ce", false, "create a new entry")
	deleteEntry = flag.Int("de", -1, "delete an existing entry; default is -1")
	deleteProj  = flag.Int("dp", -1, "delete an existing project; default is -1")
	editProj    = flag.Int("ep", -1, "rename an existing project; default is empty string")
	markdown    = flag.Bool("md", false, "output all entries to markdown file")
)

/* functions */

/* queries */

// mainMenu: flag action handling
func handleFlags(db *gorm.DB) {
	flag.Parse()
	var entries []Entry
	db.Find(&entries) // contains all data from table
	if *cEntry != false {
		createEntry(db)
	}
	if *deleteEntry != -1 {
		DeleteEntry(*deleteEntry, db)
	}
	if *deleteProj != -1 {
		DeleteProject(*deleteProj, db)
	}
	if *markdown != false {
		OutputMarkdown(entries)
	}
	if *editProj != -1 {
		RenameProject(*editProj, db)
	}
}

/* other */

// OpenFileInEditor: a new file in nvim or default editor; helper function
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

// CaptureInputFromEditor: temp file, edit it, delete it
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

// projectPrompt: input validation to create new projects or edit existing
func projectPrompt(db *gorm.DB) Project {
	var input int
	PrintProjects(db)
	fmt.Println("Project ID: ")
	fmt.Scanf("%d", &input)
	// read in input + assign to project
	fmt.Printf("selection is %d \n", input)
	return NewProject(input, db)
}

// createEntry: write and save entry
func createEntry(db *gorm.DB) error {
	message, err := CaptureInputFromEditor()
	if err != nil {
		return errors.Wrap(err, "could not open editor")
	}
	// convert []byte to string can be done vvv
	fmt.Println(string(message[:]))
	myproject := projectPrompt(db)
	myproject.SaveNewEntry(string(message[:]), db)
	return nil
}

func main() {
	// setup
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{
		PrepareStmt: true, // caches queries for faster calls
	})
	if err != nil {
		panic("failed to connect database")
	}
	db.AutoMigrate(&Entry{}, &Project{})
	handleFlags(db)
}

// https://gorm.io/docs/#Quick-Start
