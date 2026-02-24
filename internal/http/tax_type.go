package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iamarpitzala/aca-reca-backend/internal/application/usecase"
	utils "github.com/iamarpitzala/aca-reca-backend/util"
)

type TaxTypeHandler struct {
	taxTypeUC *usecase.TaxTypeService
}

func NewTaxTypeHandler(taxTypeUC *usecase.TaxTypeService) *TaxTypeHandler {
	return &TaxTypeHandler{taxTypeUC: taxTypeUC}
}

// GetTaxTypes returns all tax types (tbl_tax_type) for form field lookups.
// GET /api/v1/tax-types
func (h *TaxTypeHandler) GetTaxTypes(c *gin.Context) {
	response, err := h.taxTypeUC.GetTaxTypes(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	utils.JSONResponse(c, http.StatusOK, "tax types retrieved successfully", response, nil)
}
