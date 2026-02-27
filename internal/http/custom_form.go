package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iamarpitzala/aca-reca-backend/internal/application/usecase"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain/form"
	utils "github.com/iamarpitzala/aca-reca-backend/util"
)

// CustomFormHandler handles custom forms and entries.
// Enforces clinic-level access control (RBAC) for all operations.
type CustomFormHandler struct {
	formUC       *usecase.CustomFormService
	userClinicUC *usecase.UserClinicService
}

func NewCustomFormHandler(formUC *usecase.CustomFormService, userClinicUC *usecase.UserClinicService) *CustomFormHandler {
	return &CustomFormHandler{formUC: formUC, userClinicUC: userClinicUC}
}

func (h *CustomFormHandler) Create(c *gin.Context) {
	userID, ok := GetAuthUserID(c)
	if !ok {
		return
	}
	var req form.FormRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if req.Status == "" {
		req.Status = "DRAFT"
	}
	if req.CalculationMethod == "" {
		req.CalculationMethod = "NET"
	}
	clinicID := req.ClinicID
	if !RequireClinicAccess(c, h.userClinicUC, clinicID) {
		return
	}
	resp, err := h.formUC.Create(c.Request.Context(), &req, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	utils.JSONResponse(c, http.StatusCreated, utils.MsgCustomFormCreated, resp, nil)
}

func (h *CustomFormHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	resp, err := h.formUC.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	clinicID := resp.ClinicID
	if !RequireClinicAccess(c, h.userClinicUC, clinicID) {
		return
	}
	utils.JSONResponse(c, http.StatusOK, utils.MsgCustomFormRetrieved, resp, nil)
}

// GetSectionTypes returns the list of section types for the Build form (Step 1).
// Does not require clinic access; returns static config (INCOME: Collection; EXPENSES: Cost, Service and facility, Other cost).
func (h *CustomFormHandler) GetSectionTypes(c *gin.Context) {
	list := []form.SectionTypeOption{
		{ID: "collection", Category: "INCOME", Title: "Collection", Description: "Income and fees collected", Order: 1},
		{ID: "cost", Category: "EXPENSES", Title: "Cost", Description: "Direct costs and expenses", Order: 2},
		{ID: "service_and_facility", Category: "EXPENSES", Title: "Service and facility", Description: "Service & facility fee component", Order: 3},
		{ID: "other_cost", Category: "EXPENSES", Title: "Other cost", Description: "Other costs and deductions", Order: 4},
	}
	utils.JSONResponse(c, http.StatusOK, utils.MsgFormSectionTypesRetrieved, list, nil)
}

func (h *CustomFormHandler) GetByClinicID(c *gin.Context) {
	clinicID := c.Param("clinicId")
	if !RequireClinicAccess(c, h.userClinicUC, clinicID) {
		return
	}
	list, err := h.formUC.GetByClinicID(c.Request.Context(), clinicID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	utils.JSONResponse(c, http.StatusOK, utils.MsgCustomFormsRetrieved, list, nil)
}

func (h *CustomFormHandler) GetPublishedByClinicID(c *gin.Context) {
	clinicID := c.Param("clinicId")
	if !RequireClinicAccess(c, h.userClinicUC, clinicID) {
		return
	}
	list, err := h.formUC.GetPublishedByClinicID(c.Request.Context(), clinicID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	utils.JSONResponse(c, http.StatusOK, utils.MsgPublishedCustomFormsRetrieved, list, nil)
}

func (h *CustomFormHandler) Update(c *gin.Context) {
	id := c.Param("id")
	fm, err := h.formUC.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	clinicID := fm.ClinicID
	if !RequireClinicAccess(c, h.userClinicUC, clinicID) {
		return
	}
	userID, ok := GetAuthUserID(c)
	if !ok {
		return
	}
	var req form.FormRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := h.formUC.UpdateByRequest(c.Request.Context(), id, &req, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	utils.JSONResponse(c, http.StatusOK, utils.MsgCustomFormUpdated, resp, nil)
}

func (h *CustomFormHandler) Publish(c *gin.Context) {
	userID, ok := GetAuthUserID(c)
	if !ok {
		return
	}
	id := c.Param("id")
	fm, err := h.formUC.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	clinicID := fm.ClinicID
	if !RequireClinicAccess(c, h.userClinicUC, clinicID) {
		return
	}
	resp, err := h.formUC.Publish(c.Request.Context(), id, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	utils.JSONResponse(c, http.StatusOK, utils.MsgFormPublished, resp, nil)
}

func (h *CustomFormHandler) Unpublish(c *gin.Context) {
	id := c.Param("id")
	form, err := h.formUC.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	clinicID := form.ClinicID
	if !RequireClinicAccess(c, h.userClinicUC, clinicID) {
		return
	}
	resp, err := h.formUC.Unpublish(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	utils.JSONResponse(c, http.StatusOK, utils.MsgFormUnpublished, resp, nil)
}

func (h *CustomFormHandler) Archive(c *gin.Context) {
	id := c.Param("id")
	form, err := h.formUC.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	clinicID := form.ClinicID
	if !RequireClinicAccess(c, h.userClinicUC, clinicID) {
		return
	}
	resp, err := h.formUC.Archive(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	utils.JSONResponse(c, http.StatusOK, utils.MsgFormArchived, resp, nil)
}

func (h *CustomFormHandler) Delete(c *gin.Context) {
	id := c.Param("id")
	form, err := h.formUC.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	clinicID := form.ClinicID
	if !RequireClinicAccess(c, h.userClinicUC, clinicID) {
		return
	}
	if err := h.formUC.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	utils.JSONResponse(c, http.StatusOK, utils.MsgCustomFormDeleted, nil, nil)
}

func (h *CustomFormHandler) Duplicate(c *gin.Context) {
	userID, ok := GetAuthUserID(c)
	if !ok {
		return
	}
	id := c.Param("id")
	fm, err := h.formUC.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	clinicID := fm.ClinicID
	if !RequireClinicAccess(c, h.userClinicUC, clinicID) {
		return
	}
	var req form.FormRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.ClinicID = clinicID
	resp, err := h.formUC.Create(c.Request.Context(), &req, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	utils.JSONResponse(c, http.StatusCreated, utils.MsgFormDuplicated, resp, nil)
}
