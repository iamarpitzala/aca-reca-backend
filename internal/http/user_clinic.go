package http

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
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
		UserID   string `json:"userId" binding:"required"`
		ClinicID string `json:"clinicId" binding:"required"`
		Role     string `json:"role"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	authUserID, ok := GetAuthUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": util.ErrUserNotAuthenticated})
		return
	}

	role, err := h.userClinicUC.UserRoleInClinic(c.Request.Context(), authUserID, req.ClinicID)
	if err != nil || role == "" {
		c.JSON(http.StatusForbidden, gin.H{"error": util.ErrAccessDenied})
		return
	}
	if !strings.EqualFold(role, util.RoleOwner) {
		c.JSON(http.StatusForbidden, gin.H{"error": util.ErrAccessDeniedOwnerCanAddUsers})
		return
	}

	userClinic, err := h.userClinicUC.AssociateUserWithClinic(c.Request.Context(), req.UserID, req.ClinicID, req.Role)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": util.MsgUserAssociatedWithClinic, "userClinic": userClinic})
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
	userID := c.Param("userId")
	if !RequireClinicAccess(c, h.userClinicUC, userID) {
		return
	}
	userClinics, err := h.userClinicUC.GetUserClinics(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": util.MsgUserClinicsRetrievedSuccessfully, "userClinics": userClinics})
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
	clinicID := c.Param("clinicId")
	if !RequireClinicAccess(c, h.userClinicUC, clinicID) {
		return
	}

	authUserID, ok := GetAuthUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": util.ErrUserNotAuthenticated})
		return
	}

	hasAccess, err := h.userClinicUC.UserHasAccessToClinic(c.Request.Context(), authUserID, clinicID)
	if err != nil || !hasAccess {
		c.JSON(http.StatusForbidden, gin.H{"error": util.ErrAccessDenied})
		return
	}

	clinicUsers, err := h.userClinicUC.GetClinicUsers(c.Request.Context(), clinicID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": util.MsgClinicUsersRetrievedSuccessfully, "clinicUsers": clinicUsers})
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
	id := c.Param("id")
	if !RequireClinicAccess(c, h.userClinicUC, id) {
		return
	}

	authUserID, ok := GetAuthUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": util.ErrUserNotAuthenticated})
		return
	}

	assoc, err := h.userClinicUC.GetByID(c.Request.Context(), id)
	if err != nil || assoc == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": util.ErrUserClinicAssociationNotFound})
		return
	}

	role, err := h.userClinicUC.UserRoleInClinic(c.Request.Context(), authUserID, assoc.ClinicID)
	if err != nil || role == "" {
		c.JSON(http.StatusForbidden, gin.H{"error": util.ErrAccessDenied})
		return
	}
	if !strings.EqualFold(role, util.RoleOwner) {
		c.JSON(http.StatusForbidden, gin.H{"error": util.ErrAccessDeniedOwnerCanRemoveUsers})
		return
	}

	err = h.userClinicUC.RemoveUserFromClinic(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": util.MsgUserRemovedFromClinic, "id": id})
}
