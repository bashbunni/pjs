package repos

import (
	"github.com/bashbunni/project-management/database/dbconn"
	"github.com/bashbunni/project-management/database/models"
)

type EntryRepository interface {
	DeleteEntryByID(id uint) error
	DeleteEntriesInProjectOfID(id uint) error
	GetEntriesByProjectID(id uint) ([]models.Entry, error)
	CreateEntryInProjectOfID(message []byte, projectID uint) error
}

type entryRepo struct {
	dbConn dbconn.GormWrapper
}

func NewEntryRepo(db dbconn.GormWrapper) EntryRepository {
	return entryRepo{dbConn: db}
}

func (e entryRepo) DeleteEntryByID(id uint) error {
	return e.dbConn.Delete(models.Entry{}, id).Error()
}

func (e entryRepo) DeleteEntriesInProjectOfID(id uint) error {
	return e.dbConn.Where("project_id = ?", id).Delete(models.Entry{}).Error()
}

func (e entryRepo) GetEntriesByProjectID(id uint) ([]models.Entry, error) {
	var entries []models.Entry
	if err := e.dbConn.Where("project_id = ?", id).Find(&entries).Error(); err != nil {
		return nil, err
	}

	return entries, nil
}

func (e entryRepo) CreateEntryInProjectOfID(message []byte, id uint) error {
	entry := models.Entry{Message: string(message), ProjectID: id}
	return e.dbConn.Create(&entry).Error()
}
