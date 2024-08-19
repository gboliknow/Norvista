package database

import (
	"Norvista/internal/models"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresStorage struct {
	db *gorm.DB
}

func NewPostgresStorage(connStr string) (*PostgresStorage, error) {
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		log.Fatal("Unable to connect to database:", err)
	}

	fmt.Println("Connected to PostgreSQL!")
	return &PostgresStorage{db: db}, nil
}

func (s *PostgresStorage) InitializeDatabase() (*gorm.DB, error) {
	// Migrate the schema
	if err := s.db.AutoMigrate(
		&models.User{},
		&models.Movie{},
		&models.Showtime{},
		&models.Reservation{}); err != nil {
		return nil, err
	}
	return s.db, nil
}
