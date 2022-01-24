module github.com/bashbunni/project-management

go 1.17

replace github.com/bashbunni/project-management/models => ./models

require (
	github.com/charmbracelet/bubbles v0.10.2
	github.com/charmbracelet/bubbletea v0.19.3
	github.com/charmbracelet/lipgloss v0.4.0
	github.com/containerd/console v1.0.3 // indirect
	github.com/mattn/go-isatty v0.0.14 // indirect
	github.com/mattn/go-sqlite3 v1.14.10 // indirect
	github.com/muesli/ansi v0.0.0-20211031195517-c9f0611b6c70 // indirect
	github.com/pkg/errors v0.9.1
	golang.org/x/sys v0.0.0-20220114195835-da31bd327af9 // indirect
	golang.org/x/term v0.0.0-20210927222741-03fcf44c2211 // indirect
	gorm.io/driver/sqlite v1.2.6
	gorm.io/gorm v1.22.5
)
