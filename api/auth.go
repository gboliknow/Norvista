package api

import (
	"Norvista/internal/config"
	"Norvista/internal/models"

	"fmt"

	"net/http"

	"time"
	"unicode"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

// func WithJWTAuth(handlerFunc http.HandlerFunc, store Store) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		tokenString, err := utility.GetTokenFromRequest(r)
// 		if err != nil {
// 			errorHandler(w, "missing or invalid token")
// 			return
// 		}

// 		token, err := validateJWT(tokenString)
// 		if err != nil {
// 			log.Printf("Failed to authenticate token: %v", err)
// 			errorHandler(w, "permission denied")
// 			return
// 		}

// 		if !token.Valid {
// 			log.Printf("Token is invalid")
// 			errorHandler(w, "permission denied")
// 			return
// 		}

// 		claims, ok := token.Claims.(jwt.MapClaims)
// 		if !ok {
// 			log.Printf("Invalid token claims")
// 			errorHandler(w, "permission denied")
// 			return
// 		}

// 		userIDStr, ok := claims["userID"].(string)
// 		if !ok {
// 			log.Printf("UserID not found or invalid in token")
// 			errorHandler(w, "permission denied")
// 			return
// 		}

// 		userID, err := strconv.ParseInt(userIDStr, 10, 64)
// 		if err != nil {
// 			log.Printf("Failed to parse userID from token: %v", err)
// 			errorHandler(w, "permission denied")
// 			return
// 		}

// 		// _, err = store.GetUserByID(userID)
// 		// if err != nil {
// 		// 	log.Printf("Failed to get user by ID: %v", err)
// 		// 	errorHandler(w, "permission denied")
// 		// 	return
// 		// }

// 		handlerFunc(w, r)
// 	}
// }

// func errorHandler(w http.ResponseWriter, errorString string) {
// 	utility.WriteJSON(w, http.StatusUnauthorized, errorString, nil)
// }

// func validateJWT(tokenString string) (*jwt.Token, error) {
// 	 secret := os.Getenv("JWT_SECRET")
// 	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
// 		}

// 		return []byte(secret), nil
// 	})
// }

func CreateJWT(secret []byte, userID string) (string, error) {
	// Create a new JWT token with userID and expiration claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":    userID,
		"expiresAt": time.Now().Add(time.Hour * 24 * 120).Unix(), // 120 days expiration
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

// func getUserIDFromToken(tokenString string, secret []byte) (int64, error) {
// 	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
// 		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
// 		}
// 		return secret, nil
// 	})
// 	if err != nil {
// 		return 0, err
// 	}

// 	claims, ok := token.Claims.(jwt.MapClaims)
// 	if !ok || !token.Valid {
// 		return 0, fmt.Errorf("invalid token")
// 	}

// 	userIDStr, ok := claims["userID"].(string)
// 	if !ok {
// 		return 0, fmt.Errorf("userID not found in token")
// 	}

// 	userID, err := strconv.ParseInt(userIDStr, 10, 64)
// 	if err != nil {
// 		return 0, fmt.Errorf("invalid userID format")
// 	}

// 	return userID, nil
// }

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
