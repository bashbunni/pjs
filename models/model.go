package models

type model interface {
	GetID() uint
}

type modelImpl struct {
	ID uint
}

func (m *modelImpl) GetID() {
	return m.ID
}
