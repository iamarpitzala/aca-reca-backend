package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/application/usecase"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
	utils "github.com/iamarpitzala/aca-reca-backend/util"
)

type BASSnapshotHandler struct {
	basUC *usecase.BASSnapshotService
}

func NewBASSnapshotHandler(basUC *usecase.BASSnapshotService) *BASSnapshotHandler {
	return &BASSnapshotHandler{basUC: basUC}
}

// Generate creates a new BAS snapshot from ledger data.
// POST /api/v1/clinic/:id/bas-snapshot/generate
func (h *BASSnapshotHandler) Generate(c *gin.Context) {
	clinicID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.ErrInvalidClinicID})
		return
	}

	var req domain.GenerateBASRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": utils.ErrUserNotAuthenticated})
		return
	}
	userIDUUID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": utils.ErrInvalidUserID})
		return
	}

	result, err := h.basUC.Generate(c.Request.Context(), clinicID, &req, userIDUUID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  utils.MsgBASSnapshotGenerated,
		"snapshot": result.Snapshot,
		"lines":    result.Lines,
	})
}

// Create creates a new BAS snapshot with manual values.
// POST /api/v1/clinic/:id/bas-snapshot
func (h *BASSnapshotHandler) Create(c *gin.Context) {
	clinicID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.ErrInvalidClinicID})
		return
	}

	var req domain.GenerateBASRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": utils.ErrUserNotAuthenticated})
		return
	}
	userIDUUID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": utils.ErrInvalidUserID})
		return
	}

	result, err := h.basUC.Generate(c.Request.Context(), clinicID, &req, userIDUUID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  utils.MsgBASSnapshotCreated,
		"snapshot": result.Snapshot,
		"lines":    result.Lines,
	})
}

// GetByID retrieves a BAS snapshot by ID.
// GET /api/v1/bas-snapshot/:id
func (h *BASSnapshotHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.ErrInvalidID})
		return
	}

	result, err := h.basUC.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  utils.MsgBASSnapshotRetrieved,
		"snapshot": result.Snapshot,
		"lines":    result.Lines,
	})
}

// GetByClinicID lists all BAS snapshots for a clinic.
// GET /api/v1/clinic/:id/bas-snapshots
func (h *BASSnapshotHandler) GetByClinicID(c *gin.Context) {
	clinicID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.ErrInvalidClinicID})
		return
	}

	snapshots, err := h.basUC.GetByClinicID(c.Request.Context(), clinicID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":   utils.MsgBASSnapshotsRetrieved,
		"snapshots": snapshots,
		"total":     len(snapshots),
	})
}

// Update modifies a DRAFT BAS snapshot.
// PUT /api/v1/bas-snapshot/:id
func (h *BASSnapshotHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.ErrInvalidID})
		return
	}

	var req domain.UpdateBASSnapshotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	snapshot, err := h.basUC.Update(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  utils.MsgBASSnapshotUpdated,
		"snapshot": snapshot,
	})
}

// Finalise changes a DRAFT snapshot to FINALISED.
// POST /api/v1/bas-snapshot/:id/finalise
func (h *BASSnapshotHandler) Finalise(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.ErrInvalidID})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": utils.ErrUserNotAuthenticated})
		return
	}
	userIDUUID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": utils.ErrInvalidUserID})
		return
	}

	snapshot, err := h.basUC.Finalise(c.Request.Context(), id, userIDUUID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  utils.MsgBASSnapshotFinalised,
		"snapshot": snapshot,
	})
}

// Lock changes a FINALISED snapshot to LOCKED.
// POST /api/v1/bas-snapshot/:id/lock
func (h *BASSnapshotHandler) Lock(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.ErrInvalidID})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": utils.ErrUserNotAuthenticated})
		return
	}
	userIDUUID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": utils.ErrInvalidUserID})
		return
	}

	snapshot, err := h.basUC.Lock(c.Request.Context(), id, userIDUUID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  utils.MsgBASSnapshotLocked,
		"snapshot": snapshot,
	})
}

// Delete soft-deletes a BAS snapshot.
// DELETE /api/v1/bas-snapshot/:id
func (h *BASSnapshotHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.ErrInvalidID})
		return
	}

	if err := h.basUC.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": utils.MsgBASSnapshotDeleted})
}
