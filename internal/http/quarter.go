package http

import (
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/application/usecase"
)

type QuarterHandler struct {
	QuarterUC *usecase.QuarterService
}

func NewQuarterHandler(QuarterUC *usecase.QuarterService) *QuarterHandler {
	return &QuarterHandler{
		QuarterUC: QuarterUC,
	}
}

// CalculateForClinic calculates quarters for a clinic based on financial settings
// @Summary Calculate quarters for a clinic
// @Description Get calculated quarters for a clinic based on its financial year settings (system-driven)
// @Tags Quarter
// @Accept json
// @Produce json
// @Param id path string true "Clinic ID"
// @Param yearsBack query int false "Years to look back" default(1)
// @Param yearsForward query int false "Years to look forward" default(1)
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /clinic/{id}/quarters [get]
func (h *QuarterHandler) CalculateForClinic(c *gin.Context) {
	clinicIDStr := c.Param("id")
	clinicID, err := uuid.Parse(clinicIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid clinic ID"})
		return
	}

	// Parse query parameters with defaults
	yearsBack := 1
	if yearsBackStr := c.Query("yearsBack"); yearsBackStr != "" {
		yearsBack, err = strconv.Atoi(yearsBackStr)
		if err != nil || yearsBack < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid yearsBack parameter"})
			return
		}
	}

	yearsForward := 1
	if yearsForwardStr := c.Query("yearsForward"); yearsForwardStr != "" {
		yearsForward, err = strconv.Atoi(yearsForwardStr)
		if err != nil || yearsForward < 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid yearsForward parameter"})
			return
		}
	}

	quarters, err := h.QuarterUC.CalculateQuartersForClinic(
		c.Request.Context(),
		clinicID,
		yearsBack,
		yearsForward,
	)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": quarters,
		"message": "quarters calculated successfully",
	})
}

// GetQuarterForDate finds the quarter containing a specific date for a clinic
// @Summary Get quarter for date
// @Description Find which quarter contains a specific date for a clinic
// @Tags Quarter
// @Accept json
// @Produce json
// @Param id path string true "Clinic ID"
// @Param date query string true "Date (RFC3339 format, e.g., 2025-02-13T00:00:00Z)"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /clinic/{id}/quarter/date [get]
func (h *QuarterHandler) GetQuarterForDate(c *gin.Context) {
	clinicIDStr := c.Param("id")
	clinicID, err := uuid.Parse(clinicIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid clinic ID"})
		return
	}

	dateStr := c.Query("date")
	if dateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "date parameter is required"})
		return
	}

	date, err := time.Parse(time.RFC3339, dateStr)
	if err != nil {
		// Try parsing as date only
		date, err = time.Parse("2006-01-02", dateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date format. Use RFC3339 (2006-01-02T15:04:05Z07:00) or date (2006-01-02)"})
			return
		}
	}

	quarter, err := h.QuarterUC.GetQuarterForDate(c.Request.Context(), clinicID, date)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": quarter,
		"message": "quarter found successfully",
	})
}
