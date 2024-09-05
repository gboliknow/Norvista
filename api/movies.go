package api

import (
	"Norvista/internal/models"
	"Norvista/internal/utility"
	"net/http"

	"github.com/gin-gonic/gin"
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
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Movie not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}
	existingMovie.ID = movieID
	if err := c.ShouldBindJSON(&existingMovie); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	if err := s.store.UpdateMovie(existingMovie); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update movie"})
		return
	}

	c.JSON(http.StatusOK, existingMovie)
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
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload. 'movie_id', 'start_time', and 'end_time' are required."})
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

	// Fetch the created showtime for confirmation
	var createdShowtime models.Showtime
	if err := s.store.GetShowtimeByID(showtime.ID, &createdShowtime); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch showtime"})
		return
	}

	utility.WriteJSON(c.Writer, http.StatusCreated, "Showtime created successfully", createdShowtime)
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
	var showtime models.Showtime

	if err := s.store.GetShowtimeByID(showtimeID, &showtime); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Showtime not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch showtime"})
		}
		return
	}

	utility.WriteJSON(c.Writer, http.StatusCreated, "Showtime fetched successfully", showtime)
}
