package domain

import "time"

type Meta struct {
	deleted   bool
	createdAt time.Time
	updatedAt time.Time
	version   uint32
}

func NewMeta() Meta {
	return Meta{
		deleted:   false,
		createdAt: time.Time{},
		updatedAt: time.Time{},
		version:   0,
	}
}

func (m Meta) GetDeleted() bool {
	return m.deleted
}

func (m Meta) GetCreatedAt() time.Time {
	return m.createdAt
}

func (m Meta) GetUpdatedAt() time.Time {
	return m.updatedAt
}

func (m Meta) GetVersion() uint32 {
	return m.version
}

func (m *Meta) HydrateMeta(deleted bool, createdAt, updatedAt time.Time, version uint32) {
	m.deleted = deleted
	m.createdAt = createdAt
	m.updatedAt = updatedAt
	m.version = version
}
