package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iamarpitzala/aca-reca-backend/internal/application/port"
	"github.com/iamarpitzala/aca-reca-backend/internal/application/usecase"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain/form"
	utils "github.com/iamarpitzala/aca-reca-backend/util"
)

// FieldEntryHandler handles field entry HTTP endpoints.
// Enforces clinic-level access control (RBAC) for all operations.
type FieldEntryHandler struct {
	entryUC      *usecase.FieldEntryService
	userClinicUC *usecase.UserClinicService
	formRepo     port.CustomFormRepository
}

func NewFieldEntryHandler(entryUC *usecase.FieldEntryService, userClinicUC *usecase.UserClinicService, formRepo port.CustomFormRepository) *FieldEntryHandler {
	return &FieldEntryHandler{
		entryUC:      entryUC,
		userClinicUC: userClinicUC,
		formRepo:     formRepo,
	}
}

// Create handles POST /entry
func (h *FieldEntryHandler) Create(c *gin.Context) {
	userID, ok := GetAuthUserID(c)
	if !ok {
		return
	}

	var req form.EntryFieldRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get form to resolve clinic for access control
	fm, err := h.formRepo.GetByID(c.Request.Context(), req.FormID)
	if err != nil || fm == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "form not found"})
		return
	}

	if !RequireClinicAccess(c, h.userClinicUC, fm.ClinicID) {
		return
	}

	resp, err := h.entryUC.Create(c.Request.Context(), &req, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	utils.JSONResponse(c, http.StatusCreated, utils.MsgEntryCreated, resp, nil)
}

// GetByID handles GET /field-entries/:id
func (h *FieldEntryHandler) GetByID(c *gin.Context) {
	id := c.Param("id")

	resp, err := h.entryUC.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if !RequireClinicAccess(c, h.userClinicUC, resp.ClinicID) {
		return
	}

	utils.JSONResponse(c, http.StatusOK, utils.MsgFieldEntryRetrieved, resp, nil)
}

// GetByFormID handles GET /field-entries/form/:formId
func (h *FieldEntryHandler) GetByFormID(c *gin.Context) {
	formID := c.Param("formId")

	fm, err := h.formRepo.GetByID(c.Request.Context(), formID)
	if err != nil || fm == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "form not found"})
		return
	}

	if !RequireClinicAccess(c, h.userClinicUC, fm.ClinicID) {
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
	clinicID := c.Param("clinicId")

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
	id := c.Param("id")

	existing, err := h.entryUC.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if !RequireClinicAccess(c, h.userClinicUC, existing.ClinicID) {
		return
	}

	var req form.EntryFieldUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req.ID = id

	resp, err := h.entryUC.Update(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	utils.JSONResponse(c, http.StatusOK, utils.MsgFieldEntryUpdated, resp, nil)
}

// Delete handles DELETE /field-entries/:id
func (h *FieldEntryHandler) Delete(c *gin.Context) {
	id := c.Param("id")

	existing, err := h.entryUC.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if !RequireClinicAccess(c, h.userClinicUC, existing.ClinicID) {
		return
	}

	if err := h.entryUC.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	utils.JSONResponse(c, http.StatusOK, utils.MsgFieldEntryDeleted, nil, nil)
}
