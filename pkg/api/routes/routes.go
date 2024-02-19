package routes

import (
	"github.com/organization_api/pkg/api/handlers"
	"github.com/organization_api/pkg/api/middleware"

	"github.com/gin-gonic/gin"
)

// Routes sets up the application's HTTP routes.
func Routes(router *gin.Engine) {

	// Define authentication routes.
	auth := router.Group("/auth")
	{
		auth.POST("/signup", handlers.SignupHandler)              // Handle user registration
		auth.POST("/signin", handlers.SignInHandler)              // Handle user login
		auth.POST("/refresh-token", handlers.RefreshTokenHandler) // Handle token refresh
	}

	// Define organization routes, secured with authentication.
	organization := router.Group("/api")
	organization.Use(middleware.AuthMiddleware())
	{
		organization.POST("organization", handlers.CreateOrganizationHandler)                                                  // Handle organization creation
		organization.GET("/organization/:organization_id", middleware.InviteMiddleware(), handlers.GetOrganizationByIdHandler) // Handle organization retrieval with invitation check
		organization.GET("/organization", handlers.GetAllOrganizationsHandler)                                                 // Handle all organizations retrieval
		organization.PUT("/organization/:organization_id", handlers.UpdateOrganizationHandler)                                 // Handle organization update
		organization.DELETE("/organization/:organization_id", handlers.DeleteOrganizationHandler)                              // Handle organization deletion
		organization.POST("/organization/:organization_id/invite", handlers.InviteUserToOrganizationHandler)                   // Handle organization invitation
	}
}
