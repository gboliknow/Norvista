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

	CreateMovie(movie *models.Movie) (*models.Movie, error)
	UpdateMovie(movie *models.Movie) error
	GetAllMovies() ([]models.Movie, error)
	GetMovieByID(movieID string) (*models.Movie, error)
	DeleteMovie(movieID string) error
}

type Storage struct {
	db *gorm.DB
}

func NewStore(db *gorm.DB) *Storage {
	return &Storage{
		db: db,
	}
}

// Auth
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

//Movies

func (s *Storage) CreateMovie(movie *models.Movie) (*models.Movie, error) {
	movie.ID = uuid.New().String()
	movie.CreatedAt = time.Now()
	if err := s.db.Create(movie).Error; err != nil {
		return nil, err
	}
	return movie, nil
}

// UpdateMovie updates the details of an existing movie.
func (s *Storage) UpdateMovie(movie *models.Movie) error {
	return s.db.Save(movie).Error
}

func (s *Storage) GetAllMovies() ([]models.Movie, error) {
	var movies []models.Movie
	if err := s.db.Find(&movies).Error; err != nil {
		return nil, err
	}
	return movies, nil
}

func (s *Storage) GetMovieByID(movieID string) (*models.Movie, error) {
	var movie models.Movie
	if err := s.db.Preload("Showtimes").First(&movie, "id = ?", movieID).Error; err != nil {
		return nil, err
	}
	return &movie, nil
}

func (s *Storage) DeleteMovie(movieID string) error {
	var movie models.Movie
	if err := s.db.Where("id = ?", movieID).First(&movie).Error; err != nil {
		return err
	}
	if err := s.db.Delete(&movie).Error; err != nil {
		return err
	}
	return nil
}
