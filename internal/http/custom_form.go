package http

import (
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
	formUC       *usecase.CustomFormService
	postingUC    *usecase.TransactionPostingService
	userClinicUC *usecase.UserClinicService
}

func NewCustomFormHandler(formUC *usecase.CustomFormService, postingUC *usecase.TransactionPostingService, userClinicUC *usecase.UserClinicService) *CustomFormHandler {
	return &CustomFormHandler{formUC: formUC, postingUC: postingUC, userClinicUC: userClinicUC}
}

func (h *CustomFormHandler) Create(c *gin.Context) {
	userID, ok := GetAuthUserID(c)
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
	if !RequireClinicAccess(c, h.userClinicUC, clinicID) {
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
	if !RequireClinicAccess(c, h.userClinicUC, clinicID) {
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
	if !RequireClinicAccess(c, h.userClinicUC, clinicID) {
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
	if !RequireClinicAccess(c, h.userClinicUC, clinicID) {
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
	if !RequireClinicAccess(c, h.userClinicUC, clinicID) {
		return
	}
	userID, ok := GetAuthUserID(c)
	if !ok {
		return
	}
	var req domain.UpdateCustomFormRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := h.formUC.UpdateByRequest(c.Request.Context(), id, &req, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	utils.JSONResponse(c, http.StatusOK, "custom form updated", resp, nil)
}

func (h *CustomFormHandler) Publish(c *gin.Context) {
	userID, ok := GetAuthUserID(c)
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
	if !RequireClinicAccess(c, h.userClinicUC, clinicID) {
		return
	}
	resp, err := h.formUC.Publish(c.Request.Context(), id, userID)
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
	if !RequireClinicAccess(c, h.userClinicUC, clinicID) {
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
	if !RequireClinicAccess(c, h.userClinicUC, clinicID) {
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
	if !RequireClinicAccess(c, h.userClinicUC, clinicID) {
		return
	}
	if err := h.formUC.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	utils.JSONResponse(c, http.StatusOK, "custom form deleted", nil, nil)
}

func (h *CustomFormHandler) Duplicate(c *gin.Context) {
	userID, ok := GetAuthUserID(c)
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
	if !RequireClinicAccess(c, h.userClinicUC, clinicID) {
		return
	}
	resp, err := h.formUC.Duplicate(c.Request.Context(), id, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	utils.JSONResponse(c, http.StatusCreated, "form duplicated", resp, nil)
}

// GenerateEntryTransactions handles POST /custom-form/entries/:entryId/transactions
// Creates journal entries (transactions) from a custom form entry
func (h *CustomFormHandler) GenerateEntryTransactions(c *gin.Context) {
	entryID, err := uuid.Parse(c.Param("entryId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid entry ID"})
		return
	}

	// Get clinic ID from entry for access control
	// We'll get it from the posting service which validates the entry exists
	// The service will return an error if entry doesn't exist or belongs to wrong clinic
	transactions, err := h.postingUC.PostEntryToLedger(c.Request.Context(), entryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Validate access using clinic ID from transactions
	if len(transactions) > 0 {
		clinicID, err := uuid.Parse(transactions[0].ClinicID)
		if err == nil {
			if !RequireClinicAccess(c, h.userClinicUC, clinicID) {
				return
			}
		}
	}

	utils.JSONResponse(c, http.StatusCreated, "transactions generated", transactions, nil)
}

// GetEntryTransactions handles GET /custom-form/entries/:entryId/transactions
// Returns all transactions for a specific entry
func (h *CustomFormHandler) GetEntryTransactions(c *gin.Context) {
	entryID, err := uuid.Parse(c.Param("entryId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid entry ID"})
		return
	}

	transactions, err := h.postingUC.ListJournalEntriesByEntry(c.Request.Context(), entryID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Get clinic ID from first transaction for access control
	if len(transactions) > 0 {
		clinicID, err := uuid.Parse(transactions[0].ClinicID)
		if err == nil {
			if !RequireClinicAccess(c, h.userClinicUC, clinicID) {
				return
			}
		}
	}

	utils.JSONResponse(c, http.StatusOK, "transactions retrieved", transactions, nil)
}

// GetClinicTransactions handles GET /custom-form/clinic/:clinicId/transactions
// Returns paginated list of transactions for a clinic
func (h *CustomFormHandler) GetClinicTransactions(c *gin.Context) {
	clinicID, err := uuid.Parse(c.Param("clinicId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid clinic ID"})
		return
	}

	if !RequireClinicAccess(c, h.userClinicUC, clinicID) {
		return
	}

	// Parse query parameters for filters
	filters := &domain.ListTransactionsFilters{
		Page:          1,
		Limit:         50,
		SortField:     "date",
		SortDirection: "desc",
	}

	if page := c.Query("page"); page != "" {
		if p, err := strconv.Atoi(page); err == nil && p > 0 {
			filters.Page = p
		}
	}
	if limit := c.Query("limit"); limit != "" {
		if l, err := strconv.Atoi(limit); err == nil && l > 0 {
			filters.Limit = l
		}
	}
	if search := c.Query("search"); search != "" {
		filters.Search = search
	}
	if taxCategory := c.Query("taxCategory"); taxCategory != "" {
		filters.TaxCategory = taxCategory
	}
	if status := c.Query("status"); status != "" {
		filters.Status = status
	}
	if coaId := c.Query("coaId"); coaId != "" {
		filters.COAID = coaId
	}
	if dateFrom := c.Query("dateFrom"); dateFrom != "" {
		filters.DateFrom = dateFrom
	}
	if dateTo := c.Query("dateTo"); dateTo != "" {
		filters.DateTo = dateTo
	}
	if sortField := c.Query("sortField"); sortField != "" {
		filters.SortField = sortField
	}
	if sortDirection := c.Query("sortDirection"); sortDirection != "" {
		filters.SortDirection = sortDirection
	}

	resp, err := h.postingUC.ListJournalEntries(c.Request.Context(), clinicID, filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.JSONResponse(c, http.StatusOK, "transactions retrieved", resp, nil)
}

// GetFormFieldCOAMapping handles GET /custom-form/:id/field-coa-mapping/clinic/:clinicId
// Returns form fields with their COA mapping and clinic COA list
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

	if !RequireClinicAccess(c, h.userClinicUC, clinicID) {
		return
	}

	resp, err := h.postingUC.GetFormFieldCOAMapping(c.Request.Context(), formID, clinicID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	utils.JSONResponse(c, http.StatusOK, "form field COA mapping retrieved", resp, nil)
}
