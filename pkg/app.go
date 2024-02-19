package pkg

import (
	"github.com/organization_api/pkg/api/routes"

	"github.com/gin-gonic/gin"
)

// Run starts the web server and registers the API routes.
func Init() {
	// Initialize the Gin router with default middleware.
	router := gin.Default()
	// Register the API routes with the router.
	routes.Routes(router)

	// Start the web server on port  8080.
	router.Run(":8080")
}
