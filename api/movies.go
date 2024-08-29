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
	adminGroup := r.Group("/movies")
	adminGroup.Use(AuthMiddleware())
	adminGroup.Use(RequireAdminMiddleware(s.store))
	{
		adminGroup.POST("/", s.handleCreateMovie)
		adminGroup.PUT("/:id", s.handleUpdateMovie)
		adminGroup.DELETE("/:id", s.handleDeleteMovie)
	}
	r.GET("/movies", s.handleGetAllMovies)
	r.GET("/movies/:id", s.handleGetMovie)
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
