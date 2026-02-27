package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iamarpitzala/aca-reca-backend/internal/application/usecase"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain/form"
	utils "github.com/iamarpitzala/aca-reca-backend/util"
)

// CustomFormEntryHandler handles custom form entry HTTP endpoints (tbl_custom_form_entry + tbl_custom_form_entry_value).
type CustomFormEntryHandler struct {
	entryUC      *usecase.CustomFormEntryService
	userClinicUC *usecase.UserClinicService
}

func NewCustomFormEntryHandler(entryUC *usecase.CustomFormEntryService, userClinicUC *usecase.UserClinicService) *CustomFormEntryHandler {
	return &CustomFormEntryHandler{
		entryUC:      entryUC,
		userClinicUC: userClinicUC,
	}
}

// Create handles POST /entries — create an entry with values for a form version.
func (h *CustomFormEntryHandler) Create(c *gin.Context) {
	userID, ok := GetAuthUserID(c)
	if !ok {
		return
	}
	var req form.EntryCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Resolve clinic from form version for access control
	formVersionID := req.FormVersionID
	if formVersionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "formVersionId required"})
		return
	}
	// We validate form version and form in use case; clinic check requires form. Use case returns error if form not found.
	// For access we need clinic ID: get form version -> form -> clinic. Handler can call a use-case method that returns clinic ID for a form version.
	// For simplicity: require clinicId in query or body and RequireClinicAccess. Or add GetClinicIDForFormVersion to entry use case.
	// Entry use case has formRepo and versionRepo - we can add GetClinicIDForFormVersion(versionID) and use it here.
	// For now: create entry; use case will validate form version and form. Then for access we need to ensure user has clinic access. So after create we have entry; we don't have clinic on entry. So we need GetClinicIDForFormVersion in use case. Let me add it.
	clinicID, err := h.entryUC.GetClinicIDForFormVersion(c.Request.Context(), req.FormVersionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if !RequireClinicAccess(c, h.userClinicUC, clinicID) {
		return
	}
	resp, err := h.entryUC.Create(c.Request.Context(), &req, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	utils.JSONResponse(c, http.StatusCreated, utils.MsgEntryCreated, resp, nil)
}

// GetByID handles GET /entries/:entryId
func (h *CustomFormEntryHandler) GetByID(c *gin.Context) {
	id := c.Param("entryId")
	resp, err := h.entryUC.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	clinicID, err := h.entryUC.GetClinicIDForEntry(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if !RequireClinicAccess(c, h.userClinicUC, clinicID) {
		return
	}
	utils.JSONResponse(c, http.StatusOK, utils.MsgFieldEntryRetrieved, resp, nil)
}

// GetByFormID handles GET /entries/form/:formId
func (h *CustomFormEntryHandler) GetByFormID(c *gin.Context) {
	formID := c.Param("formId")
	f, err := h.entryUC.GetFormByIDForAccess(c.Request.Context(), formID)
	if err != nil || f == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "form not found"})
		return
	}
	if !RequireClinicAccess(c, h.userClinicUC, f.ClinicID) {
		return
	}
	list, err := h.entryUC.GetByFormID(c.Request.Context(), formID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	utils.JSONResponse(c, http.StatusOK, utils.MsgFieldEntriesRetrieved, list, nil)
}

// GetByClinicID handles GET /entries/clinic/:clinicId
func (h *CustomFormEntryHandler) GetByClinicID(c *gin.Context) {
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

// Update handles PUT /entries/:entryId
func (h *CustomFormEntryHandler) Update(c *gin.Context) {
	entryID := c.Param("entryId")
	clinicID, err := h.entryUC.GetClinicIDForEntry(c.Request.Context(), entryID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if !RequireClinicAccess(c, h.userClinicUC, clinicID) {
		return
	}
	var req form.EntryUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := h.entryUC.Update(c.Request.Context(), entryID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	utils.JSONResponse(c, http.StatusOK, utils.MsgFieldEntryUpdated, resp, nil)
}

// Delete handles DELETE /entries/:entryId
func (h *CustomFormEntryHandler) Delete(c *gin.Context) {
	id := c.Param("entryId")
	clinicID, err := h.entryUC.GetClinicIDForEntry(c.Request.Context(), id)
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
