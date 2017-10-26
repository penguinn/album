package db

import "github.com/jinzhu/gorm"

type Model interface {
	ConnectionName() string
}

func ReadModel(m Model) (*gorm.DB, error) {
	r, e := Read(m.ConnectionName())
	if e != nil {
		return nil, e
	}
	return r.Model(m), nil
}

func WriteModel(m Model) (*gorm.DB, error) {
	r, e := Write(m.ConnectionName())
	if e != nil {
		return nil, e
	}
	return r.Model(m), nil
}

func MustReadModel(m Model) *gorm.DB {

	return MustRead(m.ConnectionName()).Model(m)

}

func MustWriteModel(m Model) *gorm.DB {
	return MustWrite(m.ConnectionName()).Model(m)
}
