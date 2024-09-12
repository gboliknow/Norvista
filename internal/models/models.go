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

type Seat struct {
	ID         string         `gorm:"primaryKey"`
	ShowtimeID string         `gorm:"type:uuid;not null"` // Foreign key to Showtime
	SeatNumber string         `gorm:"type:varchar(10);not null"`
	IsReserved bool           `gorm:"default:false"`
	CreatedAt  time.Time      `gorm:"index"`
	DeletedAt  gorm.DeletedAt `gorm:"index"`
	Showtime   Showtime       `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"` // Showtime association
}

type Reservation struct {
	ID         string         `gorm:"primaryKey"`
	UserID     string         `gorm:"type:uuid;not null"` // Foreign key to User
	ShowtimeID string         `gorm:"type:uuid;not null"` // Foreign key to Showtime
	SeatID     string         `gorm:"type:uuid;not null"` // Foreign key to Seat
	CreatedAt  time.Time      `gorm:"index"`
	DeletedAt  gorm.DeletedAt `gorm:"index"`
	User       User           `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"` // User association
	Showtime   Showtime       `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"` // Showtime association
	Seat       Seat           `gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"` // Seat association
}

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

type MovieUpdateResponse struct {
	ID          string         `json:"ID"`
	Title       string         `json:"Title"`
	Description string         `json:"Description"`
	Genre       string         `json:"Genre"`
	PosterURL   string         `json:"PosterURL"`
	ReleaseDate string         `json:"ReleaseDate"`
	CreatedAt   time.Time      `json:"CreatedAt"`
	DeletedAt   *time.Time     `json:"DeletedAt,omitempty"`
	Showtimes   []ShowtimeLite `json:"Showtimes"` // Custom struct without nested Movie
}

type ShowtimeLite struct {
	ID        string     `json:"ID"`
	MovieID   string     `json:"MovieID"`
	StartTime time.Time  `json:"StartTime"`
	EndTime   time.Time  `json:"EndTime"`
	CreatedAt time.Time  `json:"CreatedAt"`
	DeletedAt *time.Time `json:"DeletedAt,omitempty"`
}


type SeatLite struct {
	ID         string         `gorm:"primaryKey"`
	ShowtimeID string         `gorm:"type:uuid;not null"` // Foreign key to Showtime
	SeatNumber string         `gorm:"type:varchar(10);not null"`
	IsReserved bool           `gorm:"default:false"`
	CreatedAt  time.Time      `gorm:"index"`
	DeletedAt  gorm.DeletedAt `gorm:"index"` // Showtime association
}
type ReservationRequest struct {
	ShowtimeID string   `json:"showtimeID" binding:"required"` // Showtime ID
	SeatNumbers []string `json:"seatNumbers" binding:"required"` // Array of seat numbers
}

type ReservationLite struct {
	ID         string         `gorm:"primaryKey"`
	UserID     string         `gorm:"type:uuid;not null"` // Foreign key to User
	ShowtimeID string         `gorm:"type:uuid;not null"` // Foreign key to Showtime
	SeatID     string         `gorm:"type:uuid;not null"` // Foreign key to Seat
	CreatedAt  time.Time      `gorm:"index"`
	DeletedAt  gorm.DeletedAt `gorm:"index"`
	UserName   string    `json:"userName"`     // User's name for display
	SeatNumber string    `json:"seatNumber"`   // Seat number
	Showtime   time.Time `json:"showtime"`
}