package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iamarpitzala/aca-reca-backend/internal/application/usecase"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain/user"
	utils "github.com/iamarpitzala/aca-reca-backend/util"
)

type UserHandler struct {
	authUC       *usecase.AuthService
	userClinicUC *usecase.UserClinicService
}

func NewUserHandler(authUC *usecase.AuthService, userClinicUC *usecase.UserClinicService) *UserHandler {
	return &UserHandler{
		authUC:       authUC,
		userClinicUC: userClinicUC,
	}
}

// GetMe returns the current authenticated user from JWT
// GET /api/v1/users/me
// @Summary Get current user (me)
// @Description Get current authenticated user from JWT token
// @Tags User
// @Accept json
// @Produce json
// @Success 200 {object} domain.User
// @Failure 401 {object} domain.H
// @Failure 500 {object} domain.H
// @Router /users/me [get]
func (h *UserHandler) GetMe(c *gin.Context) {
	userID, ok := GetAuthUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": utils.ErrUserNotAuthenticated})
		return
	}

	user, err := h.authUC.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// GetCurrentUser returns the current authenticated user
// GET /api/v1/users/:user_id
// @Summary Get current user
// @Description Get current user by user ID
// @Tags User
// @Accept json
// @Produce json
// @Param userId path string true "User ID"
// @Success 200 {object} domain.User
// @Failure 400 {object} domain.H
// @Failure 500 {object} domain.H
// @Router /users/{userId} [get]
func (h *UserHandler) GetCurrentUser(c *gin.Context) {
	userID := c.Param("userId")
	if !RequireClinicAccess(c, h.userClinicUC, userID) {
		return
	}

	user, err := h.authUC.GetUserByID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// UpdateCurrentUser updates the current authenticated user
// PUT /api/v1/users/:userId
// @Summary Update current user
// @Description Update current user by user ID
// @Tags User
// @Accept json
// @Produce json
// @Param userId path string true "User ID"
// @Param updateUserRequest body domain.UpdateUserRequest true "Update user request"
// @Success 200 {object} domain.User
// @Failure 400 {object} domain.H
// @Failure 500 {object} domain.H
// @Router /users/{userId} [put]
func (h *UserHandler) UpdateCurrentUser(c *gin.Context) {
	userID := c.Param("userId")
	if !RequireClinicAccess(c, h.userClinicUC, userID) {
		return
	}

	var req user.UserRequest

	if err := utils.BindAndValidate(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.authUC.UpdateUser(c.Request.Context(), userID, req.FirstName, req.LastName, req.Phone)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, user)
}

// GetActiveSessions returns all active sessions for the current user
// GET /api/v1/users/:userId/sessions
// @Summary Get active sessions
// @Description Get active sessions by user ID
// @Tags User
// @Accept json
// @Produce json
// @Param userId path string true "User ID"
// @Success 200 {array} domain.Session
// @Failure 400 {object} domain.H
// @Failure 500 {object} domain.H
// @Router /users/{userId}/sessions [get]
func (h *UserHandler) GetActiveSessions(c *gin.Context) {
	userID := c.Param("userId")
	if !RequireClinicAccess(c, h.userClinicUC, userID) {
		return
	}
	sessions, err := h.authUC.GetUserSessions(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "sessions retrieved successfully", "sessions": sessions})
}

// RevokeSession revokes a specific session
// DELETE /api/v1/users/:userId/sessions/:sessionId
// @Summary Revoke a session
// @Description Revoke a session by session ID and user ID
// @Tags User
// @Accept json
// @Produce json
// @Param userId path string true "User ID"
// @Param sessionId path string true "Session ID"
// @Success 200 {object} domain.H
// @Failure 400 {object} domain.H
// @Failure 500 {object} domain.H
// @Router /users/{userId}/sessions/{sessionId} [delete]
func (h *UserHandler) RevokeSession(c *gin.Context) {
	userID := c.Param("userId")
	if !RequireClinicAccess(c, h.userClinicUC, userID) {
		return
	}
	sessionID := c.Param("sessionId")

	// Verify session belongs to user
	sessions, err := h.authUC.GetUserSessions(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	found := false
	for _, session := range sessions {
		if session.ID == sessionID {
			found = true
			break
		}
	}

	if !found {
		c.JSON(http.StatusForbidden, gin.H{"error": utils.ErrSessionNotFound})
		return
	}

	if err := h.authUC.RevokeSession(c.Request.Context(), sessionID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "session revoked successfully", "sessionId": sessionID})
}
