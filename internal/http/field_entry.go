package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/application/usecase"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
	utils "github.com/iamarpitzala/aca-reca-backend/util"
)

// FieldEntryHandler handles field entry HTTP endpoints.
// Enforces clinic-level access control (RBAC) for all operations.
type FieldEntryHandler struct {
	entryUC     *usecase.FieldEntryService
	userClinicUC *usecase.UserClinicService
}

func NewFieldEntryHandler(entryUC *usecase.FieldEntryService, userClinicUC *usecase.UserClinicService) *FieldEntryHandler {
	return &FieldEntryHandler{
		entryUC:     entryUC,
		userClinicUC: userClinicUC,
	}
}

// Create handles POST /entry
func (h *FieldEntryHandler) Create(c *gin.Context) {
	userID, ok := GetAuthUserID(c)
	if !ok {
		return
	}

	var req domain.CreateEntryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Validate form ID
	formID, err := uuid.Parse(req.FormID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.ErrInvalidFormID})
		return
	}

	// Get clinic ID from form for access control
	clinicID, err := h.entryUC.GetClinicIDFromForm(c.Request.Context(), formID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if !RequireClinicAccess(c, h.userClinicUC, clinicID) {
		return
	}

	resp, err := h.entryUC.CreateEntry(c.Request.Context(), &req, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	utils.JSONResponse(c, http.StatusCreated, utils.MsgEntryCreated, resp, nil)
}

// GetByID handles GET /field-entries/:id
func (h *FieldEntryHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.ErrInvalidEntryID})
		return
	}

	// Get clinic ID for access control
	clinicID, err := h.entryUC.GetClinicIDFromEntry(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if !RequireClinicAccess(c, h.userClinicUC, clinicID) {
		return
	}

	resp, err := h.entryUC.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	utils.JSONResponse(c, http.StatusOK, utils.MsgFieldEntryRetrieved, resp, nil)
}

// GetNetDetails handles GET /entry/:id/net-details
func (h *FieldEntryHandler) GetNetDetails(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.ErrInvalidEntryID})
		return
	}

	clinicID, err := h.entryUC.GetClinicIDFromEntry(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if !RequireClinicAccess(c, h.userClinicUC, clinicID) {
		return
	}

	resp, err := h.entryUC.GetNetDetails(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": utils.ErrNetDetailsNotFound})
		return
	}

	utils.JSONResponse(c, http.StatusOK, utils.MsgNetDetailsRetrieved, resp, nil)
}

// GetGrossDetails handles GET /entry/:id/gross-details
func (h *FieldEntryHandler) GetGrossDetails(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.ErrInvalidEntryID})
		return
	}

	clinicID, err := h.entryUC.GetClinicIDFromEntry(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if !RequireClinicAccess(c, h.userClinicUC, clinicID) {
		return
	}

	resp, err := h.entryUC.GetGrossDetails(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": utils.ErrGrossDetailsNotFound})
		return
	}

	utils.JSONResponse(c, http.StatusOK, utils.MsgGrossDetailsRetrieved, resp, nil)
}

// GetByFormID handles GET /field-entries/form/:formId
func (h *FieldEntryHandler) GetByFormID(c *gin.Context) {
	formID, err := uuid.Parse(c.Param("formId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.ErrInvalidFormID})
		return
	}

	// Get clinic ID for access control
	clinicID, err := h.entryUC.GetClinicIDFromForm(c.Request.Context(), formID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if !RequireClinicAccess(c, h.userClinicUC, clinicID) {
		return
	}

	list, err := h.entryUC.GetByFormID(c.Request.Context(), formID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.JSONResponse(c, http.StatusOK, utils.MsgFieldEntriesRetrieved, list, nil)
}

// GetByClinicID handles GET /field-entries/clinic/:clinicId
func (h *FieldEntryHandler) GetByClinicID(c *gin.Context) {
	clinicID, err := uuid.Parse(c.Param("clinicId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.ErrInvalidClinicID})
		return
	}

	if !RequireClinicAccess(c, h.userClinicUC, clinicID) {
		return
	}

	list, err := h.entryUC.GetByClinicID(c.Request.Context(), clinicID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.JSONResponse(c, http.StatusOK, utils.MsgFieldEntriesRetrieved, list, nil)
}

// Update handles PUT /field-entries/:id
func (h *FieldEntryHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.ErrInvalidEntryID})
		return
	}

	// Get clinic ID for access control
	clinicID, err := h.entryUC.GetClinicIDFromEntry(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if !RequireClinicAccess(c, h.userClinicUC, clinicID) {
		return
	}

	var req struct {
		Value float64 `json:"value" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	resp, err := h.entryUC.Update(c.Request.Context(), id, req.Value)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	utils.JSONResponse(c, http.StatusOK, utils.MsgFieldEntryUpdated, resp, nil)
}

// Delete handles DELETE /field-entries/:id
func (h *FieldEntryHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.ErrInvalidEntryID})
		return
	}

	// Get clinic ID for access control
	clinicID, err := h.entryUC.GetClinicIDFromEntry(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if !RequireClinicAccess(c, h.userClinicUC, clinicID) {
		return
	}

	if err := h.entryUC.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	utils.JSONResponse(c, http.StatusOK, utils.MsgFieldEntryDeleted, nil, nil)
}
