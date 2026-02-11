package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/application/usecase"
)

type UserClinicHandler struct {
	userClinicUC *usecase.UserClinicService
}

func NewUserClinicHandler(userClinicUC *usecase.UserClinicService) *UserClinicHandler {
	return &UserClinicHandler{
		userClinicUC: userClinicUC,
	}
}

// AssociateUserWithClinic associates a user with a clinic
// POST /api/v1/user-clinic
// @Summary Associate user with clinic
// @Description Associate a user with a clinic
// @Tags UserClinic
// @Accept json
// @Produce json
// @Param request body object true "Association request" example({"userId":"uuid","clinicId":"uuid","role":"owner"})
// @Success 201 {object} domain.H
// @Failure 400 {object} domain.H
// @Failure 500 {object} domain.H
// @Router /user-clinic [post]
func (h *UserClinicHandler) AssociateUserWithClinic(c *gin.Context) {
	var req struct {
		UserID   uuid.UUID `json:"userId" binding:"required"`
		ClinicID uuid.UUID `json:"clinicId" binding:"required"`
		Role     string    `json:"role"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userClinic, err := h.userClinicUC.AssociateUserWithClinic(c.Request.Context(), req.UserID, req.ClinicID, req.Role)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "user associated with clinic successfully", "userClinic": userClinic})
}

// GetUserClinics retrieves all clinics for a user
// GET /api/v1/user-clinic/user/:userId
// @Summary Get user's clinics
// @Description Get all clinics associated with a user. Users can only fetch their own clinics.
// @Tags UserClinic
// @Accept json
// @Produce json
// @Param userId path string true "User ID"
// @Success 200 {object} domain.H
// @Failure 400 {object} domain.H
// @Failure 403 {object} domain.H
// @Failure 500 {object} domain.H
// @Router /user-clinic/user/{userId} [get]
func (h *UserClinicHandler) GetUserClinics(c *gin.Context) {
	userIDStr := c.Param("userId")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid user ID"})
		return
	}

	// Verify the requesting user can only fetch their own clinics
	authUserID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	authUserUUID, ok := authUserID.(uuid.UUID)
	if !ok || authUserUUID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied: you can only view your own clinics"})
		return
	}

	userClinics, err := h.userClinicUC.GetUserClinics(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user clinics retrieved successfully", "userClinics": userClinics})
}

// GetClinicUsers retrieves all users for a clinic
// GET /api/v1/user-clinic/clinic/:clinicId
// @Summary Get clinic's users
// @Description Get all users associated with a clinic
// @Tags UserClinic
// @Accept json
// @Produce json
// @Param clinicId path string true "Clinic ID"
// @Success 200 {object} domain.H
// @Failure 400 {object} domain.H
// @Failure 500 {object} domain.H
// @Router /user-clinic/clinic/{clinicId} [get]
func (h *UserClinicHandler) GetClinicUsers(c *gin.Context) {
	clinicIDStr := c.Param("clinicId")
	clinicID, err := uuid.Parse(clinicIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid clinic ID"})
		return
	}

	clinicUsers, err := h.userClinicUC.GetClinicUsers(c.Request.Context(), clinicID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "clinic users retrieved successfully", "clinicUsers": clinicUsers})
}

// RemoveUserFromClinic removes a user-clinic association
// DELETE /api/v1/user-clinic/:id
// @Summary Remove user from clinic
// @Description Remove a user-clinic association
// @Tags UserClinic
// @Accept json
// @Produce json
// @Param id path string true "User-Clinic Association ID"
// @Success 200 {object} domain.H
// @Failure 400 {object} domain.H
// @Failure 404 {object} domain.H
// @Failure 500 {object} domain.H
// @Router /user-clinic/{id} [delete]
func (h *UserClinicHandler) RemoveUserFromClinic(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid association ID"})
		return
	}

	err = h.userClinicUC.RemoveUserFromClinic(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user removed from clinic successfully", "id": id})
}
