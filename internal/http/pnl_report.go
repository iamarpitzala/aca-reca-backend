package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iamarpitzala/aca-reca-backend/internal/application/usecase"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
	utils "github.com/iamarpitzala/aca-reca-backend/util"
)

type PnlReportHandler struct {
	pnlUC *usecase.PnlReportService
}

func NewPnlReportHandler(pnlUC *usecase.PnlReportService) *PnlReportHandler {
	return &PnlReportHandler{pnlUC: pnlUC}
}

// Generate creates a new P&L report from ledger data
// POST /api/v1/reports/pnl/generate
func (h *PnlReportHandler) Generate(c *gin.Context) {
	var req domain.GeneratePnlRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, ok := GetAuthUserID(c)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": utils.ErrUserNotAuthenticated})
		return
	}

	result, err := h.pnlUC.Generate(c.Request.Context(), &req, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": utils.MsgPnlReportGenerated,
		"report":  result.Report,
		"lines":   result.Lines,
	})
}

// GetByID retrieves a P&L report by ID
// GET /api/v1/reports/pnl/:id
func (h *PnlReportHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	result, err := h.pnlUC.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": utils.MsgPnlReportRetrieved,
		"report":  result.Report,
		"lines":   result.Lines,
	})
}

// GetByClinicID lists all P&L reports for a clinic
// GET /api/v1/reports/pnl/clinic/:clinicId
func (h *PnlReportHandler) GetByClinicID(c *gin.Context) {
	clinicID := c.Param("clinicId")

	reports, err := h.pnlUC.GetByClinicID(c.Request.Context(), clinicID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": utils.MsgPnlReportsRetrieved,
		"reports": reports,
		"total":   len(reports),
	})
}

// Finalize changes a DRAFT report to FINAL
// POST /api/v1/reports/pnl/:id/finalize
func (h *PnlReportHandler) Finalize(c *gin.Context) {
	id := c.Param("id")

	report, err := h.pnlUC.Finalize(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": utils.MsgPnlReportFinalized,
		"report":  report,
	})
}

// Regenerate re-aggregates ledger data for a DRAFT report
// POST /api/v1/reports/pnl/:id/regenerate
func (h *PnlReportHandler) Regenerate(c *gin.Context) {
	id := c.Param("id")

	result, err := h.pnlUC.Regenerate(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": utils.MsgPnlReportRegenerated,
		"report":  result.Report,
		"lines":   result.Lines,
	})
}

// Delete soft-deletes a P&L report
// DELETE /api/v1/reports/pnl/:id
func (h *PnlReportHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.pnlUC.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": utils.MsgPnlReportDeleted})
}
