package utils

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt"
)

// GenerateJWTToken generates a JWT token for organization members
func GenerateJWTToken(userID, email, organizationID, role string) (string, error) {
	// Create claims with multiple fields
	claims := jwt.MapClaims{
		"user_id":         userID,
		"email":           email,
		"organization_id": organizationID,
		"role":            role,
		"exp":             time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	}

	// Create token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Get JWT secret from environment variable
	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		jwtSecret = "your-secret-key" // Fallback secret, should be replaced with a proper secret
	}

	// Generate encoded token
	tokenString, err := token.SignedString([]byte(jwtSecret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}
