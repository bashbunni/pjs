module github.com/bashbunni/project-management

go 1.16

replace github.com/bashbunni/project-management/models => ./models

require (
	github.com/charmbracelet/bubbles v0.9.0
	github.com/charmbracelet/bubbletea v0.19.1
	github.com/charmbracelet/lipgloss v0.4.0
	github.com/google/go-cmp v0.5.6
	github.com/pkg/errors v0.9.1
	gorm.io/driver/sqlite v1.1.4
	gorm.io/gorm v1.21.10
)
