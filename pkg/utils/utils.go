package utils

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

// Constants for JWT expiration times and secret key.
const (
	AccessTokenExpiry  = time.Hour * 1
	RefreshTokenExpiry = time.Hour * 72
	SecretKey          = "secret_key"
)

// Claims holds the standard JWT claims plus additional custom fields.
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

// GenerateTokens creates JWT access and refresh tokens for a user.
func GenerateTokens(username, email string) (accessToken string, refreshToken string, err error) {
	// Define the claims of the access and refresh tokens.
	accessClaims := jwt.MapClaims{
		"username": username,
		"email":    email,
		"exp":      time.Now().Add(AccessTokenExpiry).Unix(),
	}
	refreshClaims := jwt.MapClaims{
		"username": username,
		"email":    email,
		"exp":      time.Now().Add(RefreshTokenExpiry).Unix(),
	}

	// Create the access and refresh token objects.
	accessTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	refreshTokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, refreshClaims)

	// Sign the tokens with the secret key.
	accessToken, err = accessTokenObj.SignedString([]byte(SecretKey))
	if err != nil {
		return "", "", err
	}
	refreshToken, err = refreshTokenObj.SignedString([]byte(SecretKey))
	if err != nil {
		return "", "", err
	}

	return accessToken, refreshToken, nil
}

// HashPassword secures a plaintext password using bcrypt.
func HashPassword(password string) (string, error) {
	// Generate a bcrypt hash from the plaintext password.
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

// VerifyRefreshToken checks the validity of a refresh token and returns the username and email.
func VerifyRefreshToken(refreshToken string) (string, string, error) {
	// Parse and validate the refresh token.
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(SecretKey), nil
	})

	// Handle parsing errors.
	if err != nil {
		return "", "", err
	}

	// Extract and return the username and email from the token claims.
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		username, ok1 := claims["username"].(string)
		email, ok2 := claims["email"].(string)
		if !ok1 || !ok2 {
			return "", "", errors.New("invalid claims")
		}
		return username, email, nil
	} else {
		return "", "", errors.New("invalid refresh token")
	}
}

// ValidateToken parses and validates a JWT token string.
func ValidateToken(tokenString string) (*Claims, error) {
	// Parse the token with the custom claims structure.
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	if err != nil {
		return nil, err
	}

	// Return the validated claims if successful.
	if claims, ok := token.Claims.(*Claims); ok && token.Valid {
		return claims, nil
	}

	return nil, err
}

// CheckPasswordHash compares a plaintext password with a bcrypt hash.
func CheckPasswordHash(password, hash string) (bool, error) {
	// Compare the hashed password with the plaintext password.
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	if err != nil {
		// Handle password mismatch or unexpected errors.
		if err == bcrypt.ErrMismatchedHashAndPassword {
			return false, nil
		}
		return false, err
	}

	// Passwords match.
	return true, nil
}

// GetEmailFromToken extracts the email claim from a JWT token string.
func GetEmailFromToken(tokenString string) (string, error) {
	// Parse the token to retrieve the claims.
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	if err != nil {
		return "", err
	}

	// Extract and return the email from the token claims.
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		email, ok := claims["email"].(string)
		if !ok {
			return "", errors.New("invalid email claim")
		}
		return email, nil
	} else {
		return "", errors.New("invalid token")
	}
}
