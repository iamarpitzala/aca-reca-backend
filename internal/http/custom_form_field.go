package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iamarpitzala/aca-reca-backend/internal/application/usecase"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain/form"
	utils "github.com/iamarpitzala/aca-reca-backend/util"
)

// CustomFormFieldHandler handles custom form field, field config, and formula source HTTP endpoints.
type CustomFormFieldHandler struct {
	fieldUC       *usecase.CustomFormFieldService
	userClinicUC  *usecase.UserClinicService
	formUC        *usecase.CustomFormService
}

func NewCustomFormFieldHandler(
	fieldUC *usecase.CustomFormFieldService,
	userClinicUC *usecase.UserClinicService,
	formUC *usecase.CustomFormService,
) *CustomFormFieldHandler {
	return &CustomFormFieldHandler{
		fieldUC:      fieldUC,
		userClinicUC: userClinicUC,
		formUC:       formUC,
	}
}

// --- Form fields ---

func (h *CustomFormFieldHandler) CreateField(c *gin.Context) {
	formID := c.Param("formId")
	fm, err := h.formUC.GetByID(c.Request.Context(), formID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "form not found"})
		return
	}
	if !RequireClinicAccess(c, h.userClinicUC, fm.ClinicID) {
		return
	}
	var req form.FieldRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.FormID = formID
	if formVersionId := c.Query("formVersionId"); formVersionId != "" {
		req.FormVersionID = formVersionId
	}
	if req.FormVersionID == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "formVersionId is required (query or body)"})
		return
	}
	resp, err := h.fieldUC.CreateField(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	utils.JSONResponse(c, http.StatusCreated, "form field created", resp, nil)
}

func (h *CustomFormFieldHandler) GetFieldByID(c *gin.Context) {
	id := c.Param("fieldId")
	resp, err := h.fieldUC.GetFieldByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	fm, _ := h.formUC.GetByID(c.Request.Context(), resp.FormID)
	if fm != nil && !RequireClinicAccess(c, h.userClinicUC, fm.ClinicID) {
		return
	}
	utils.JSONResponse(c, http.StatusOK, "form field retrieved", resp, nil)
}

func (h *CustomFormFieldHandler) GetFieldsByFormID(c *gin.Context) {
	formID := c.Param("formId")
	if !RequireClinicAccessFromForm(c, h.formUC, h.userClinicUC, formID) {
		return
	}
	list, err := h.fieldUC.GetFieldsByFormID(c.Request.Context(), formID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	utils.JSONResponse(c, http.StatusOK, "form fields retrieved", list, nil)
}

func (h *CustomFormFieldHandler) GetFieldsByFormVersionID(c *gin.Context) {
	formID := c.Param("formId")
	formVersionID := c.Param("formVersionId")
	if !RequireClinicAccessFromForm(c, h.formUC, h.userClinicUC, formID) {
		return
	}
	list, err := h.fieldUC.GetFieldsByFormVersionID(c.Request.Context(), formVersionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	utils.JSONResponse(c, http.StatusOK, "form fields retrieved", list, nil)
}

func (h *CustomFormFieldHandler) UpdateField(c *gin.Context) {
	fieldID := c.Param("fieldId")
	existing, err := h.fieldUC.GetFieldByID(c.Request.Context(), fieldID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	fm, _ := h.formUC.GetByID(c.Request.Context(), existing.FormID)
	if fm != nil && !RequireClinicAccess(c, h.userClinicUC, fm.ClinicID) {
		return
	}
	var req form.FieldRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := h.fieldUC.UpdateField(c.Request.Context(), fieldID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	utils.JSONResponse(c, http.StatusOK, "form field updated", resp, nil)
}

func (h *CustomFormFieldHandler) DeleteField(c *gin.Context) {
	fieldID := c.Param("fieldId")
	existing, err := h.fieldUC.GetFieldByID(c.Request.Context(), fieldID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	fm, _ := h.formUC.GetByID(c.Request.Context(), existing.FormID)
	if fm != nil && !RequireClinicAccess(c, h.userClinicUC, fm.ClinicID) {
		return
	}
	if err := h.fieldUC.DeleteField(c.Request.Context(), fieldID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	utils.JSONResponse(c, http.StatusOK, "form field deleted", nil, nil)
}

// RequireClinicAccessFromForm loads form by id and enforces clinic access.
func RequireClinicAccessFromForm(c *gin.Context, formUC *usecase.CustomFormService, userClinicUC *usecase.UserClinicService, formID string) bool {
	fm, err := formUC.GetByID(c.Request.Context(), formID)
	if err != nil || fm == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "form not found"})
		return false
	}
	return RequireClinicAccess(c, userClinicUC, fm.ClinicID)
}

// --- Field configs ---

func (h *CustomFormFieldHandler) CreateFieldConfig(c *gin.Context) {
	fieldID := c.Param("fieldId")
	existing, err := h.fieldUC.GetFieldByID(c.Request.Context(), fieldID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "form field not found"})
		return
	}
	if !RequireClinicAccessFromForm(c, h.formUC, h.userClinicUC, existing.FormID) {
		return
	}
	var req form.FieldConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.FormFieldID = fieldID
	resp, err := h.fieldUC.CreateFieldConfig(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	utils.JSONResponse(c, http.StatusCreated, "field config created", resp, nil)
}

func (h *CustomFormFieldHandler) GetFieldConfigByID(c *gin.Context) {
	configID := c.Param("configId")
	resp, err := h.fieldUC.GetFieldConfigByID(c.Request.Context(), configID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	field, _ := h.fieldUC.GetFieldByID(c.Request.Context(), resp.FormFieldID)
	if field != nil {
		fm, _ := h.formUC.GetByID(c.Request.Context(), field.FormID)
		if fm != nil && !RequireClinicAccess(c, h.userClinicUC, fm.ClinicID) {
			return
		}
	}
	utils.JSONResponse(c, http.StatusOK, "field config retrieved", resp, nil)
}

func (h *CustomFormFieldHandler) GetFieldConfigsByFormFieldID(c *gin.Context) {
	fieldID := c.Param("fieldId")
	existing, err := h.fieldUC.GetFieldByID(c.Request.Context(), fieldID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "form field not found"})
		return
	}
	if !RequireClinicAccessFromForm(c, h.formUC, h.userClinicUC, existing.FormID) {
		return
	}
	list, err := h.fieldUC.GetFieldConfigsByFormFieldID(c.Request.Context(), fieldID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	utils.JSONResponse(c, http.StatusOK, "field configs retrieved", list, nil)
}

func (h *CustomFormFieldHandler) UpdateFieldConfig(c *gin.Context) {
	configID := c.Param("configId")
	existing, err := h.fieldUC.GetFieldConfigByID(c.Request.Context(), configID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	field, _ := h.fieldUC.GetFieldByID(c.Request.Context(), existing.FormFieldID)
	if field != nil {
		fm, _ := h.formUC.GetByID(c.Request.Context(), field.FormID)
		if fm != nil && !RequireClinicAccess(c, h.userClinicUC, fm.ClinicID) {
			return
		}
	}
	var req form.FieldConfigRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.FormFieldID = existing.FormFieldID
	resp, err := h.fieldUC.UpdateFieldConfig(c.Request.Context(), configID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	utils.JSONResponse(c, http.StatusOK, "field config updated", resp, nil)
}

func (h *CustomFormFieldHandler) DeleteFieldConfig(c *gin.Context) {
	configID := c.Param("configId")
	existing, err := h.fieldUC.GetFieldConfigByID(c.Request.Context(), configID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	field, _ := h.fieldUC.GetFieldByID(c.Request.Context(), existing.FormFieldID)
	if field != nil {
		fm, _ := h.formUC.GetByID(c.Request.Context(), field.FormID)
		if fm != nil && !RequireClinicAccess(c, h.userClinicUC, fm.ClinicID) {
			return
		}
	}
	if err := h.fieldUC.DeleteFieldConfig(c.Request.Context(), configID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	utils.JSONResponse(c, http.StatusOK, "field config deleted", nil, nil)
}

// --- Formula sources ---

func (h *CustomFormFieldHandler) CreateFormulaSource(c *gin.Context) {
	configID := c.Param("configId")
	cfg, err := h.fieldUC.GetFieldConfigByID(c.Request.Context(), configID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "field config not found"})
		return
	}
	field, _ := h.fieldUC.GetFieldByID(c.Request.Context(), cfg.FormFieldID)
	if field != nil {
		if !RequireClinicAccessFromForm(c, h.formUC, h.userClinicUC, field.FormID) {
			return
		}
	}
	var req form.FormulaSourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.FieldConfigID = configID
	resp, err := h.fieldUC.CreateFormulaSource(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	utils.JSONResponse(c, http.StatusCreated, "formula source created", resp, nil)
}

func (h *CustomFormFieldHandler) GetFormulaSourceByID(c *gin.Context) {
	sourceID := c.Param("sourceId")
	resp, err := h.fieldUC.GetFormulaSourceByID(c.Request.Context(), sourceID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	cfg, _ := h.fieldUC.GetFieldConfigByID(c.Request.Context(), resp.FieldConfigID)
	if cfg != nil {
		field, _ := h.fieldUC.GetFieldByID(c.Request.Context(), cfg.FormFieldID)
		if field != nil {
			fm, _ := h.formUC.GetByID(c.Request.Context(), field.FormID)
			if fm != nil && !RequireClinicAccess(c, h.userClinicUC, fm.ClinicID) {
				return
			}
		}
	}
	utils.JSONResponse(c, http.StatusOK, "formula source retrieved", resp, nil)
}

func (h *CustomFormFieldHandler) GetFormulaSourcesByFieldConfigID(c *gin.Context) {
	configID := c.Param("configId")
	cfg, err := h.fieldUC.GetFieldConfigByID(c.Request.Context(), configID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "field config not found"})
		return
	}
	field, _ := h.fieldUC.GetFieldByID(c.Request.Context(), cfg.FormFieldID)
	if field != nil {
		if !RequireClinicAccessFromForm(c, h.formUC, h.userClinicUC, field.FormID) {
			return
		}
	}
	list, err := h.fieldUC.GetFormulaSourcesByFieldConfigID(c.Request.Context(), configID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	utils.JSONResponse(c, http.StatusOK, "formula sources retrieved", list, nil)
}

func (h *CustomFormFieldHandler) UpdateFormulaSource(c *gin.Context) {
	sourceID := c.Param("sourceId")
	existing, err := h.fieldUC.GetFormulaSourceByID(c.Request.Context(), sourceID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	cfg, _ := h.fieldUC.GetFieldConfigByID(c.Request.Context(), existing.FieldConfigID)
	if cfg != nil {
		field, _ := h.fieldUC.GetFieldByID(c.Request.Context(), cfg.FormFieldID)
		if field != nil {
			fm, _ := h.formUC.GetByID(c.Request.Context(), field.FormID)
			if fm != nil && !RequireClinicAccess(c, h.userClinicUC, fm.ClinicID) {
				return
			}
		}
	}
	var req form.FormulaSourceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.FieldConfigID = existing.FieldConfigID
	resp, err := h.fieldUC.UpdateFormulaSource(c.Request.Context(), sourceID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	utils.JSONResponse(c, http.StatusOK, "formula source updated", resp, nil)
}

func (h *CustomFormFieldHandler) DeleteFormulaSource(c *gin.Context) {
	sourceID := c.Param("sourceId")
	existing, err := h.fieldUC.GetFormulaSourceByID(c.Request.Context(), sourceID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	cfg, _ := h.fieldUC.GetFieldConfigByID(c.Request.Context(), existing.FieldConfigID)
	if cfg != nil {
		field, _ := h.fieldUC.GetFieldByID(c.Request.Context(), cfg.FormFieldID)
		if field != nil {
			fm, _ := h.formUC.GetByID(c.Request.Context(), field.FormID)
			if fm != nil && !RequireClinicAccess(c, h.userClinicUC, fm.ClinicID) {
				return
			}
		}
	}
	if err := h.fieldUC.DeleteFormulaSource(c.Request.Context(), sourceID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	utils.JSONResponse(c, http.StatusOK, "formula source deleted", nil, nil)
}
