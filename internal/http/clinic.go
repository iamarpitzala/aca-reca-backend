package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
	"github.com/iamarpitzala/aca-reca-backend/internal/service"
)

type ClinicHandler struct {
	clinicService *service.ClinicService
}

func NewClinicHandler(clinicService *service.ClinicService) *ClinicHandler {
	return &ClinicHandler{
		clinicService: clinicService,
	}
}

// CreateClinic creates a new clinic
// POST /api/v1/clinic
// @Summary Create a new clinic
// @Description Create a new clinic with the given information
// @Tags Clinic
// @Accept json
// @Produce json
// @Param clinic body domain.Clinic true "Clinic information"
// @Success 201 {object} domain.H
// @Failure 400 {object} domain.H
// @Failure 500 {object} domain.H
// @Router /clinic [post]
func (h *ClinicHandler) CreateClinic(c *gin.Context) {
	var clinic domain.Clinic
	if err := c.ShouldBindJSON(&clinic); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := h.clinicService.CreateClinic(c.Request.Context(), &clinic)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "clinic created successfully", "clinic_id": clinic.ID})
}

// GetClinic retrieves a clinic by ID
// GET /api/v1/clinic/:id
// @Summary Retrieve a clinic by ID
// @Description Retrieve a clinic by ID
// @Tags Clinic
// @Accept json
// @Produce json
// @Param id path string true "Clinic ID"
// @Success 200 {object} domain.H
// @Failure 400 {object} domain.H
// @Failure 404 {object} domain.H
// @Failure 500 {object} domain.H
// @Router /clinic/{id} [get]
func (h *ClinicHandler) GetClinic(c *gin.Context) {
	id := c.Param("id")
	idUUID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid clinic ID"})
		return
	}
	clinic, err := h.clinicService.GetClinicByID(c.Request.Context(), idUUID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "clinic retrieved successfully", "clinic_id": clinic.ID})
}

// UpdateClinic updates a clinic by ID
// PUT /api/v1/clinic/:id
// @Summary Update a clinic by ID
// @Description Update a clinic by ID with the given information
// @Tags Clinic
// @Accept json
// @Produce json
// @Param id path string true "Clinic ID"
// @Param clinic body domain.Clinic true "Clinic information"
// @Success 200 {object} domain.H
// @Failure 400 {object} domain.H
// @Failure 404 {object} domain.H
// @Failure 500 {object} domain.H
// @Router /clinic/{id} [put]
func (h *ClinicHandler) UpdateClinic(c *gin.Context) {
	id := c.Param("id")
	idUUID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid clinic ID"})
		return
	}
	var clinic domain.Clinic
	if err := c.ShouldBindJSON(&clinic); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	clinic.ID = idUUID
	err = h.clinicService.UpdateClinic(c.Request.Context(), &clinic)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "clinic updated successfully", "clinic_id": clinic.ID})
}

// DeleteClinic deletes a clinic by ID
// DELETE /api/v1/clinic/:id
// @Summary Delete a clinic by ID
// @Description Delete a clinic by ID
// @Tags Clinic
// @Accept json
// @Produce json
// @Param id path string true "Clinic ID"
// @Success 200 {object} domain.H
// @Failure 400 {object} domain.H
// @Failure 404 {object} domain.H
// @Failure 500 {object} domain.H
// @Router /clinic/{id} [delete]
func (h *ClinicHandler) DeleteClinic(c *gin.Context) {
	id := c.Param("id")
	idUUID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid clinic ID"})
		return
	}
	err = h.clinicService.DeleteClinic(c.Request.Context(), idUUID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "clinic deleted successfully", "clinic_id": idUUID})
}

// GetAllClinics retrieves all clinics
// GET /api/v1/clinic
// @Summary Retrieve all clinics
// @Description Retrieve all clinics
// @Tags Clinic
// @Accept json
// @Produce json
// @Success 200 {object} domain.H
// @Failure 404 {object} domain.H
// @Failure 500 {object} domain.H
// @Router /clinic [get]
func (h *ClinicHandler) GetAllClinics(c *gin.Context) {
	clinics, err := h.clinicService.GetAllClinics(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "clinics retrieved successfully", "clinics": clinics})
}

// GetClinicByABNNumber retrieves a clinic by ABN number
// GET /api/v1/clinic/abn/:abnNumber
// @Summary Retrieve a clinic by ABN number
// @Description Retrieve a clinic by ABN number
// @Tags Clinic
// @Accept json
// @Produce json
// @Param abnNumber path string true "ABN number"
// @Success 200 {object} domain.H
// @Failure 400 {object} domain.H
// @Failure 404 {object} domain.H
// @Failure 500 {object} domain.H
// @Router /clinic/abn/{abnNumber} [get]
func (h *ClinicHandler) GetClinicByABNNumber(c *gin.Context) {
	abnNumber := c.Param("abnNumber")
	clinic, err := h.clinicService.GetClinicByABNNumber(c.Request.Context(), abnNumber)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "clinic retrieved successfully", "clinic": clinic})
}
