package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
	"github.com/iamarpitzala/aca-reca-backend/internal/service"
)

type QuarterHandler struct {
	QuarterService *service.QuarterService
}

func NewQuarterHandler(QuarterService *service.QuarterService) *QuarterHandler {
	return &QuarterHandler{
		QuarterService: QuarterService,
	}
}

// Create creates a new financial quarter
// @Summary Create a new financial quarter
// @Description Create a new financial quarter with the given information
// @Tags Quarter
// @Accept json
// @Produce json
// @Param quarter body domain.Quarter true "Financial Quarter"
// @Success 201 {object} domain.Quarter
// @Failure 400 {object} map[string]string
// @Router /quarter [post]
func (h *QuarterHandler) Create(c *gin.Context) {
	var quarter domain.Quarter
	if err := c.ShouldBindJSON(&quarter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	err := h.QuarterService.CreateQuarter(c.Request.Context(), &quarter)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "financial quarter created successfully", "quarter": quarter})
}

// Get retrieves a financial quarter by its ID
// @Summary Get a financial quarter
// @Description Retrieve a financial quarter by its ID
// @Tags Quarter
// @Accept json
// @Produce json
// @Param id path string true "Quarter ID"
// @Success 200 {object} domain.Quarter
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /quarter/{id} [get]
func (h *QuarterHandler) Get(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid financial quarter ID"})
		return
	}

	gt, err := h.QuarterService.GetQuaterByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "financial quarter retrieved successfully", "quarter": gt})
}

// Delete deletes a financial quarter by its ID
// @Summary Delete a financial quarter
// @Description Delete a financial quarter by its ID
// @Tags Quarter
// @Accept json
// @Produce json
// @Param id path string true "Quarter ID"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /quarter/{id} [delete]
func (h *QuarterHandler) Delete(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid financial quarter ID"})
		return
	}

	err = h.QuarterService.DeleteQuarter(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "financial quarter deleted successfully"})
}

// Update updates a financial quarter by its ID
// @Summary Update a financial quarter
// @Description Update a financial quarter's details by its ID
// @Tags Quarter
// @Accept json
// @Produce json
// @Param id path string true "Quarter ID"
// @Param quarter body domain.Quarter true "Financial Quarter"
// @Success 200 {object} map[string]string
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /quarter/{id} [put]
func (h *QuarterHandler) Update(c *gin.Context) {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid financial quarter ID"})
		return
	}

	var quarter domain.Quarter
	if err := c.ShouldBindJSON(&quarter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	quarter.ID = id

	err = h.QuarterService.UpdateQuarter(c.Request.Context(), &quarter)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "financial quarter updated successfully"})
}

// List lists all financial quarters
// @Summary List financial quarters
// @Description Get a list of all financial quarters
// @Tags Quarter
// @Accept json
// @Produce json
// @Success 200 {object} []domain.Quarter
// @Failure 404 {object} map[string]string
// @Router /quarter [get]
func (h *QuarterHandler) List(c *gin.Context) {

	list, err := h.QuarterService.ListQuarter(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": list})
}
