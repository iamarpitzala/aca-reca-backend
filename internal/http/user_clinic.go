package http

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/application/usecase"
	"github.com/iamarpitzala/aca-reca-backend/util"
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

	authUserID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	authUserUUID, ok := authUserID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	role, err := h.userClinicUC.UserRoleInClinic(c.Request.Context(), authUserUUID, req.ClinicID)
	if err != nil || role == "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied: you do not have access to this clinic"})
		return
	}
	if !strings.EqualFold(role, util.RoleOwner) {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied: only the clinic owner can add users"})
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

// GetClinicUsers retrieves all users for a clinic. Requester must have access to the clinic.
// GET /api/v1/user-clinic/clinic/:clinicId
// @Summary Get clinic's users
// @Description Get all users associated with a clinic. User must have access to the clinic.
// @Tags UserClinic
// @Accept json
// @Produce json
// @Param clinicId path string true "Clinic ID"
// @Success 200 {object} domain.H
// @Failure 400 {object} domain.H
// @Failure 403 {object} domain.H
// @Failure 500 {object} domain.H
// @Router /user-clinic/clinic/{clinicId} [get]
func (h *UserClinicHandler) GetClinicUsers(c *gin.Context) {
	clinicIDStr := c.Param("clinicId")
	clinicID, err := uuid.Parse(clinicIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid clinic ID"})
		return
	}

	authUserID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	authUserUUID, ok := authUserID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	hasAccess, err := h.userClinicUC.UserHasAccessToClinic(c.Request.Context(), authUserUUID, clinicID)
	if err != nil || !hasAccess {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied: you do not have access to this clinic"})
		return
	}

	clinicUsers, err := h.userClinicUC.GetClinicUsers(c.Request.Context(), clinicID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "clinic users retrieved successfully", "clinicUsers": clinicUsers})
}

// RemoveUserFromClinic removes a user-clinic association. Only the clinic owner can remove users.
// DELETE /api/v1/user-clinic/:id
// @Summary Remove user from clinic
// @Description Remove a user-clinic association. Only the clinic owner can perform this action.
// @Tags UserClinic
// @Accept json
// @Produce json
// @Param id path string true "User-Clinic Association ID"
// @Success 200 {object} domain.H
// @Failure 400 {object} domain.H
// @Failure 403 {object} domain.H
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

	authUserID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	authUserUUID, ok := authUserID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}

	assoc, err := h.userClinicUC.GetByID(c.Request.Context(), id)
	if err != nil || assoc == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user-clinic association not found"})
		return
	}

	role, err := h.userClinicUC.UserRoleInClinic(c.Request.Context(), authUserUUID, assoc.ClinicID)
	if err != nil || role == "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied: you do not have access to this clinic"})
		return
	}
	if !strings.EqualFold(role, util.RoleOwner) {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied: only the clinic owner can remove users"})
		return
	}

	err = h.userClinicUC.RemoveUserFromClinic(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user removed from clinic successfully", "id": id})
}
