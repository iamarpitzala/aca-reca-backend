package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
	"github.com/iamarpitzala/aca-reca-backend/internal/service"
	utils "github.com/iamarpitzala/aca-reca-backend/util"
)

type FinancialFormHandler struct {
	financialFormService *service.FinancialFormService
}

func NewFinancialFormHandler(financialFormService *service.FinancialFormService) *FinancialFormHandler {
	return &FinancialFormHandler{
		financialFormService: financialFormService,
	}
}

// CreateFinancialForm creates a new financial form
// POST /api/v1/financial-form
// @Summary Create a new financial form
// @Description Create a new financial form with configuration
// @Tags FinancialForm
// @Accept json
// @Produce json
// @Param form body domain.FinancialFormRequest true "Financial form information"
// @Success 201 {object} domain.H
// @Failure 400 {object} domain.H
// @Failure 500 {object} domain.H
// @Router /financial-form [post]
func (h *FinancialFormHandler) CreateFinancialForm(c *gin.Context) {
	var form domain.FinancialFormRequest
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.financialFormService.CreateFinancialForm(c.Request.Context(), &form)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	utils.JSONResponse(c, http.StatusCreated, "financial form created successfully", nil, nil)
}

// GetFinancialForm retrieves a financial form by ID
// GET /api/v1/financial-form/:id
// @Summary Get financial form by ID
// @Description Retrieve a financial form by ID
// @Tags FinancialForm
// @Accept json
// @Produce json
// @Param id path string true "Financial Form ID"
// @Success 200 {object} domain.H
// @Failure 400 {object} domain.H
// @Failure 404 {object} domain.H
// @Failure 500 {object} domain.H
// @Router /financial-form/{id} [get]
func (h *FinancialFormHandler) GetFinancialForm(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid financial form ID"})
		return
	}

	form, err := h.financialFormService.GetFinancialFormByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	utils.JSONResponse(c, http.StatusOK, "financial form retrieved successfully", form, nil)
}

// GetFinancialFormsByClinic retrieves all financial forms for a clinic
// GET /api/v1/financial-form/clinic/:clinicId
// @Summary Get clinic's financial forms
// @Description Get all financial forms for a clinic
// @Tags FinancialForm
// @Accept json
// @Produce json
// @Param clinicId path string true "Clinic ID"
// @Success 200 {object} domain.H
// @Failure 400 {object} domain.H
// @Failure 500 {object} domain.H
// @Router /financial-form/clinic/{clinicId} [get]
func (h *FinancialFormHandler) GetFinancialFormsByClinic(c *gin.Context) {
	clinicIDStr := c.Param("clinicId")
	clinicID, err := uuid.Parse(clinicIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid clinic ID"})
		return
	}

	forms, err := h.financialFormService.GetFinancialFormsByClinicID(c.Request.Context(), clinicID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	utils.JSONResponse(c, http.StatusOK, "financial forms retrieved successfully", forms, nil)
}

// UpdateFinancialForm updates a financial form
// PUT /api/v1/financial-form/:id
// @Summary Update financial form
// @Description Update a financial form by ID
// @Tags FinancialForm
// @Accept json
// @Produce json
// @Param id path string true "Financial Form ID"
// @Param form body domain.FinancialForm true "Financial form information"
// @Success 200 {object} domain.H
// @Failure 400 {object} domain.H
// @Failure 404 {object} domain.H
// @Failure 500 {object} domain.H
// @Router /financial-form/{id} [put]
func (h *FinancialFormHandler) UpdateFinancialForm(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid financial form ID"})
		return
	}

	var form domain.FinancialForm
	if err := c.ShouldBindJSON(&form); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	form.ID = id
	err = h.financialFormService.UpdateFinancialForm(c.Request.Context(), &form)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	utils.JSONResponse(c, http.StatusOK, "financial form updated successfully", form, nil)
}

// DeleteFinancialForm deletes a financial form
// DELETE /api/v1/financial-form/:id
// @Summary Delete financial form
// @Description Delete a financial form by ID
// @Tags FinancialForm
// @Accept json
// @Produce json
// @Param id path string true "Financial Form ID"
// @Success 200 {object} domain.H
// @Failure 400 {object} domain.H
// @Failure 404 {object} domain.H
// @Failure 500 {object} domain.H
// @Router /financial-form/{id} [delete]
func (h *FinancialFormHandler) DeleteFinancialForm(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid financial form ID"})
		return
	}

	err = h.financialFormService.DeleteFinancialForm(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	utils.JSONResponse(c, http.StatusOK, "financial form deleted successfully", nil, nil)
}
