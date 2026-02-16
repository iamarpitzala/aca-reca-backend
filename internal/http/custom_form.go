package http

import (
	"net/http"

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
