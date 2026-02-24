package http

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/iamarpitzala/aca-reca-backend/internal/application/usecase"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain/clinic"
	utils "github.com/iamarpitzala/aca-reca-backend/util"
)

// FinanicalYearHandler handles financial year related API endpoints.
type FinanicalYearHandler struct {
	financialYearUC *usecase.FinancialYearService
}

func NewFinanicalYearHandler(financialYearUC *usecase.FinancialYearService) *FinanicalYearHandler {
	return &FinanicalYearHandler{
		financialYearUC: financialYearUC,
	}
}

// GetAllFinancialYears returns all financial years for a clinic
// GET /api/v1/clinics/:clinicId/financial-years
func (h *FinanicalYearHandler) GetAllFinancialYears(c *gin.Context) {
	clinicID := c.Param("clinicId")
	if clinicID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "clinicId required"})
		return
	}

	fys, err := h.financialYearUC.GetByClinicID(c.Request.Context(), clinicID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.JSONResponse(c, http.StatusOK, "financial years retrieved successfully", fys, nil)
}

// GetFinancialYearByID returns a single financial year by its id
// GET /api/v1/clinics/:clinicId/financial-years/:id
func (h *FinanicalYearHandler) GetFinancialYearByID(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id required"})
		return
	}
	fy, err := h.financialYearUC.GetByID(c.Request.Context(), idInt)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	utils.JSONResponse(c, http.StatusOK, "financial year retrieved successfully", fy, nil)
}

// CreateFinancialYear creates a link to a master financial year for the clinic, or creates a new master FY if financialYearId is not set.
// POST /api/v1/clinics/:clinicId/financial-years
func (h *FinanicalYearHandler) CreateFinancialYear(c *gin.Context) {
	clinicID := c.Param("clinicId")
	var req clinic.FinancialYearRequest
	if err := utils.BindAndValidate(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.FinancialYearID > 0 {
		cfy, err := h.financialYearUC.CreateLink(c.Request.Context(), clinicID, req.FinancialYearID)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		utils.JSONResponse(c, http.StatusCreated, "financial year linked successfully", cfy, nil)
		return
	}
	fy := req.ToFinancialYear()
	if err := h.financialYearUC.Create(c.Request.Context(), fy); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	utils.JSONResponse(c, http.StatusCreated, "financial year created successfully", fy, nil)
}

// UpdateFinancialYear updates an existing master financial year
// PUT /api/v1/clinics/:clinicId/financial-years/:id
func (h *FinanicalYearHandler) UpdateFinancialYear(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id required"})
		return
	}
	var req clinic.FinancialYearRequest
	if err := utils.BindAndValidate(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.ID = idInt
	fy := req.ToFinancialYear()
	if err := h.financialYearUC.Update(c.Request.Context(), fy); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	utils.JSONResponse(c, http.StatusOK, "financial year updated successfully", nil, nil)
}

// DeleteFinancialYear deletes a financial year
// DELETE /api/v1/clinics/:clinicId/financial-years/:id
func (h *FinanicalYearHandler) DeleteFinancialYear(c *gin.Context) {
	id := c.Param("id")
	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "id required"})
		return
	}
	err = h.financialYearUC.Delete(c.Request.Context(), idInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	utils.JSONResponse(c, http.StatusOK, "financial year deleted successfully", nil, nil)
}
