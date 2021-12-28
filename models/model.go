package models

type model interface {
	GetID() uint
}

type modelImpl struct {
	ID uint
}

func (m *modelImpl) GetID() uint {
	return m.ID
}
