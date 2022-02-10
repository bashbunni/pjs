package frontend

import "github.com/bashbunni/project-management/models"

type updateEntryListMsg struct {
	entries []models.Entry
}
type errMsg struct {error}
// TODO: have this implement Error()

type createProjectListMsg struct {
	project models.Project
}
