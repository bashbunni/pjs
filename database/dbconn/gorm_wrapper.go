package dbconn

import (
	"errors"

	"gorm.io/gorm"
)

type GormWrapper interface {
	Error() error
	AutoMigrate(...interface{}) error
	Create(interface{}) GormWrapper
	Delete(interface{}, ...interface{}) GormWrapper
	Where(interface{}, ...interface{}) GormWrapper
	First(interface{}, ...interface{}) GormWrapper
	Find(interface{}, ...interface{}) GormWrapper
	Unscoped() GormWrapper
	Save(interface{}) GormWrapper
}

type wrapper struct {
	db *gorm.DB
}

func Wrap(db *gorm.DB) GormWrapper {
	return &wrapper{
		db: db,
	}
}

func (w *wrapper) Error() error {
	return w.db.Error
}

func (w *wrapper) AutoMigrate(m ...interface{}) error {
	if w.db == nil {
		return errors.New("unable to migrate tables for nil connection")
	}
	return w.db.AutoMigrate(m)
}

func (w *wrapper) Create(value interface{}) GormWrapper {
	w.db.Create(value)
	return w
}

func (w *wrapper) Delete(value interface{}, args ...interface{}) GormWrapper {
	w.db.Delete(value, args...)
	return w
}

func (w *wrapper) Where(query interface{}, args ...interface{}) GormWrapper {
	w.db.Where(query, args...)
	return w
}

func (w *wrapper) First(dest interface{}, conds ...interface{}) GormWrapper {
	w.db = w.db.First(dest, conds)
	return w
}

func (w *wrapper) Find(dest interface{}, conds ...interface{}) GormWrapper {
	w.db.Find(dest, conds)
	return w
}

func (w *wrapper) Unscoped() GormWrapper {
	w.db = w.db.Unscoped()
	return w
}

func (w *wrapper) Save(value interface{}) GormWrapper {
	w.db.Save(value)
	return w
}
