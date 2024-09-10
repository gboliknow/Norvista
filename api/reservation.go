package api

import (
	"Norvista/internal/models"
	"Norvista/internal/utility"
	"fmt"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ReservationService struct {
	store Store
}

func NewReservationService(s Store) *ReservationService {
	return &ReservationService{store: s}
}

func (s *ReservationService) ReservationRoutes(r *gin.RouterGroup) {

	//for users
	reservationGroup := r.Group("/reservation")
	reservationGroup.Use(AuthMiddleware())
	{
		reservationGroup.GET("/me", s.getUserReservations)
		reservationGroup.POST("/", s.handleReserveSeats)
	}

	//for admin
	adminReservationGroup := r.Group("/reservation/admin")
	adminReservationGroup.Use(AuthMiddleware())
	adminReservationGroup.Use(RequireAdminMiddleware(s.store))
	{
		adminReservationGroup.POST("/:id", s.getReservationsByShowtime)
	}

	//for unautheticated users
	r.GET("/seats/:id", s.handleGetSeatsForShowtime)
}

func (s *ReservationService) handleGetSeatsForShowtime(c *gin.Context) {
	showtimeID := c.Param("id")
	seats, err := s.store.GetSeatsByShowtimeID(showtimeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch seats"})
		return
	}

	seatLites := make([]models.SeatLite, len(seats))
	for i, seat := range seats {
		seatLites[i] = ConvertToSeatLite(seat)
	}
	utility.WriteJSON(c.Writer, http.StatusOK, "Seats Fetched successfully", seatLites)
}

func (s *ReservationService) handleReserveSeats(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	var reserveRequest *models.ReservationRequest
	if err := c.ShouldBindJSON(&reserveRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	var seatIDs []string
	for _, seatNumber := range reserveRequest.SeatNumbers {
		seat, err := s.store.GetSeatBySeatNumber(seatNumber)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("Failed to fetch seat %s", seatNumber)})
			return
		}
		if seat.IsReserved {
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Seat %s is already reserved", seat.SeatNumber)})
			return
		}

		seatIDs = append(seatIDs, seat.ID)
	}

	for _, seatID := range seatIDs {
		if err := s.store.ReserveSeat(seatID); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to reserve seat"})
			return
		}
		reservation := models.Reservation{
			ID:         uuid.New().String(),
			UserID:     userID.(string),
			ShowtimeID: reserveRequest.ShowtimeID,
			SeatID:     seatID,
		}
		if err := s.store.CreateReservation(&reservation); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create reservation"})
			return
		}
	}
	utility.WriteJSON(c.Writer, http.StatusOK, "Seats reserved successfully", nil)
}

func (s *ReservationService) getUserReservations(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}
	reservations, err := s.store.GetReservationsByUser(userID.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve reservations"})
		return
	}
	utility.WriteJSON(c.Writer, http.StatusOK, "Reservations fetched successfully", reservations)
}

func (s *ReservationService) getReservationsByShowtime(c *gin.Context) {
	showtimeID := c.Param("showtimeID")
	reservations, err := s.store.GetReservationsByShowtime(showtimeID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve reservations"})
		return
	}
	utility.WriteJSON(c.Writer, http.StatusOK, "Reservations fetched successfully", reservations)
}

func ConvertToSeatLite(seat models.Seat) models.SeatLite {
	return models.SeatLite{
		ID:         seat.ID,
		ShowtimeID: seat.ShowtimeID,
		SeatNumber: seat.SeatNumber,
		IsReserved: seat.IsReserved,
		CreatedAt:  seat.CreatedAt,
		DeletedAt:  seat.DeletedAt,
	}
}
