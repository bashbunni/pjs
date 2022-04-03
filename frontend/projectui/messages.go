package projectui

import "github.com/charmbracelet/bubbles/list"

type errMsg struct{ error } // TODO: have this implement Error()
type updateProjectListMsg struct{}
type renameProjectMsg []list.Item
