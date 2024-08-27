package api

import (
	"Norvista/internal/models"
	"Norvista/internal/utility"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type UserService struct {
	store Store
}

func NewUserService(s Store) *UserService {
	return &UserService{store: s}
}

func (s *UserService) RegisterRoutes(r *gin.RouterGroup) {
	r.POST("/users/register", s.handleUserRegister)
	r.POST("/users/login", s.handleUserLogin)
	r.GET("/users/me", AuthMiddleware(), s.handleGetUserInfo)
}

func (s *UserService) handleUserRegister(c *gin.Context) {
	var payload models.User
	if err := c.ShouldBindJSON(&payload); err != nil {
		utility.WriteJSON(c.Writer, http.StatusBadRequest, "Invalid request payload", nil)
		return
	}

	if payload.Role == "" {
		payload.Role = "user"
	}
	if err := validateUserPayload(&payload); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := HashPassword(payload.Password)
	if err != nil {
		utility.WriteJSON(c.Writer, http.StatusInternalServerError, "Error creating user", nil)
		return
	}
	payload.Password = hashedPassword

	u, err := s.store.CreateUser(&payload)
	if err != nil {
		// Check if the error is due to a unique constraint (email already exists)
		if err == gorm.ErrDuplicatedKey {
			utility.WriteJSON(c.Writer, http.StatusConflict, "Email already exists", nil)
		} else {
			utility.WriteJSON(c.Writer, http.StatusInternalServerError, "Error creating user", nil)
		}
		return
	}

	token, err := createAndSetAuthCookie(u.ID, c.Writer)
	if err != nil {
		utility.WriteJSON(c.Writer, http.StatusInternalServerError, "Error creating user", nil)
		return
	}

	utility.WriteJSON(c.Writer, http.StatusCreated, "Successful", token)
}

func (s *UserService) handleUserLogin(c *gin.Context) {
	var loginRequest models.LoginRequest
	if err := c.ShouldBindJSON(&loginRequest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	// Find the user by email
	var user models.User
	if err := s.store.FindUserByEmail(loginRequest.Email, &user); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "User not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		}
		return
	}

	if !CheckPasswordHash(loginRequest.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}
	// Generate JWT token
	token, err := createAndSetAuthCookie(user.ID, c.Writer)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	responseData := models.UserResponse{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		Address:   user.Address,
		Phone:     user.Phone,
		Role:      user.Role,
	}

	// Return the user and token
	c.JSON(http.StatusOK, gin.H{
		"user":  responseData,
		"token": token,
	})
}

func (s *UserService) handleGetUserInfo(c *gin.Context) {
	userID, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "permission denied"})
		return
	}

	user, err := s.store.FindUserByID(userID.(string))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	responseData := models.UserResponse{
		ID:        user.ID,
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Email:     user.Email,
		CreatedAt: user.CreatedAt,
		Address:   user.Address,
		Phone:     user.Phone,
		Role:      user.Role,
	}

	// Return the user data
	utility.WriteJSON(c.Writer, http.StatusOK, "Fetched User", responseData)

}
