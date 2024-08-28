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
	FindUserByID(userID string) (*models.User, error)
	UpdateUser(user *models.User) error
	GetAllUsers() ([]models.User, error)
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
		user.Role = "user"
	}
	user.ID = uuid.New().String()
	user.CreatedAt = time.Now()
	if err := s.db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}

func (db *Storage) FindUserByEmail(email string, user *models.User) error {
	return db.db.Where("email = ?", email).First(user).Error
}
func (s *Storage) FindUserByID(userID string) (*models.User, error) {
	var user models.User
	if err := s.db.Where("id = ?", userID).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *Storage) UpdateUser(user *models.User) error {
	return s.db.Save(user).Error
}

func (s *Storage) GetAllUsers() ([]models.User, error) {
    var users []models.User
    if err := s.db.Find(&users).Error; err != nil {
        return nil, err
    }
    return users, nil
}
