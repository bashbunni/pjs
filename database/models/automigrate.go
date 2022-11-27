package models

import "github.com/bashbunni/project-management/database/dbconn"

type Model interface{}

var models = []Model{}

func AutoMigrate(db dbconn.GormWrapper) error {
	if models != nil {
		if err := db.AutoMigrate(models); err != nil {
			return err
		}
	}

	return nil
}

func registerForAutomigration(m Model) {
	models = append(models, m)
}
