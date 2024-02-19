package utils

import (
	"errors" // Importing errors package for creating custom errors

	"github.com/organization_api/pkg/database/mongodb/models" // Importing models package for user struct
)

// ValidateUser validates a user's name and password.
func ValidateUser(user models.User) error {
	if err := ValidateUsername(user.Name); err != nil {
		return err
	}

	if err := ValidatePassword(user.Password); err != nil {
		return err
	}

	return nil
}

// ValidateUsername checks if a username is empty.
func ValidateUsername(username string) error {
	if username == "" {
		return errors.New("Username is required")
	}

	return nil
}

// ValidatePassword checks if a password meets the required length.
func ValidatePassword(password string) error {
	if password == "" {
		return errors.New("Password is required")
	}
	if len(password) < 8 {
		return errors.New("Password must be at least 8 characters long")
	}

	return nil
}
