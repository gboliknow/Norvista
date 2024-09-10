package database

import (
	"Norvista/api"
	"Norvista/internal/models"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PostgresStorage struct {
	db *gorm.DB
}

func NewPostgresStorage(connStr string) (*PostgresStorage, error) {
	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{})
	if err != nil {
		log.Error().Err(err).Msg("Unable to connect to database")
	}

	log.Info().Msg("Connected to PostgreSQL!")
	return &PostgresStorage{db: db}, nil
}

func (s *PostgresStorage) InitializeDatabase() (*gorm.DB, error) {
	// Migrate the schema
	if err := s.db.AutoMigrate(
		&models.User{},
		&models.Showtime{},
		&models.Movie{},
		&models.Showtime{},
		&models.Seat{},
		&models.Reservation{},
	); err != nil {
		log.Error().Err(err).Msg("Failed to migrate database schema")
		return nil, err
	}
	log.Info().Msg("Database schema migrated successfully")

	if err := s.SeedAdminUser(); err != nil {
		log.Error().Err(err).Msg("Failed to seed initial admin user")
		return nil, err
	}
	return s.db, nil
}

func (s *PostgresStorage) SeedAdminUser() error {
	// Check if an admin already exists
	var adminCount int64
	err := s.db.Model(&models.User{}).Where("role = ?", "admin").Count(&adminCount).Error
	if err != nil {
		return fmt.Errorf("error checking for existing admin: %w", err)
	}
	adminPass := os.Getenv("Default_Admin_Password")
	adminEmail := os.Getenv("Admin_Email")
	if adminCount == 0 {
		hashedPassword, err := api.HashPassword(adminPass)
		if err != nil {
			return fmt.Errorf("error hashing admin password: %w", err)
		}

		initialAdmin := &models.User{
			ID:        uuid.New().String(),
			Email:     adminEmail,
			FirstName: "Admin",
			LastName:  "User",
			Password:  hashedPassword,
			Role:      "admin",
			CreatedAt: time.Now(),
		}
		if err := s.db.Create(initialAdmin).Error; err != nil {
			return fmt.Errorf("error creating initial admin user: %w", err)
		}
		log.Info().Msg("Initial admin user created successfully")
	} else {
		log.Info().Msg("Admin user already exists")
	}

	return nil
}
