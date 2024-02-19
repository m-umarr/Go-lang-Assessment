package handlers

import (
	"net/http"

	"github.com/organization_api/pkg/database/mongodb/models"
	"github.com/organization_api/pkg/database/mongodb/repository"
	"github.com/organization_api/pkg/utils"

	"github.com/gin-gonic/gin"
)

// SignInHandler authenticates a user based on their AuthCreds.
func SignInHandler(c *gin.Context) {
	// Parse the incoming JSON payload containing user AuthCreds.
	var AuthCreds *models.AuthCreds
	newrepo := repository.NewUserRepository()

	err := c.ShouldBindJSON(&AuthCreds)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
		return
	}

	// Check that both username and password are provided.
	if AuthCreds.Email == "" || AuthCreds.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username and password cannot be empty"})
		return
	}

	// Find the user by email in the database.
	userFound, err := newrepo.FindUserByEmail(AuthCreds.Email)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid AuthCreds"})
		return
	}

	// Verify the provided password against the stored hash.
	isMatch, err := utils.CheckPasswordHash(AuthCreds.Password, userFound.Password)
	if err != nil || !isMatch {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid AuthCreds"})
		return
	}

	// Generate authentication tokens for the authenticated user.
	access_token, refresh_token, err := utils.GenerateTokens(userFound.Name, userFound.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
		return
	}

	// Respond with success message and tokens.
	c.JSON(http.StatusOK, models.AuthResponse{
		Message:      "SignIn successful",
		AccessToken:  access_token,
		RefreshToken: refresh_token,
	})
}

// SignupHandler handles the creation of a new user account.
func SignupHandler(c *gin.Context) {
	// Parse and validate the incoming JSON payload.
	var user models.User
	repo := repository.NewUserRepository()

	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
		return
	}

	// Validate the user data before proceeding.
	if err := utils.ValidateUser(user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Hash the user's password for secure storage.
	hash, err := utils.HashPassword(user.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}
	user.Password = hash

	// Attempt to create the user in the database.
	createdUser, err := repo.CreateUser(&user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "User not created"})
		return
	}

	// Generate authentication tokens for the newly created user.
	access_token, refresh_token, err := utils.GenerateTokens(createdUser.Name, createdUser.Email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
		return
	}

	// Respond with success message and tokens.
	c.JSON(http.StatusCreated, models.AuthResponse{
		Message:      "User created successfully",
		AccessToken:  access_token,
		RefreshToken: refresh_token,
	})
}

// RefreshTokenHandler issues new access and refresh tokens using a valid refresh token.
func RefreshTokenHandler(c *gin.Context) {
	// Parse the incoming JSON payload containing the refresh token.
	var request models.RefreshToken

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
		return
	}

	// Verify the refresh token and extract the associated username and email.
	username, email, err := utils.VerifyRefreshToken(request.Token)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	// Generate new access and refresh tokens for the user.
	accessToken, refreshToken, err := utils.GenerateTokens(username, email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate tokens"})
		return
	}

	// Prepare the response with the new tokens.
	response := models.AuthResponse{
		Message:      "Tokens refreshed",
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}

	c.JSON(http.StatusOK, response)
}
