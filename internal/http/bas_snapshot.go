package http

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/application/usecase"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
)

type BASSnapshotHandler struct {
	basUC *usecase.BASSnapshotService
}

func NewBASSnapshotHandler(basUC *usecase.BASSnapshotService) *BASSnapshotHandler {
	return &BASSnapshotHandler{
		basUC: basUC,
	}
}

// getAuthUserID extracts user ID from JWT context
func (h *BASSnapshotHandler) getAuthUserID(c *gin.Context) (uuid.UUID, bool) {
	userIDVal, exists := c.Get("user_id")
	if !exists {
		return uuid.Nil, false
	}
	userUUID, ok := userIDVal.(uuid.UUID)
	return userUUID, ok
}

// CreateBASSnapshot creates a new BAS snapshot (draft)
// POST /api/v1/clinic/:id/bas-snapshot
func (h *BASSnapshotHandler) CreateBASSnapshot(c *gin.Context) {
	clinicIDStr := c.Param("id")
	clinicID, err := uuid.Parse(clinicIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid clinic ID"})
		return
	}
	
	var req domain.BASSnapshotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	snapshot, err := h.basUC.Create(c.Request.Context(), clinicID, &req)
	if err != nil {
		if err == usecase.ErrBASAlreadyFinalised {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusCreated, snapshot)
}

// GetBASSnapshot retrieves a BAS snapshot by ID
// GET /api/v1/bas-snapshot/:id
func (h *BASSnapshotHandler) GetBASSnapshot(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid BAS snapshot ID"})
		return
	}
	
	snapshot, err := h.basUC.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, snapshot)
}

// GetBASSnapshotsByClinic retrieves all BAS snapshots for a clinic
// GET /api/v1/clinic/:id/bas-snapshots
func (h *BASSnapshotHandler) GetBASSnapshotsByClinic(c *gin.Context) {
	clinicIDStr := c.Param("id")
	clinicID, err := uuid.Parse(clinicIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid clinic ID"})
		return
	}
	
	snapshots, err := h.basUC.GetByClinicID(c.Request.Context(), clinicID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"snapshots": snapshots})
}

// FinaliseBAS finalises a BAS snapshot
// POST /api/v1/bas-snapshot/:id/finalise
func (h *BASSnapshotHandler) FinaliseBAS(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid BAS snapshot ID"})
		return
	}
	
	userID, ok := h.getAuthUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	
	snapshot, err := h.basUC.Finalise(c.Request.Context(), id, userID)
	if err != nil {
		if err == usecase.ErrBASAlreadyFinalised {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, snapshot)
}

// LockBAS locks a BAS snapshot
// POST /api/v1/bas-snapshot/:id/lock
func (h *BASSnapshotHandler) LockBAS(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid BAS snapshot ID"})
		return
	}
	
	snapshot, err := h.basUC.Lock(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, snapshot)
}

// UpdateBASSnapshot updates a BAS snapshot (only if not finalised)
// PUT /api/v1/bas-snapshot/:id
func (h *BASSnapshotHandler) UpdateBASSnapshot(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid BAS snapshot ID"})
		return
	}
	
	var req domain.BASSnapshotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	snapshot, err := h.basUC.Update(c.Request.Context(), id, &req)
	if err != nil {
		if err == usecase.ErrBASAlreadyFinalised {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, snapshot)
}

// GetConsolidatedGSTSummary retrieves consolidated GST data across multiple clinics
// POST /api/v1/reports/consolidated-gst-summary
func (h *BASSnapshotHandler) GetConsolidatedGSTSummary(c *gin.Context) {
	var req struct {
		ClinicIDs  []string `json:"clinicIds" binding:"required"`
		PeriodStart string  `json:"periodStart" binding:"required"`
		PeriodEnd   string  `json:"periodEnd" binding:"required"`
	}
	
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	clinicIDs := make([]uuid.UUID, len(req.ClinicIDs))
	for i, idStr := range req.ClinicIDs {
		id, err := uuid.Parse(idStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid clinic ID in list"})
			return
		}
		clinicIDs[i] = id
	}
	
	periodStart, err := time.Parse("2006-01-02", req.PeriodStart)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid periodStart format (expected YYYY-MM-DD)"})
		return
	}
	
	periodEnd, err := time.Parse("2006-01-02", req.PeriodEnd)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid periodEnd format (expected YYYY-MM-DD)"})
		return
	}
	
	snapshots, err := h.basUC.GetConsolidatedGSTSummary(c.Request.Context(), clinicIDs, periodStart, periodEnd)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, gin.H{"snapshots": snapshots})
}
