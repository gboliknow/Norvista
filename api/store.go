package api

import (
	"Norvista/internal/models"
	"Norvista/internal/utility"
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
	CreateShowtime(showtime *models.Showtime) error
	GetShowtimeByID(showtimeID string) (*models.Showtime, error)
	DeleteShowtime(showtimeID string) error
	GetAllShowtimes(showtimes *[]models.Showtime) error
	UpdateShowtime(showtime *models.Showtime) error

	//reservation
	CreateSeats(seats []models.Seat) error
	GetSeatsByShowtimeID(showtimeID string) ([]models.Seat, error)
	CreateReservation(reservation *models.Reservation) error
	GetSeatByID(seatID string) (*models.Seat, error)
	ReserveSeat(seatID string) error
	GetReservationsByUser(userID string) ([]models.Reservation, error)
	GetReservationsByShowtime(showtimeID string) ([]models.Reservation, error)
	GetSeatBySeatNumber(seatNumber string) (*models.Seat, error)
	CancelReservation(reservationID string) error
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

func (s *Storage) UpdateMovie(movie *models.Movie) error {
	return s.db.Omit("Showtimes").Save(movie).Error
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

func (s *Storage) CreateShowtime(showtime *models.Showtime) error {
	showtime.ID = uuid.New().String()
	return s.db.Create(showtime).Error
}

func (s *Storage) GetShowtimeByID(showtimeID string) (*models.Showtime, error) {
	var showtime models.Showtime
	err := s.db.Where("id = ?", showtimeID).First(&showtime).Error
	if err != nil {
		return nil, err
	}
	return &showtime, nil
}

func (s *Storage) DeleteShowtime(showtimeID string) error {
	var showtime models.Showtime
	if err := s.db.Where("id = ?", showtimeID).First(&showtime).Error; err != nil {
		return err
	}
	if err := s.db.Delete(&showtime).Error; err != nil {
		return err
	}
	return nil
}

func (s *Storage) GetAllShowtimes(showtimes *[]models.Showtime) error {
	return s.db.Preload("Movie").Find(showtimes).Error
}

func (s *Storage) UpdateShowtime(showtime *models.Showtime) error {
	return s.db.Save(showtime).Error
}

// reservation
func (s *Storage) CreateSeats(seats []models.Seat) error {
	return s.db.Create(&seats).Error
}

func (s *Storage) GetSeatsByShowtimeID(showtimeID string) ([]models.Seat, error) {
	var seats []models.Seat
	if err := s.db.Where("showtime_id = ?", showtimeID).Find(&seats).Error; err != nil {
		return nil, err
	}
	return seats, nil
}

func (s *Storage) GetSeatByID(seatID string) (*models.Seat, error) {
	var seat models.Seat
	if err := s.db.Where("id = ?", seatID).First(&seat).Error; err != nil {
		return nil, err
	}
	return &seat, nil
}

func (s *Storage) GetSeatBySeatNumber(seatNumber string) (*models.Seat, error) {
	var seat models.Seat
	if err := s.db.Where("seat_number = ?", seatNumber).First(&seat).Error; err != nil {
		return nil, err
	}
	return &seat, nil
}

func (s *Storage) GetReservationsByShowtime(showtimeID string) ([]models.Reservation, error) {
	var reservations []models.Reservation
	err := s.db.Where("showtime_id = ?", showtimeID).Find(&reservations).Error
	return reservations, err
}

func (s *Storage) GetReservationsByUser(userID string) ([]models.Reservation, error) {
	var reservations []models.Reservation
	//confirmed it avoids the "N+1 query problem,"
	err := s.db.Preload("User").Preload("Seat").Preload("Showtime").Where("user_id = ?", userID).Find(&reservations).Error
	return reservations, err
}

func (s *Storage) ReserveSeat(seatID string) error {
	return s.db.Model(&models.Seat{}).Where("id = ?", seatID).Update("is_reserved", true).Error
}

func (s *Storage) CreateReservation(reservation *models.Reservation) error {
	return s.db.Create(reservation).Error
}

func (s *Storage) CancelReservation(reservationID string) error {
	tx := s.db.Begin()
	var reservation models.Reservation
	if err := tx.Preload("Showtime").Where("id = ?", reservationID).First(&reservation).Error; err != nil {
		tx.Rollback()
		return err
	}

	if time.Until(reservation.Showtime.StartTime) <= 24*time.Hour {
		tx.Rollback()
		return utility.ErrCancellationTooSoon
	}

	if err := tx.Delete(&reservation).Error; err != nil {

		tx.Rollback()
		return utility.ErrFailedToDelete
	}

	if err := tx.Model(&models.Seat{}).Where("id = ?", reservation.SeatID).Update("is_reserved", false).Error; err != nil {
		tx.Rollback()
		return utility.ErrFailedToUpdateSeat
	}

	if err := tx.Commit().Error; err != nil {
		return utility.ErrFailedToCommit
	}

	return nil

}
