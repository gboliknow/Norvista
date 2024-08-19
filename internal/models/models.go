package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// User represents a user model with UUID as the primary key.
type User struct {
	ID        string         `gorm:"primaryKey"`
	Email     string         `gorm:"type:varchar(255);unique;not null"`
	FirstName string         `gorm:"type:varchar(255);not null"`
	LastName  string         `gorm:"type:varchar(255);not null"`
	Password  string         `gorm:"type:varchar(255);not null"`
	Role      string         `gorm:"type:varchar(50);not null"`
	CreatedAt time.Time      `gorm:"index"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Phone     string         `gorm:"type:varchar(20)"`
	Address   string         `gorm:"type:varchar(255)"`
}

// Movie represents a movie model with UUID as the primary key.
type Movie struct {
	ID          string         `gorm:"primaryKey"`
	Title       string         `gorm:"type:varchar(255);not null"`
	Description string         `gorm:"type:text"`
	ReleaseDate string         `gorm:"type:date"`
	CreatedAt   time.Time      `gorm:"index"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
}

// Showtime represents a showtime model with UUID as the primary key.
type Showtime struct {
	ID        string         `gorm:"primaryKey"`
	MovieID   string         `gorm:"type:uuid;not null"`
	StartTime string         `gorm:"type:timestamp;not null"`
	EndTime   string         `gorm:"type:timestamp;not null"`
	CreatedAt time.Time      `gorm:"index"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Movie     Movie          `gorm:"foreignKey:MovieID"`
}

// Reservation represents a reservation model with UUID as the primary key.
type Reservation struct {
	ID         string         `gorm:"primaryKey"`
	UserID     string         `gorm:"type:uuid;not null"`
	ShowtimeID string         `gorm:"type:uuid;not null"`
	Seats      int            `gorm:"type:int;not null"`
	CreatedAt  time.Time      `gorm:"index"`
	DeletedAt  gorm.DeletedAt `gorm:"index"`
	User       User           `gorm:"foreignKey:UserID"`
	Showtime   Showtime       `gorm:"foreignKey:ShowtimeID"`
}

// Initialize the models in the same file
func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	user.ID = uuid.New().String()
	return
}

func (movie *Movie) BeforeCreate(tx *gorm.DB) (err error) {
	movie.ID = uuid.New().String()
	return
}

func (showtime *Showtime) BeforeCreate(tx *gorm.DB) (err error) {
	showtime.ID = uuid.New().String()
	return
}

func (reservation *Reservation) BeforeCreate(tx *gorm.DB) (err error) {
	reservation.ID = uuid.New().String()
	return
}
