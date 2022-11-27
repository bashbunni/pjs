package models

import "github.com/bashbunni/project-management/database/dbconn"

type Model interface{}

var models = []Model{}

func AutoMigrate(db dbconn.GormWrapper) error {
	for _, m := range models {
		if err := db.AutoMigrate(m); err != nil {
			return err
		}
	}
	return nil
}

func registerForAutomigration(m Model) {
	models = append(models, m)
}
