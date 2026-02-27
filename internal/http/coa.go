package http

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/iamarpitzala/aca-reca-backend/internal/application/usecase"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain/coa"
	"github.com/iamarpitzala/aca-reca-backend/util"
	utils "github.com/iamarpitzala/aca-reca-backend/util"
)

type COAHandler struct {
	coaUC *usecase.COAService
}

func NewCOAHandler(coaUC *usecase.COAService) *COAHandler {
	return &COAHandler{
		coaUC: coaUC,
	}
}

// CreateCOA creates a new coa
// POST /api/v1/coa
// @Summary Create a new coa
// @Description Create a new coa with the given information
// @Tags COA
// @Accept json
// @Produce json
// @Param coa body coa.COARequest true "COA information"
// @Success 201 {object} domain.COAResponse
// @Failure 400 {object} domain.H
// @Failure 500 {object} domain.H
// @Router /coa [post]
func (h *COAHandler) CreateCOA(c *gin.Context) {
	userID, ok := GetAuthUserID(c)
	if !ok {
		return
	}
	var coa coa.COARequest
	if err := utils.BindAndValidate(c, &coa); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := h.coaUC.CreateCOA(c.Request.Context(), &coa, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	utils.JSONResponse(c, http.StatusCreated, "coa created successfully", nil, nil)
}

// GetAllCOAs returns all chart of accounts entries
// GET /api/v1/coa
// @Summary Get all COAs
// @Description Get all chart of accounts entries
// @Tags COA
// @Success 200 {array} domain.COAResponse
// @Router /coa [get]
func (h *COAHandler) GetAllCOAs(c *gin.Context) {
	response, err := h.coaUC.GetAllCOAs(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	utils.JSONResponse(c, http.StatusOK, "coas retrieved successfully", response, nil)
}

func (h *COAHandler) GetAllCOAType(c *gin.Context) {
	response, err := h.coaUC.GetCOAType(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	utils.JSONResponse(c, http.StatusOK, "account types retrieved successfully", response, nil)
}

// GetAllAccountTax gets all account tax types
// GET /api/v1/coa/tax
// @Summary Get all account tax types
// @Description Get all account tax types
// @Tags COA
// @Accept json
// @Produce json
// @Success 200 {object} domain.AccountTax
// @Failure 500 {object} domain.H
// @Router /coa/tax [get]
func (h *COAHandler) GetAllAccountTax(c *gin.Context) {
	response, err := h.coaUC.GetAccountTax(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	utils.JSONResponse(c, http.StatusOK, "account taxes retrieved successfully", response, nil)
}

// GetCOAByAccountTaxID gets accounts by account tax id
// GET /api/v1/coa/account-tax/:id
// @Summary Get a coa by account tax id
// @Description Get a coa by account tax id
// @Tags COA
// @Accept json
// @Produce json
// @Param id path string true "COA ID"
// @Success 200 {object} domain.COAResponse
// @Failure 400 {object} domain.H
// @Failure 404 {object} domain.H
// @Failure 500 {object} domain.H
func (h *COAHandler) GetCOAByAccountTaxID(c *gin.Context) {
	accountTaxId := c.Param("id")
	accountTaxIdInt, ok := util.ToInt(accountTaxId)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.ErrInvalidAccountTaxID})
		return
	}
	response, err := h.coaUC.GetCOAByAccountTaxID(c.Request.Context(), accountTaxIdInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	utils.JSONResponse(c, http.StatusOK, "coa retrieved successfully", response, nil)
}

// GetCOAByID gets a coa by id
// GET /api/v1/coa/:id
// @Summary Get a coa by id
// @Description Get a coa by id
// @Tags COA
// @Accept json
// @Produce json
// @Param id path string true "COA ID"
// @Success 200 {object} domain.COAResponse
// @Failure 400 {object} domain.H
// @Failure 404 {object} domain.H
// @Failure 500 {object} domain.H
func (h *COAHandler) GetCOAByID(c *gin.Context) {
	id := c.Param("id")

	response, err := h.coaUC.GetCOAByID(c.Request.Context(), id)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": utils.ErrCOANotFound})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	utils.JSONResponse(c, http.StatusOK, "coa retrieved successfully", response, nil)
}

// GetCOAByCode gets a coa by code
// GET /api/v1/coa/code/:code
// @Summary Get a coa by code
// @Description Get a coa by code
// @Tags COA
// @Accept json
// @Produce json
// @Param code path string true "COA Code"
// @Success 200 {object} domain.COAResponse
// @Failure 400 {object} domain.H
// @Failure 404 {object} domain.H
// @Failure 500 {object} domain.H
func (h *COAHandler) GetCOAByCode(c *gin.Context) {
	code := c.Param("code")
	response, err := h.coaUC.GetCOAByCode(c.Request.Context(), code)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": utils.ErrCOANotFound})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	utils.JSONResponse(c, http.StatusOK, "coa retrieved successfully", response, nil)
}

// GetCOAByAccountTypeID gets a coa by account type id
// GET /api/v1/coa/account-type/:id
// @Summary Get a coa by account type id
// @Description Get a coa by account type id
// @Tags COA
// @Accept json
// @Produce json
// @Param id path string true "COA ID"
// @Success 200 {object} domain.COAResponse
// @Failure 400 {object} domain.H
// @Failure 404 {object} domain.H
// @Failure 500 {object} domain.H
func (h *COAHandler) GetCOAByAccountTypeID(c *gin.Context) {
	accountTypeId := c.Param("id")
	accountTypeIdInt, ok := util.ToInt(accountTypeId)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.ErrInvalidAccountTypeID})
		return
	}
	sortBy := c.DefaultQuery("sort", "code")
	if sortBy != "code" && sortBy != "name" {
		sortBy = "code"
	}
	sortOrder := c.DefaultQuery("order", "asc")
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "asc"
	}
	response, err := h.coaUC.GetCOAByAccountTypeID(c.Request.Context(), accountTypeIdInt, sortBy, sortOrder)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": utils.ErrCOANotFound})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	utils.JSONResponse(c, http.StatusOK, "coa retrieved successfully", response, nil)
}

// GetCOAsByAccountType returns all accounts with optional sort. Query: sort=code|name, order=asc|desc
// GET /api/v1/coa/account-types
// @Summary Get all accounts (optionally sorted)
// @Description Get all accounts. Query params: sort (code|name, default code), order (asc|desc, default asc)
// @Tags COA
// @Param sort query string false "Sort field: code or name"
// @Param order query string false "Sort order: asc or desc"
// @Success 200 {array} domain.COAResponse
// @Router /coa/account-types [get]
func (h *COAHandler) GetCOAsByAccountType(c *gin.Context) {
	sortBy := c.DefaultQuery("sort", "code")
	if sortBy != "code" && sortBy != "name" {
		sortBy = "code"
	}
	sortOrder := c.DefaultQuery("order", "asc")
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "asc"
	}
	response, err := h.coaUC.GetCOAsByAccountType(c.Request.Context(), sortBy, sortOrder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	utils.JSONResponse(c, http.StatusOK, "coas retrieved successfully", response, nil)
}

// UpdateCOA updates a coa
// PUT /api/v1/coa/:id
// @Summary Update a coa
// @Description Update a coa with the given information
// @Tags COA
// @Accept json
// @Produce json
// @Param id path string true "COA ID"
// @Param coa body domain.COARequest true "COA information"
// @Success 200 {object} domain.COAResponse
// @Failure 400 {object} domain.H
// @Failure 404 {object} domain.H
// @Failure 500 {object} domain.H
func (h *COAHandler) UpdateCOA(c *gin.Context) {
	id := c.Param("id")

	var coa coa.COARequest
	if err := util.BindAndValidate(c, &coa); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	coaRepo := coa.ToRepo()
	coaRepo.ID = id
	err := h.coaUC.UpdateCOA(c.Request.Context(), coaRepo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	utils.JSONResponse(c, http.StatusOK, "coa updated successfully", coaRepo.ToResponse(), nil)
}

// DeleteCOA deletes a coa
// DELETE /api/v1/coa
// @Summary Delete a coa
// @Description Delete a coa
// @Tags COA
// @Accept json
// @Produce json
// @Param ids body domain.BulkDeleteCOARequest true "COA IDs"
// @Success 200 {object} domain.COAResponse
// @Failure 400 {object} domain.H
// @Failure 404 {object} domain.H
// @Failure 500 {object} domain.H
func (h *COAHandler) DeleteCOA(c *gin.Context) {
	var req coa.BulkDeleteCOARequest
	if err := utils.BindAndValidate(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.ErrIDsRequired})
		return
	}
	if len(req.IDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.ErrIDsMustNotBeEmpty})
		return
	}
	ids := make([]string, 0, len(req.IDs))
	for _, s := range req.IDs {
		ids = append(ids, s)
	}
	err := h.coaUC.DeleteCOA(c.Request.Context(), ids)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	utils.JSONResponse(c, http.StatusOK, "coa deleted successfully", gin.H{"deleted": len(ids)}, nil)
}

// BulkUpdateCOATax updates account_tax_id for multiple accounts. PATCH /coa/bulk-tax with body { ids, accountTaxId }.
func (h *COAHandler) BulkUpdateCOATax(c *gin.Context) {
	var req coa.BulkUpdateCOATaxRequest
	if err := utils.BindAndValidate(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(req.IDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.ErrIDsMustNotBeEmpty})
		return
	}
	ids := make([]string, 0, len(req.IDs))
	for _, s := range req.IDs {
		ids = append(ids, s)
	}
	if err := h.coaUC.BulkUpdateCOATax(c.Request.Context(), ids, req.AccountTaxID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	utils.JSONResponse(c, http.StatusOK, "tax updated successfully", gin.H{"updated": len(ids)}, nil)
}

// ArchiveCOA soft-deletes (archives) multiple accounts. PATCH /coa/archive with body { ids }.
func (h *COAHandler) ArchiveCOA(c *gin.Context) {
	var req coa.BulkArchiveCOARequest
	if err := utils.BindAndValidate(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.ErrIDsRequired})
		return
	}
	if len(req.IDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.ErrIDsMustNotBeEmpty})
		return
	}
	ids := make([]string, 0, len(req.IDs))
	for _, s := range req.IDs {
		ids = append(ids, s)
	}
	if err := h.coaUC.DeleteCOA(c.Request.Context(), ids); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	utils.JSONResponse(c, http.StatusOK, "accounts archived successfully", gin.H{"archived": len(ids)}, nil)
}

func (h *COAHandler) CheckIfAccountTaxIsTaxable(c *gin.Context) {
	accountTypeID := c.Param("id")
	accountTypeIDInt, ok := util.ToInt(accountTypeID)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.ErrInvalidAccountTypeID})
		return
	}
	isTaxable, err := h.coaUC.CheckIfAccountTaxIsTaxable(c.Request.Context(), accountTypeIDInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	utils.JSONResponse(c, http.StatusOK, "account type is taxable", gin.H{"isTaxable": isTaxable}, nil)
}
