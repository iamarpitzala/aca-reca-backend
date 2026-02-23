package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/iamarpitzala/aca-reca-backend/internal/application/usecase"
	utils "github.com/iamarpitzala/aca-reca-backend/util"
)

// FinanicalQuarterHandler handles financial quarter related API endpoints.
type FinanicalQuarterHandler struct {
	financialQuarterUC *usecase.FinancialYearService
}

// NewFinanicalQuarterHandler constructs a new handler.
func NewFinanicalQuarterHandler(financialQuarterUC *usecase.FinancialYearService) *FinanicalQuarterHandler {
	return &FinanicalQuarterHandler{
		financialQuarterUC: financialQuarterUC,
	}
}

// GetAllFinancialQuarters returns all quarters for a financial year within a clinic.
// GET /api/v1/clinics/:clinicId/financial-years/:financialYearId/quarters
func (h *FinanicalQuarterHandler) GetAllFinancialQuarters(c *gin.Context) {
	clinicID := c.Param("clinicId")

	fy, err := h.financialQuarterUC.GetByClinicID(c.Request.Context(), clinicID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.JSONResponse(c, http.StatusOK, "financial quarters retrieved successfully", fy, nil)
}

// GetFinancialQuarterByID returns a single financial quarter.
// GET /api/v1/clinics/:clinicId/financial-years/:financialYearId/quarters/:id
func (h *FinanicalQuarterHandler) GetFinancialQuarterByID(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid quarter id"})
		return
	}

	quarter, err := h.financialQuarterUC.GetByID(c.Request.Context(), idInt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	utils.JSONResponse(c, http.StatusOK, "financial quarter retrieved successfully", quarter, nil)
}

// DeleteFinancialQuarter deletes a financial quarter.
// DELETE /api/v1/clinics/:clinicId/financial-years/:financialYearId/quarters/:id
func (h *FinanicalQuarterHandler) DeleteFinancialQuarter(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid quarter id"})
		return
	}
	err = h.financialQuarterUC.Delete(c.Request.Context(), idInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	utils.JSONResponse(c, http.StatusOK, "financial quarter deleted successfully", nil, nil)
}
