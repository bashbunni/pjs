module github.com/bashbunni/project-management

go 1.16

replace github.com/bashbunni/project-management/models => ./models

require (
	github.com/charmbracelet/bubbles v0.9.0
	github.com/charmbracelet/bubbletea v0.19.1
	github.com/charmbracelet/lipgloss v0.4.0
	github.com/mattn/go-sqlite3 v1.14.10 // indirect
	github.com/pkg/errors v0.9.1
	gorm.io/driver/sqlite v1.2.6
	gorm.io/gorm v1.22.5
)
