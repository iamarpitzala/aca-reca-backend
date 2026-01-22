package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
	"github.com/iamarpitzala/aca-reca-backend/internal/service"
)

type FinancialCalculationHandler struct {
	calculationService *service.FinancialCalculationService
}

func NewFinancialCalculationHandler(calculationService *service.FinancialCalculationService) *FinancialCalculationHandler {
	return &FinancialCalculationHandler{
		calculationService: calculationService,
	}
}

// CalculateFinancial performs a financial calculation
// POST /api/v1/financial-calculation/calculate
// @Summary Calculate financial values
// @Description Perform financial calculation based on form configuration and input values
// @Tags FinancialCalculation
// @Accept json
// @Produce json
// @Param request body object true "Calculation request" example({"formId":"uuid","input":{}})
// @Success 200 {object} domain.H
// @Failure 400 {object} domain.H
// @Failure 500 {object} domain.H
// @Router /financial-calculation/calculate [post]
func (h *FinancialCalculationHandler) CalculateFinancial(c *gin.Context) {
	var req struct {
		FormID uuid.UUID               `json:"formId" binding:"required"`
		Input  domain.CalculationInput `json:"input" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from context if available
	var userID *uuid.UUID
	if userIDVal, exists := c.Get("user_id"); exists {
		if uid, ok := userIDVal.(uuid.UUID); ok {
			userID = &uid
		}
	}

	result, basMapping, err := h.calculationService.CalculateFinancial(c.Request.Context(), req.FormID, req.Input, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":    "calculation completed successfully",
		"result":     result,
		"basMapping": basMapping,
	})
}

// GetCalculationHistory retrieves calculation history for a form
// GET /api/v1/financial-calculation/history/:formId
// @Summary Get calculation history
// @Description Get calculation history for a financial form
// @Tags FinancialCalculation
// @Accept json
// @Produce json
// @Param formId path string true "Financial Form ID"
// @Success 200 {object} domain.H
// @Failure 400 {object} domain.H
// @Failure 500 {object} domain.H
// @Router /financial-calculation/history/{formId} [get]
func (h *FinancialCalculationHandler) GetCalculationHistory(c *gin.Context) {
	formIDStr := c.Param("formId")
	formID, err := uuid.Parse(formIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid form ID"})
		return
	}

	history, err := h.calculationService.GetCalculationHistory(c.Request.Context(), formID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "calculation history retrieved successfully", "history": history})
}
