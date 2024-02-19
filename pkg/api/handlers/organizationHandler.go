package handlers

import (
	"fmt"
	"net/http"

	"github.com/organization_api/pkg/database/mongodb/models"
	"github.com/organization_api/pkg/database/mongodb/repository"

	"github.com/gin-gonic/gin"
)

// GetAllOrganizationsHandler lists all organizations.
func GetAllOrganizationsHandler(c *gin.Context) {
	repo := repository.NewOrganizationRepo()
	organizations, err := repo.GetAllOrganizations()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch organizations"})
		return
	}

	// Respond with a success message and the list of organizations.
	c.JSON(http.StatusOK, organizations)
}

// CreateOrganizationHandler creates a new organization record.
func CreateOrganizationHandler(c *gin.Context) {
	var org models.Organization
	repo := repository.NewOrganizationRepo()

	err := c.ShouldBindJSON(&org)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
		return
	}
	orgID, err := repo.CreateOrganization(&org)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create organization"})
		return
	}
	// Respond with a success message and the organization ID.
	c.JSON(http.StatusCreated, gin.H{"organization_id": orgID})
}

// GetOrganizationByIdHandler retrieves an organization by its ID.
func GetOrganizationByIdHandler(c *gin.Context) {
	organizationID := c.Param("organization_id")

	repo := repository.NewOrganizationRepo()
	organization, err := repo.GetOrganizationById(organizationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch organization details"})
		return
	}

	// Respond with a success message and the organization details.
	c.JSON(http.StatusOK, models.Organization{
		Id:          organization.Id,
		Name:        organization.Name,
		Description: organization.Description,
	})
}

// UpdateOrganizationHandler updates an existing organization's details.
func UpdateOrganizationHandler(c *gin.Context) {
	organizationID := c.Param("organization_id")

	var updateData models.OrganizationUpdate
	if err := c.ShouldBindJSON(&updateData); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	fmt.Println("req body", updateData)
	repo := repository.NewOrganizationRepo()

	organization, err := repo.UpdateOrganization(organizationID, &updateData)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update organization"})
		return
	}

	// Respond with a success message and the updated organization details.
	c.JSON(http.StatusOK, gin.H{
		"organization_id": organization.Id,
		"name":            organization.Name,
		"description":     organization.Description,
	})
}

// DeleteOrganizationHandler removes an organization from the database.
func DeleteOrganizationHandler(c *gin.Context) {
	organizationID := c.Param("organization_id")

	repo := repository.NewOrganizationRepo()

	err := repo.DeleteOrganization(organizationID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete organization"})
		return
	}
	// Respond with a success message.
	c.JSON(http.StatusOK, gin.H{"message": "Organization deleted successfully"})
}

// InviteUserToOrganizationHandler sends an invitation to join an organization.
func InviteUserToOrganizationHandler(c *gin.Context) {
	organizationID := c.Param("organization_id")
	var requestBody models.InviterequestBody

	if err := c.ShouldBindJSON(&requestBody); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON payload"})
		return
	}

	repo := repository.NewOrganizationRepo()
	// Call the InviteUserToOrganization method in the repository
	err := repo.InviteUserToOrganization(organizationID, requestBody.UserEmail)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to invite user to organization"})
		return
	}

	// Respond with a success message.
	c.JSON(http.StatusOK, gin.H{"message": "User invited to organization"})
}
