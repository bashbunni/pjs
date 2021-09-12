package models

import (
	"fmt"

	"gorm.io/gorm"
)

const notFound uint = 0

type Project struct {
	gorm.Model
	Name string
}

func PrintProjects(db *gorm.DB) {
	if hasProjects(db) {
		projects := GetAllProjects(db)
		for _, project := range projects {
			fmt.Printf(Format, project.ID, project.Name)
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
func getProject(projectId int, db *gorm.DB) Project {
	var project Project
	db.Where("id = ?", projectId).Find(&project)
	return project
}

// getAllProjects: return all projects
func GetAllProjects(db *gorm.DB) []Project {
	var projects []Project
	if hasProjects(db) {
		db.Find(&projects)
	}
	return projects
}

// DeleteProject: delete a project by id
func DeleteProject(projectID int, db *gorm.DB) {
	// what if projectID does not exist?
	db.Where("project_id = ?", projectID).Delete(&Entry{})
	db.Delete(&Project{}, projectID)
}

func GetOrCreateProject(projectID int, db *gorm.DB) Project {
	proj := getProject(projectID, db)
	if proj.ID == notFound {
		return CreateProject("", db)
	}
	return proj
}

func RenameProject(projectID int, db *gorm.DB) {
	name := newProjectPrompt()
	var project Project
	db.Where("id = ?", projectID).First(&project)
	project.Name = name
	db.Save(&project)
}

/* PROMPTS */
func newProjectPrompt() string {
	var name string
	fmt.Println("what would you like to name your project?")
	fmt.Scanf("%s", &name)
	return name
}

func CreateProject(name string, db *gorm.DB) Project {
	if name == "" {
		name = newProjectPrompt()
	}
	proj := Project{Name: name}
	db.Create(&proj)
	return proj
}
