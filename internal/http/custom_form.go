package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/application/usecase"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
	utils "github.com/iamarpitzala/aca-reca-backend/util"
)

// CustomFormHandler handles custom forms, entries, and journal posting.
// Enforces clinic-level access control (RBAC) for all operations.
type CustomFormHandler struct {
	formUC        *usecase.CustomFormService
	postingUC     *usecase.TransactionPostingService
	userClinicUC  *usecase.UserClinicService
}

func NewCustomFormHandler(formUC *usecase.CustomFormService, postingUC *usecase.TransactionPostingService, userClinicUC *usecase.UserClinicService) *CustomFormHandler {
	return &CustomFormHandler{formUC: formUC, postingUC: postingUC, userClinicUC: userClinicUC}
}

func (h *CustomFormHandler) getAuthUserID(c *gin.Context) (uuid.UUID, bool) {
	v, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
		return uuid.Nil, false
	}
	userID, ok := v.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user context"})
		return uuid.Nil, false
	}
	return userID, true
}

// checkClinicAccess verifies the authenticated user has access to the clinic.
func (h *CustomFormHandler) checkClinicAccess(c *gin.Context, clinicID uuid.UUID) bool {
	userID, ok := h.getAuthUserID(c)
	if !ok {
		return false
	}
	hasAccess, err := h.userClinicUC.UserHasAccessToClinic(c.Request.Context(), userID, clinicID)
	if err != nil || !hasAccess {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied: you do not have access to this clinic"})
		return false
	}
	return true
}

func (h *CustomFormHandler) Create(c *gin.Context) {
	userID, ok := h.getAuthUserID(c)
	if !ok {
		return
	}
	var req domain.CreateCustomFormRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	clinicID, err := uuid.Parse(req.ClinicID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid clinic ID"})
		return
	}
	if !h.checkClinicAccess(c, clinicID) {
		return
	}
	resp, err := h.formUC.Create(c.Request.Context(), &req, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	utils.JSONResponse(c, http.StatusCreated, "custom form created", resp, nil)
}

func (h *CustomFormHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid form ID"})
		return
	}
	resp, err := h.formUC.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	clinicID, _ := uuid.Parse(resp.ClinicID)
	if !h.checkClinicAccess(c, clinicID) {
		return
	}
	utils.JSONResponse(c, http.StatusOK, "custom form retrieved", resp, nil)
}

func (h *CustomFormHandler) GetByClinicID(c *gin.Context) {
	clinicID, err := uuid.Parse(c.Param("clinicId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid clinic ID"})
		return
	}
	if !h.checkClinicAccess(c, clinicID) {
		return
	}
	list, err := h.formUC.GetByClinicID(c.Request.Context(), clinicID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	utils.JSONResponse(c, http.StatusOK, "custom forms retrieved", list, nil)
}

func (h *CustomFormHandler) GetPublishedByClinicID(c *gin.Context) {
	clinicID, err := uuid.Parse(c.Param("clinicId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid clinic ID"})
		return
	}
	if !h.checkClinicAccess(c, clinicID) {
		return
	}
	list, err := h.formUC.GetPublishedByClinicID(c.Request.Context(), clinicID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	utils.JSONResponse(c, http.StatusOK, "published custom forms retrieved", list, nil)
}

func (h *CustomFormHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid form ID"})
		return
	}
	form, err := h.formUC.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	clinicID, _ := uuid.Parse(form.ClinicID)
	if !h.checkClinicAccess(c, clinicID) {
		return
	}
	var req domain.UpdateCustomFormRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := h.formUC.UpdateByRequest(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	utils.JSONResponse(c, http.StatusOK, "custom form updated", resp, nil)
}

func (h *CustomFormHandler) Publish(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid form ID"})
		return
	}
	form, err := h.formUC.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	clinicID, _ := uuid.Parse(form.ClinicID)
	if !h.checkClinicAccess(c, clinicID) {
		return
	}
	resp, err := h.formUC.Publish(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	utils.JSONResponse(c, http.StatusOK, "form published", resp, nil)
}

func (h *CustomFormHandler) Unpublish(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid form ID"})
		return
	}
	form, err := h.formUC.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	clinicID, _ := uuid.Parse(form.ClinicID)
	if !h.checkClinicAccess(c, clinicID) {
		return
	}
	resp, err := h.formUC.Unpublish(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	utils.JSONResponse(c, http.StatusOK, "form unpublished", resp, nil)
}

func (h *CustomFormHandler) Archive(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid form ID"})
		return
	}
	form, err := h.formUC.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	clinicID, _ := uuid.Parse(form.ClinicID)
	if !h.checkClinicAccess(c, clinicID) {
		return
	}
	resp, err := h.formUC.Archive(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	utils.JSONResponse(c, http.StatusOK, "form archived", resp, nil)
}

func (h *CustomFormHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid form ID"})
		return
	}
	form, err := h.formUC.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	clinicID, _ := uuid.Parse(form.ClinicID)
	if !h.checkClinicAccess(c, clinicID) {
		return
	}
	if err := h.formUC.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	utils.JSONResponse(c, http.StatusOK, "custom form deleted", nil, nil)
}

func (h *CustomFormHandler) Duplicate(c *gin.Context) {
	userID, ok := h.getAuthUserID(c)
	if !ok {
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid form ID"})
		return
	}
	form, err := h.formUC.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	clinicID, _ := uuid.Parse(form.ClinicID)
	if !h.checkClinicAccess(c, clinicID) {
		return
	}
	resp, err := h.formUC.Duplicate(c.Request.Context(), id, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	utils.JSONResponse(c, http.StatusCreated, "form duplicated", resp, nil)
}

// Entry handlers

func (h *CustomFormHandler) CreateEntry(c *gin.Context) {
	userID, ok := h.getAuthUserID(c)
	if !ok {
		return
	}
	var req domain.CreateEntryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	clinicID, err := uuid.Parse(req.ClinicID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid clinic ID"})
		return
	}
	if !h.checkClinicAccess(c, clinicID) {
		return
	}
	resp, err := h.formUC.CreateEntryFromRequest(c.Request.Context(), &req, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	utils.JSONResponse(c, http.StatusCreated, "entry created", resp, nil)
}

func (h *CustomFormHandler) GetEntryByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("entryId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid entry ID"})
		return
	}
	entry, err := h.formUC.GetEntryByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if !h.checkClinicAccess(c, entry.ClinicID) {
		return
	}
	resp, err := h.formUC.GetEntryResponseByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	utils.JSONResponse(c, http.StatusOK, "entry retrieved", resp, nil)
}

func (h *CustomFormHandler) GetEntriesByFormID(c *gin.Context) {
	formID, err := uuid.Parse(c.Param("formId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid form ID"})
		return
	}
	form, err := h.formUC.GetByID(c.Request.Context(), formID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	clinicID, _ := uuid.Parse(form.ClinicID)
	if !h.checkClinicAccess(c, clinicID) {
		return
	}
	list, err := h.formUC.GetEntriesResponseByFormID(c.Request.Context(), formID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	utils.JSONResponse(c, http.StatusOK, "entries retrieved", list, nil)
}

func (h *CustomFormHandler) GetEntriesByClinicID(c *gin.Context) {
	clinicID, err := uuid.Parse(c.Param("clinicId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid clinic ID"})
		return
	}
	if !h.checkClinicAccess(c, clinicID) {
		return
	}
	list, err := h.formUC.GetEntriesResponseByClinicID(c.Request.Context(), clinicID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	utils.JSONResponse(c, http.StatusOK, "entries retrieved", list, nil)
}

func (h *CustomFormHandler) UpdateEntry(c *gin.Context) {
	id, err := uuid.Parse(c.Param("entryId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid entry ID"})
		return
	}
	entry, err := h.formUC.GetEntryByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if !h.checkClinicAccess(c, entry.ClinicID) {
		return
	}
	var req domain.UpdateEntryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := h.formUC.UpdateEntryFromRequest(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	utils.JSONResponse(c, http.StatusOK, "entry updated", resp, nil)
}

func (h *CustomFormHandler) DeleteEntry(c *gin.Context) {
	id, err := uuid.Parse(c.Param("entryId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid entry ID"})
		return
	}
	entry, err := h.formUC.GetEntryByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if !h.checkClinicAccess(c, entry.ClinicID) {
		return
	}
	if err := h.formUC.DeleteEntry(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	utils.JSONResponse(c, http.StatusOK, "entry deleted", nil, nil)
}

func (h *CustomFormHandler) PreviewCalculations(c *gin.Context) {
	var req domain.PreviewCalculationsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	formID, err := uuid.Parse(req.FormID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid form ID"})
		return
	}
	form, err := h.formUC.GetByID(c.Request.Context(), formID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	clinicID, _ := uuid.Parse(form.ClinicID)
	if !h.checkClinicAccess(c, clinicID) {
		return
	}
	deductions := req.Deductions
	if len(deductions) == 0 {
		deductions = nil
	}
	calculations, err := h.formUC.PreviewCalculations(c.Request.Context(), formID, req.Values, deductions)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	var out map[string]interface{}
	if err := json.Unmarshal(calculations, &out); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid calculations"})
		return
	}
	utils.JSONResponse(c, http.StatusOK, "calculations", out, nil)
}

// Journal entry handlers (post entry to ledger)

func (h *CustomFormHandler) GenerateEntryTransactions(c *gin.Context) {
	entryID, err := uuid.Parse(c.Param("entryId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid entry ID"})
		return
	}
	entry, err := h.formUC.GetEntryByID(c.Request.Context(), entryID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if !h.checkClinicAccess(c, entry.ClinicID) {
		return
	}
	list, err := h.postingUC.PostEntryToLedger(c.Request.Context(), entryID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	utils.JSONResponse(c, http.StatusCreated, "transactions generated", list, nil)
}

func (h *CustomFormHandler) GetEntryTransactions(c *gin.Context) {
	entryID, err := uuid.Parse(c.Param("entryId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid entry ID"})
		return
	}
	entry, err := h.formUC.GetEntryByID(c.Request.Context(), entryID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if !h.checkClinicAccess(c, entry.ClinicID) {
		return
	}
	list, err := h.postingUC.ListJournalEntriesByEntry(c.Request.Context(), entryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	utils.JSONResponse(c, http.StatusOK, "transactions", list, nil)
}

func (h *CustomFormHandler) GetClinicTransactions(c *gin.Context) {
	clinicID, err := uuid.Parse(c.Param("clinicId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid clinic ID"})
		return
	}
	if !h.checkClinicAccess(c, clinicID) {
		return
	}
	f := &domain.ListTransactionsFilters{
		Search:        c.Query("search"),
		TaxCategory:   c.Query("taxCategory"),
		Status:        c.Query("status"),
		DateFrom:      c.Query("dateFrom"),
		DateTo:        c.Query("dateTo"),
		SortField:     c.DefaultQuery("sortField", "date"),
		SortDirection: c.DefaultQuery("sortDirection", "desc"),
		Page:          1,
		Limit:         50,
	}
	if p := c.Query("page"); p != "" {
		if v, err := strconv.Atoi(p); err == nil && v > 0 {
			f.Page = v
		}
	}
	if l := c.Query("limit"); l != "" {
		if v, err := strconv.Atoi(l); err == nil && v > 0 && v <= 100 {
			f.Limit = v
		}
	}
	resp, err := h.postingUC.ListJournalEntries(c.Request.Context(), clinicID, f)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	utils.JSONResponse(c, http.StatusOK, "transactions", resp, nil)
}

func (h *CustomFormHandler) GetFormFieldCOAMapping(c *gin.Context) {
	formID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid form ID"})
		return
	}
	clinicID, err := uuid.Parse(c.Param("clinicId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid clinic ID"})
		return
	}
	if !h.checkClinicAccess(c, clinicID) {
		return
	}
	resp, err := h.postingUC.GetFormFieldCOAMapping(c.Request.Context(), formID, clinicID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	utils.JSONResponse(c, http.StatusOK, "form field COA mapping", resp, nil)
}
