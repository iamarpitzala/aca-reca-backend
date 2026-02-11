package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/application/usecase"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
)

type ClinicFinancialSettingsHandler struct {
	settingsUC *usecase.ClinicFinancialSettingsService
}

func NewClinicFinancialSettingsHandler(settingsUC *usecase.ClinicFinancialSettingsService) *ClinicFinancialSettingsHandler {
	return &ClinicFinancialSettingsHandler{
		settingsUC: settingsUC,
	}
}

// GetFinancialSettings retrieves financial settings for a clinic
// GET /api/v1/clinic/:id/financial-settings
func (h *ClinicFinancialSettingsHandler) GetFinancialSettings(c *gin.Context) {
	clinicIDStr := c.Param("id")
	clinicID, err := uuid.Parse(clinicIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid clinic ID"})
		return
	}
	
	settings, err := h.settingsUC.GetByClinicID(c.Request.Context(), clinicID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, settings)
}

// CreateOrUpdateFinancialSettings creates or updates financial settings for a clinic
// PUT /api/v1/clinic/:id/financial-settings
func (h *ClinicFinancialSettingsHandler) CreateOrUpdateFinancialSettings(c *gin.Context) {
	clinicIDStr := c.Param("id")
	clinicID, err := uuid.Parse(clinicIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid clinic ID"})
		return
	}
	
	var req domain.ClinicFinancialSettingsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	settings, err := h.settingsUC.CreateOrUpdate(c.Request.Context(), clinicID, &req)
	if err != nil {
		if err == usecase.ErrFinancialSettingsLocked {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	
	c.JSON(http.StatusOK, settings)
}
