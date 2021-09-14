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

func hasProjects(db *gorm.DB) bool {
	var projects []Project
	if err := db.Find(&projects).Error; err != nil {
		return false
	}
	return true
}

func getProjectByID(projectId int, db *gorm.DB) Project {
	var project Project
	db.Where("id = ?", projectId).Find(&project)
	return project
}

func GetAllProjects(db *gorm.DB) []Project {
	var projects []Project
	if hasProjects(db) {
		db.Find(&projects)
	}
	return projects
}

func DeleteProject(pe *ProjectWithEntries, db *gorm.DB) {
	// what if projectID does not exist?
	DeleteEntries(pe, db)
	db.Delete(&Project{}, pe.Project.ID)
}

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

func GetOrCreateProjectByID(projectID int, db *gorm.DB) Project {
	proj := getProjectByID(projectID, db)
	if proj.ID == notFound {
		return CreateProject("", db)
	}
	return proj
}

// TODO: make pe's Project a *Project instead to simplify?
func RenameProject(pe *ProjectWithEntries, db *gorm.DB) {
	name := newProjectPrompt()
	var project Project
	db.Where("id = ?", pe.Project.ID).First(&project)
	project.Name = name
	pe.Project.Name = name
	db.Save(&project)
}
