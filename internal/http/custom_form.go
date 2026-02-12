package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
	"github.com/iamarpitzala/aca-reca-backend/internal/service"
	utils "github.com/iamarpitzala/aca-reca-backend/util"
)

type CustomFormHandler struct {
	svc            *service.CustomFormService
	transactionSvc *service.TransactionService
}

func NewCustomFormHandler(svc *service.CustomFormService, transactionSvc *service.TransactionService) *CustomFormHandler {
	return &CustomFormHandler{svc: svc, transactionSvc: transactionSvc}
}

func (h *CustomFormHandler) getUserID(c *gin.Context) (uuid.UUID, bool) {
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

func (h *CustomFormHandler) Create(c *gin.Context) {
	userID, ok := h.getUserID(c)
	if !ok {
		return
	}
	var req domain.CreateCustomFormRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := h.svc.Create(c.Request.Context(), &req, userID)
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
	resp, err := h.svc.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
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
	list, err := h.svc.GetByClinicID(c.Request.Context(), clinicID)
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
	list, err := h.svc.GetPublishedByClinicID(c.Request.Context(), clinicID)
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
	var req domain.UpdateCustomFormRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := h.svc.Update(c.Request.Context(), id, &req)
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
	resp, err := h.svc.Publish(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	utils.JSONResponse(c, http.StatusOK, "form published", resp, nil)
}

func (h *CustomFormHandler) Archive(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid form ID"})
		return
	}
	resp, err := h.svc.Archive(c.Request.Context(), id)
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
	if err := h.svc.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	utils.JSONResponse(c, http.StatusOK, "custom form deleted", nil, nil)
}

func (h *CustomFormHandler) Duplicate(c *gin.Context) {
	userID, ok := h.getUserID(c)
	if !ok {
		return
	}
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid form ID"})
		return
	}
	resp, err := h.svc.Duplicate(c.Request.Context(), id, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	utils.JSONResponse(c, http.StatusCreated, "form duplicated", resp, nil)
}

// Entry handlers

func (h *CustomFormHandler) CreateEntry(c *gin.Context) {
	userID, ok := h.getUserID(c)
	if !ok {
		return
	}
	var req domain.CreateEntryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := h.svc.CreateEntry(c.Request.Context(), &req, userID)
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
	resp, err := h.svc.GetEntryByID(c.Request.Context(), id)
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
	list, err := h.svc.GetEntriesByFormID(c.Request.Context(), formID)
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
	list, err := h.svc.GetEntriesByClinicID(c.Request.Context(), clinicID)
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
	var req domain.UpdateEntryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := h.svc.UpdateEntry(c.Request.Context(), id, &req)
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
	if err := h.svc.DeleteEntry(c.Request.Context(), id); err != nil {
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
	deductions := req.Deductions
	if len(deductions) == 0 {
		deductions = nil
	}
	calculations, err := h.svc.PreviewCalculations(c.Request.Context(), formID, req.Values, deductions)
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

// Transaction handlers (form field entries -> COA mapping)

func (h *CustomFormHandler) GenerateEntryTransactions(c *gin.Context) {
	entryID, err := uuid.Parse(c.Param("entryId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid entry ID"})
		return
	}
	list, err := h.transactionSvc.GenerateFromEntry(c.Request.Context(), entryID)
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
	list, err := h.transactionSvc.ListByEntryID(c.Request.Context(), entryID)
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
	resp, err := h.transactionSvc.List(c.Request.Context(), clinicID, f)
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
	resp, err := h.transactionSvc.GetFormFieldCOAMapping(c.Request.Context(), formID, clinicID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	utils.JSONResponse(c, http.StatusOK, "form field COA mapping", resp, nil)
}
