package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/application/usecase"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
)

type EntryHandler struct {
	service *usecase.EntryService
}

func NewEntryHandler(entryService *usecase.EntryService) *EntryHandler {
	return &EntryHandler{
		service: entryService,
	}
}

func (h *EntryHandler) AddEntry(c *gin.Context) {
	formId := c.Query("form_Id")
	formUUID, err := uuid.Parse(formId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid form ID"})
		return
	}
	var entry domain.CommonEntry
	if err := c.ShouldBindJSON(&entry); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	field, err := h.service.AddEntry(c, formUUID, entry)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"success": field})
}
