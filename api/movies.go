package api

import (
	"Norvista/internal/models"
	"Norvista/internal/utility"
	"fmt"

	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"gorm.io/gorm"
)

type MovieService struct {
	store Store
}

func NewMovieService(s Store) *MovieService {
	return &MovieService{store: s}
}

func (s *MovieService) MoviesRoutes(r *gin.RouterGroup) {
	movieGroup := r.Group("/movies")
	showtimeGroup := r.Group("/showtimes")
	movieGroup.Use(AuthMiddleware())
	movieGroup.Use(RequireAdminMiddleware(s.store))
	showtimeGroup.Use(AuthMiddleware())
	showtimeGroup.Use(RequireAdminMiddleware(s.store))
	{
		movieGroup.POST("/", s.handleCreateMovie)
		movieGroup.PUT("/:id", s.handleUpdateMovie)
		movieGroup.DELETE("/:id", s.handleDeleteMovie)
		showtimeGroup.POST("/", s.handleCreateShowtime)
		showtimeGroup.DELETE("/:id", s.handleDeleteShowtime)
		showtimeGroup.PUT("/:id", s.handleUpdateShowtime)
	}
	r.GET("/movies", s.handleGetAllMovies)
	r.GET("/movies/:id", s.handleGetMovie)
	r.GET("/showtimes", s.handleGetAllShowtimes)
	r.GET("/showtimes/:id", s.handleGetShowtime)
}

func (s *MovieService) handleCreateMovie(c *gin.Context) {
	var payload models.Movie
	if err := c.ShouldBindJSON(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if err := validateMoviePayload(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	movie, err := s.store.CreateMovie(&payload)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create movie"})
		return
	}
	utility.WriteJSON(c.Writer, http.StatusCreated, "Movie created successfully", movie)

}

func (s *MovieService) handleUpdateMovie(c *gin.Context) {
	movieID := c.Param("id")

	existingMovie, err := s.store.GetMovieByID(movieID)
	if err != nil {
		status := http.StatusInternalServerError
		message := "Internal server error"
		if err == gorm.ErrRecordNotFound {
			status = http.StatusNotFound
			message = "Movie not found"
		}
		c.JSON(status, gin.H{"error": message})
		return
	}

	if err := c.ShouldBindJSON(&existingMovie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if err := s.store.UpdateMovie(existingMovie); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update movie"})
		return
	}

	response := createMovieUpdateResponse(existingMovie)
	utility.WriteJSON(c.Writer, http.StatusOK, "Movie updated successfully", response)
}

func (s *MovieService) handleDeleteMovie(c *gin.Context) {
	movieID := c.Param("id")

	err := s.store.DeleteMovie(movieID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting movie"})
		}
		return
	}

	utility.WriteJSON(c.Writer, http.StatusCreated, "Movie deleted successfully", nil)
}

func (s *MovieService) handleGetAllMovies(c *gin.Context) {
	movies, err := s.store.GetAllMovies()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving movies"})
		return
	}
	utility.WriteJSON(c.Writer, http.StatusCreated, "Movies Fetched successfully", movies)
}

func (s *MovieService) handleGetMovie(c *gin.Context) {
	movieID := c.Param("id")
	movie, err := s.store.GetMovieByID(movieID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error retrieving movie"})
		}
		return
	}

	utility.WriteJSON(c.Writer, http.StatusCreated, "Movie Fetched successfully", movie)
}

func (s *MovieService) handleCreateShowtime(c *gin.Context) {
	var showtimeRequest models.ShowtimeRequest

	// Bind JSON data to ShowtimeRequest struct
	if err := c.ShouldBindJSON(&showtimeRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload. 'movie_id', 'start_time', and 'end_time' are required." + err.Error()})
		return
	}

	// Validate the showtime request
	if err := validateShowtime(&showtimeRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Ensure the movie exists
	movie, err := s.store.GetMovieByID(showtimeRequest.MovieID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	// Create a new Showtime instance
	showtime := models.Showtime{
		MovieID:   movie.ID, // Use the movie ID directly
		StartTime: showtimeRequest.StartTime,
		EndTime:   showtimeRequest.EndTime,
	}

	// Save the showtime in the database
	if err := s.store.CreateShowtime(&showtime); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create showtime"})
		return
	}

	numSeats := 100 // Set this dynamically or as needed
	var seats []models.Seat
	for i := 1; i <= numSeats; i++ {
		seat := models.Seat{
			ID:         uuid.New().String(),
			ShowtimeID: showtime.ID,
			SeatNumber: fmt.Sprintf("Seat-%d", i),
			IsReserved: false,
		}
		seats = append(seats, seat)
	}

	if err := s.store.CreateSeats(seats); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create seats"})
		return
	}
	createdShowtime, err := s.store.GetShowtimeByID(showtime.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch showtime"})
		return
	}

	response := models.ShowtimeLite{
		ID:        createdShowtime.ID,
		MovieID:   createdShowtime.MovieID,
		StartTime: createdShowtime.StartTime,
		EndTime:   createdShowtime.EndTime,
		CreatedAt: createdShowtime.CreatedAt,
		DeletedAt: &createdShowtime.DeletedAt.Time,
	}

	utility.WriteJSON(c.Writer, http.StatusCreated, "Showtime and seats created successfully", response)
}

func (s *MovieService) handleDeleteShowtime(c *gin.Context) {
	showtimeID := c.Param("id")

	err := s.store.DeleteShowtime(showtimeID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Showtime not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error deleting Showtime"})
		}
		return
	}

	utility.WriteJSON(c.Writer, http.StatusCreated, "Showtime deleted successfully", nil)
}

func (s *MovieService) handleGetAllShowtimes(c *gin.Context) {
	var showtimes []models.Showtime
	if err := s.store.GetAllShowtimes(&showtimes); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch showtimes"})
		return
	}

	utility.WriteJSON(c.Writer, http.StatusCreated, "Showtime fetched successfully", showtimes)
}

func (s *MovieService) handleGetShowtime(c *gin.Context) {
	showtimeID := c.Param("id")
	showtime, err := s.store.GetShowtimeByID(showtimeID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Showtime not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch showtime"})
		}
		return
	}

	utility.WriteJSON(c.Writer, http.StatusCreated, "Showtime fetched successfully", showtime)
}

func (s *MovieService) handleUpdateShowtime(c *gin.Context) {
	// Get the showtime ID from the URL parameters
	showtimeID := c.Param("id")
	existingShowtime, err := s.store.GetShowtimeByID(showtimeID)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Showtime not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	oldMovieID := existingShowtime.MovieID
	if err := c.ShouldBindJSON(&existingShowtime); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}
	if existingShowtime.MovieID != oldMovieID {
		_, err := s.store.GetMovieByID(existingShowtime.MovieID)
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid MovieID"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
	}

	if err := s.store.UpdateShowtime(existingShowtime); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update showtime"})
		return
	}

	response := models.ShowtimeLite{
		ID:        existingShowtime.ID,
		MovieID:   existingShowtime.MovieID,
		StartTime: existingShowtime.StartTime,
		EndTime:   existingShowtime.EndTime,
		CreatedAt: existingShowtime.CreatedAt,
		DeletedAt: &existingShowtime.DeletedAt.Time,
	}
	utility.WriteJSON(c.Writer, http.StatusOK, "Showtime edited successfully", response)
}

func createMovieUpdateResponse(movie *models.Movie) models.MovieUpdateResponse {
	response := models.MovieUpdateResponse{
		ID:          movie.ID,
		Title:       movie.Title,
		Description: movie.Description,
		Genre:       movie.Genre,
		PosterURL:   movie.PosterURL,
		ReleaseDate: movie.ReleaseDate,
		CreatedAt:   movie.CreatedAt,
		DeletedAt:   nil,
	}

	for _, showtime := range movie.Showtimes {
		response.Showtimes = append(response.Showtimes, models.ShowtimeLite{
			ID:        showtime.ID,
			MovieID:   showtime.MovieID,
			StartTime: showtime.StartTime,
			EndTime:   showtime.EndTime,
			CreatedAt: showtime.CreatedAt,
			DeletedAt: &showtime.DeletedAt.Time,
		})
	}
	return response
}
