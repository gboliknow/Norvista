package api

import (
	"Norvista/internal/models"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Store interface {
	CreateUser(user *models.User) (*models.User, error)
	FindUserByEmail(email string, user *models.User) error
}

type Storage struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) *Storage {
	return &Storage{
		db: db,
	}
}

func (s *Storage) CreateUser(user *models.User) (*models.User, error) {

	if user.Role == "" {
		user.Role = "user" // Default role
	}
	// Set default values before saving to the database
	user.ID = uuid.New().String()
	user.CreatedAt = time.Now()

	// Use GORM to create the user in the database
	if err := s.db.Create(user).Error; err != nil {
		return nil, err
	}

	// Return the user and nil error if successful
	return user, nil
}

func (db *Storage) FindUserByEmail(email string, user *models.User) error {
	return db.db.Where("email = ?", email).First(user).Error
}