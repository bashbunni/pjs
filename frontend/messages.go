package frontend

import "github.com/bashbunni/project-management/models"

type updateEntryListMsg struct {
	entries []models.Entry
}
type errMsg struct {error}
// TODO: have this implement Error()

// project

type updateProjectListMsg struct {}

type renameProjectMsg struct {}
