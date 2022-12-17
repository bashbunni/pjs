package models

import "gorm.io/gorm"

func init() {
	registerForAutomigration(&Entry{})
}

// Entry the entry model
type Entry struct {
	gorm.Model
	ProjectID uint `gorm:"foreignKey:Project"`
	Message   string
}
