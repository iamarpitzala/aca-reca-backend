package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iamarpitzala/aca-reca-backend/internal/application/usecase"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain/coa"
	utils "github.com/iamarpitzala/aca-reca-backend/util"
)

type ArrangmentHandler struct {
	arrangementUC *usecase.ArrangementService
}

func NewArrangmentHandler(arrangementUC *usecase.ArrangementService) *ArrangmentHandler {
	return &ArrangmentHandler{arrangementUC: arrangementUC}
}

// Create creates a new arrangement
// POST /api/v1/arrangements
func (h *ArrangmentHandler) Create(c *gin.Context) {
	userID, ok := GetAuthUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": utils.ErrUserNotAuthenticated})
		return
	}
	var req coa.ArrangementRequest
	if err := utils.BindAndValidate(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := h.arrangementUC.Create(c.Request.Context(), userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": utils.ErrArrangementFailed})
		return
	}
	utils.JSONResponse(c, http.StatusCreated, utils.MsgArrangementCreated, resp, nil)
}

// GetByID returns an arrangement by id (only if owned by current user)
// GET /api/v1/arrangements/:id
func (h *ArrangmentHandler) GetByID(c *gin.Context) {
	userID, ok := GetAuthUserID(c)
	if !ok {
		return
	}
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.ErrArrangementIDRequired})
		return
	}
	resp, err := h.arrangementUC.GetByID(c.Request.Context(), id, userID)
	if err != nil {
		if err == usecase.ErrArrangementNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": utils.ErrArrangementNotFound})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": utils.ErrArrangementFailed})
		return
	}
	utils.JSONResponse(c, http.StatusOK, utils.MsgArrangementRetrieved, resp, nil)
}

// List returns all arrangements for the current user
// GET /api/v1/arrangements
func (h *ArrangmentHandler) ListByUserID(c *gin.Context) {
	userID, ok := GetAuthUserID(c)
	if !ok {
		return
	}
	list, err := h.arrangementUC.ListByUserID(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": utils.ErrArrangementFailed})
		return
	}
	utils.JSONResponse(c, http.StatusOK, utils.MsgArrangementsRetrieved, list, nil)
}

// Update updates an arrangement (only if owned by current user)
// PUT /api/v1/arrangements/:id
func (h *ArrangmentHandler) Update(c *gin.Context) {
	userID, ok := GetAuthUserID(c)
	if !ok {
		return
	}
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.ErrArrangementIDRequired})
		return
	}
	var req coa.ArrangementRequest
	if err := utils.BindAndValidate(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := h.arrangementUC.Update(c.Request.Context(), id, userID, &req)
	if err != nil {
		if err == usecase.ErrArrangementNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": utils.ErrArrangementNotFound})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": utils.ErrArrangementFailed})
		return
	}
	utils.JSONResponse(c, http.StatusOK, utils.MsgArrangementUpdated, resp, nil)
}

// Delete soft-deletes an arrangement (only if owned by current user)
// DELETE /api/v1/arrangements/:id
func (h *ArrangmentHandler) Delete(c *gin.Context) {
	userID, ok := GetAuthUserID(c)
	if !ok {
		return
	}
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.ErrArrangementIDRequired})
		return
	}
	if err := h.arrangementUC.Delete(c.Request.Context(), id, userID); err != nil {
		if err == usecase.ErrArrangementNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": utils.ErrArrangementNotFound})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": utils.ErrArrangementFailed})
		return
	}
	utils.JSONResponse(c, http.StatusOK, utils.MsgArrangementDeleted, nil, nil)
}
