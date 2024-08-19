package api

import (
	"gorm.io/gorm"
)

type Store interface {
	// Define your methods here
}

type Storage struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) *Storage {
	return &Storage{
		db: db,
	}
}
