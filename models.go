package main

import (
	"fmt"
	"time"

	"gorm.io/gorm"
)

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

/* entries */
// SaveNewEntry: save a new entry to the db
func (p *Project) SaveNewEntry(message string, db *gorm.DB) {
	db.Create(&Entry{Message: message, ProjectId: p.ID})
}

// DeleteEntry: delete an entry by id
func DeleteEntry(pKey int, db *gorm.DB) {
	fmt.Println(pKey)
	db.Delete(&Entry{}, pKey)
}

// GetEntriesByDate: return all entries in a date range
func GetEntriesByDate(start time.Time, end time.Time, db *gorm.DB) []Entry {
	var entries []Entry
	db.Where("created_at >= ? and created_at <= ?", start, end).Find(&entries)
	return entries
}

/* project */
func PrintAll(p Project, db *gorm.DB) {
	var entries []Entry
	db.Where("project_id = ?", p.ID).Find(&entries) // note to self: queries should be snakecase
	for _, e := range entries {
		fmt.Printf(Format, e.ID, e.Message)
	}
}

func PrintProjects(db *gorm.DB) {
	if hasProjects(db) {
		projects := getAllProjects(db)
		for _, p := range projects {
			fmt.Printf(Format, p.ID, p.Name)
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

// countProjects: return the number of projects
func countProjects(db *gorm.DB) int {
	var projects []Project
	db.Find(&projects) // note to self: queries should be snakecase
	return len(projects)
}

// getProject: return a project by id
func getProject(projId int, db *gorm.DB) Project {
	var project Project
	db.Where("id = ?", projId).Find(&project)
	return project
}

// getAllProjects: return all projects
func getAllProjects(db *gorm.DB) []Project {
	var projects []Project
	if hasProjects(db) {
		db.Find(&projects)
	}
	return projects
}

// SaveNewProject: create new project
func SaveNewProject(name string, db *gorm.DB) Project {
	proj := Project{Name: name}
	db.Create(&proj)
	return proj
}

// DeleteProject: delete a project by id
func DeleteProject(pKey int, db *gorm.DB) {
	// what if pKey does not exist?
	db.Where("project_id = ?", pKey).Delete(&Entry{})
	db.Delete(&Project{}, pKey)
}

/* other */

// TODO: handle renaming existing projects and tie in w handleFlags
func RenameProject(pKey int, db *gorm.DB) Project {
	var name string
	fmt.Println("what would you like to name your project?")
	fmt.Scanf("%s", &name)
	PrintProjects(db)
	return SaveNewProject(name, db)
}

func NewProject(pKey int, db *gorm.DB) Project {
	proj := getProject(pKey, db)
	if proj.ID == 0 {
		return RenameProject(pKey, db)
	}
	return proj
}
