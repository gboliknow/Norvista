package api

import (
	"Norvista/internal/config"
	"Norvista/internal/models"
	"Norvista/internal/utility"
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"
	"unicode"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

func validateJWT(tokenString string) (*jwt.Token, error) {
	secret := os.Getenv("JWT_SECRET")
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte(secret), nil
	})
}

func CreateJWT(secret []byte, userID string) (string, error) {
	// Create a new JWT token with userID and expiration claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":    userID,
		"expiresAt": time.Now().Add(time.Hour * 24 * 1).Unix(),
	})

	// Sign the token with the provided secret
	tokenString, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func HashPassword(password string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hash), nil
}

func validatePassword(password string) error {
	if len(password) == 0 {
		return errPasswordRequired
	}

	if len(password) < 8 {
		return fmt.Errorf("password must be at least 8 characters long")
	}

	var hasUpper bool
	var hasLower bool
	var hasNumber bool
	var hasSpecial bool

	for _, char := range password {
		switch {
		case unicode.IsUpper(char):
			hasUpper = true
		case unicode.IsLower(char):
			hasLower = true
		case unicode.IsNumber(char):
			hasNumber = true
		case unicode.IsPunct(char) || unicode.IsSymbol(char):
			hasSpecial = true
		}
	}

	if !hasUpper || !hasLower || !hasNumber || !hasSpecial {
		return errPasswordStrength
	}

	return nil
}

func validateUserPayload(user *models.User) error {
	if user.Email == "" {
		return errEmailRequired
	}

	if user.FirstName == "" {
		return errFirstNameRequired
	}

	if user.LastName == "" {
		return errLastNameRequired
	}

	if err := validatePassword(user.Password); err != nil {
		return err
	}
	validRoles := map[string]bool{
		"user":  true,
		"admin": true,
		"guest": true,
	}

	if _, ok := validRoles[user.Role]; !ok {
		return fmt.Errorf("invalid role: %s", user.Role)
	}

	return nil
}

func validateMoviePayload(movie *models.Movie) error {
	if movie.Title == "" {
		return errTitleRequired
	}
	if movie.Description == "" {
		return errDescriptionRequired
	}

	if movie.Genre == "" {
		return errGenreRequired
	}

	if movie.ReleaseDate == "" {

		return errReleaseDateRequired
	}
	validGenres := map[string]bool{
		"Action":          true,
		"Comedy":          true,
		"Drama":           true,
		"Horror":          true,
		"Romance":         true,
		"Sci-Fi":          true,
		"Documentary":     true,
		"Adventure":       true,
		"Thriller":        true,
		"Mystery":        true,
		"Science Fiction": true,
	}

	if _, ok := validGenres[movie.Genre]; !ok {
		return fmt.Errorf("invalid genre: %s", movie.Genre)
	}
	return nil
}
func createAndSetAuthCookie(userID string, w http.ResponseWriter) (string, error) {
	secret := []byte(config.Envs.JWTSecret)
	token, err := CreateJWT(secret, userID)
	if err != nil {
		return "", err
	}

	http.SetCookie(w, &http.Cookie{
		Name:  "Authorization",
		Value: token,
	})

	return token, nil
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenString, err := utility.GetTokenFromRequest(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid token"})
			c.Abort()
			return
		}

		token, err := validateJWT(tokenString)
		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "permission denied"})
			c.Abort()
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token claims"})
			c.Abort()
			return
		}

		userID, ok := claims["userID"].(string)
		if !ok || userID == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "userID not found in token"})
			c.Abort()
			return
		}
		c.Set("userID", userID)
		c.Next()
	}
}

func RequireAdminMiddleware(store Store) gin.HandlerFunc {
	return func(c *gin.Context) {
		requesterID, exists := c.Get("userID")
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "permission denied"})
			c.Abort()
			return
		}
		requester, err := store.FindUserByID(requesterID.(string))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			c.Abort()
			return
		}

		if requester.Role != "admin" {
			c.JSON(http.StatusForbidden, gin.H{"error": "Access Denied: Only admins are authorized to perform this action."})
			c.Abort()
			return
		}

		c.Next()
	}
}

var (
	errEmailRequired       = errors.New("email is required")
	errFirstNameRequired   = errors.New("first name is required")
	errLastNameRequired    = errors.New("last name is required")
	errPasswordRequired    = errors.New("password is required")
	errPasswordStrength    = errors.New("password must be at least 8 characters long and include at least one uppercase letter, one lowercase letter, one number, and one special character")
	errTitleRequired       = errors.New("title is required")
	errDescriptionRequired = errors.New("description is required")
	errGenreRequired       = errors.New("genre is required")
	errReleaseDateRequired = errors.New("release date is required and must be in the format YYYY-MM-DD")
)
