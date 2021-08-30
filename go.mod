module github.com/bashbunni/project-management

go 1.16

replace github.com/bashbunni/project-management/models => ./models

require (
	github.com/pkg/errors v0.9.1
	gorm.io/driver/sqlite v1.1.4
	gorm.io/gorm v1.21.10
)
