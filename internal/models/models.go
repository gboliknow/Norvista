package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

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

type UserResponse struct {
	ID        string         `gorm:"primaryKey"`
	Email     string         `gorm:"type:varchar(255);unique;not null"`
	FirstName string         `gorm:"type:varchar(255);not null"`
	LastName  string         `gorm:"type:varchar(255);not null"`
	Role      string         `gorm:"type:varchar(50);not null"`
	CreatedAt time.Time      `gorm:"index"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Phone     string         `gorm:"type:varchar(20)"`
	Address   string         `gorm:"type:varchar(255)"`
}

type Movie struct {
	ID          string         `gorm:"primaryKey"`
	Title       string         `gorm:"type:varchar(255);not null"`
	Description string         `gorm:"type:text"`
	Genre       string         `gorm:"type:varchar(100);not null"`
	PosterURL   string         `gorm:"type:varchar(255)"`
	ReleaseDate string         `gorm:"type:date"`
	CreatedAt   time.Time      `gorm:"index"`
	DeletedAt   gorm.DeletedAt `gorm:"index"`
	Showtimes   []Showtime     `gorm:"foreignKey:MovieID"` // One-to-many relationship
}

// Showtime model
type Showtime struct {
	ID        string         `gorm:"primaryKey"`
	MovieID   string         `gorm:"type:uuid;not null"` // Foreign key to Movie
	StartTime time.Time      `gorm:"type:timestamp;not null"`
	EndTime   time.Time      `gorm:"type:timestamp;not null"`
	CreatedAt time.Time      `gorm:"index"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
	Movie     Movie          `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"` // Movie association
}

type ShowtimeRequest struct {
	MovieID   string    `json:"movieId" binding:"required"`
	StartTime time.Time `json:"startTime" binding:"required"`
	EndTime   time.Time `json:"endTime" binding:"required"`
}

// type Seat struct {
// 	ID         string `gorm:"primaryKey"`
// 	ShowtimeID string `gorm:"type:uuid;not null"`
// 	Number     string `gorm:"not null"`
// 	IsReserved bool   `gorm:"default:false"`
// 	CreatedAt  time.Time
// 	DeletedAt  gorm.DeletedAt
// }

// type Reservation struct {
// 	ID         string         `gorm:"primaryKey"`
// 	UserID     string         `gorm:"type:uuid;not null"`
// 	ShowtimeID string         `gorm:"type:uuid;not null"`
// 	Seats      []Seat         `gorm:"many2many:reservation_seats;"`
// 	TotalPrice float64        `gorm:"not null"`
// 	CreatedAt  time.Time      `gorm:"index"`
// 	DeletedAt  gorm.DeletedAt `gorm:"index"`
// 	User       User           `gorm:"foreignKey:UserID"`
// 	Showtime   Showtime       `gorm:"foreignKey:ShowtimeID"`
// }

type Response struct {
	StatusCode int         `json:"statusCode"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data,omitempty"` // Data is omitted if nil or empty
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
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

// func (reservation *Reservation) BeforeCreate(tx *gorm.DB) (err error) {
// 	reservation.ID = uuid.New().String()
// 	return
// }
