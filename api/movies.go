package api

import "github.com/gin-gonic/gin"

type MovieService struct {
	store Store
}

func NewMovieService(s Store) *MovieService {
	return &MovieService{store: s}
}

func (s *MovieService) MoviesRoutes(r *gin.RouterGroup) {

}
