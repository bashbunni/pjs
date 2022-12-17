package models

import (
	"fmt"

	"gorm.io/gorm"
)

func init() {
	registerForAutomigration(&Project{})
}

// Project the project holds entries
type Project struct {
	gorm.Model
	Name string
}

// Implement list.Item for Bubbletea TUI

// Title the project title to display in a list
func (p Project) Title() string { return p.Name }

// Description the project description to display in a list
func (p Project) Description() string { return fmt.Sprintf("%d", p.ID) }

// FilterValue choose what field to use for filtering in a Bubbletea list component
func (p Project) FilterValue() string { return p.Name }
