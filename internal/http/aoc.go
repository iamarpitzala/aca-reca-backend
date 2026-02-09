package http

import (
	"database/sql"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
	"github.com/iamarpitzala/aca-reca-backend/internal/service"
	"github.com/iamarpitzala/aca-reca-backend/util"
	utils "github.com/iamarpitzala/aca-reca-backend/util"
)

type AOCHandler struct {
	aocService *service.AOSService
}

func NewAOCHandler(aocService *service.AOSService) *AOCHandler {
	return &AOCHandler{
		aocService: aocService,
	}
}

// CreateAOC creates a new aoc
// POST /api/v1/aoc
// @Summary Create a new aoc
// @Description Create a new aoc with the given information
// @Tags AOC
// @Accept json
// @Produce json
// @Param aoc body domain.AOCRequest true "AOC information"
// @Success 201 {object} domain.AOCResponse
// @Failure 400 {object} domain.H
// @Failure 500 {object} domain.H
// @Router /aoc [post]
func (h *AOCHandler) CreateAOC(c *gin.Context) {
	var aoc domain.AOCRequest
	if err := util.BindAndValidate(c, &aoc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := h.aocService.CreateAOC(c.Request.Context(), &aoc)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	utils.JSONResponse(c, http.StatusCreated, "aoc created successfully", nil, nil)
}

// GetAllAOCs returns all chart of accounts entries
// GET /api/v1/aoc
// @Summary Get all AOCs
// @Description Get all chart of accounts entries
// @Tags AOC
// @Success 200 {array} domain.AOCResponse
// @Router /aoc [get]
func (h *AOCHandler) GetAllAOCs(c *gin.Context) {
	response, err := h.aocService.GetAllAOCs(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	utils.JSONResponse(c, http.StatusOK, "aocs retrieved successfully", response, nil)
}

func (h *AOCHandler) GetAllAOCType(c *gin.Context) {
	response, err := h.aocService.GetAOCType(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	utils.JSONResponse(c, http.StatusOK, "account types retrieved successfully", response, nil)
}

// GetAllAccountTax gets all account tax types
// GET /api/v1/aoc/tax
// @Summary Get all account tax types
// @Description Get all account tax types
// @Tags AOC
// @Accept json
// @Produce json
// @Success 200 {object} domain.AccountTax
// @Failure 500 {object} domain.H
// @Router /aoc/tax [get]
func (h *AOCHandler) GetAllAccountTax(c *gin.Context) {
	response, err := h.aocService.GetAccountTax(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	utils.JSONResponse(c, http.StatusOK, "account taxes retrieved successfully", response, nil)
}

// GetAOCByAccountTaxID gets accounts by account tax id
// GET /api/v1/aoc/account-tax/:id
// @Summary Get a aoc by account tax id
// @Description Get a aoc by account tax id
// @Tags AOC
// @Accept json
// @Produce json
// @Param id path string true "AOC ID"
// @Success 200 {object} domain.AOCResponse
// @Failure 400 {object} domain.H
// @Failure 404 {object} domain.H
// @Failure 500 {object} domain.H
func (h *AOCHandler) GetAOCByAccountTaxID(c *gin.Context) {
	accountTaxId := c.Param("id")
	accountTaxIdInt, ok := util.ToInt(accountTaxId)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account tax id"})
		return
	}
	response, err := h.aocService.GetAOCByAccountTaxID(c.Request.Context(), accountTaxIdInt)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	utils.JSONResponse(c, http.StatusOK, "aoc retrieved successfully", response, nil)
}

// GetAOCByID gets a aoc by id
// GET /api/v1/aoc/:id
// @Summary Get a aoc by id
// @Description Get a aoc by id
// @Tags AOC
// @Accept json
// @Produce json
// @Param id path string true "AOC ID"
// @Success 200 {object} domain.AOCResponse
// @Failure 400 {object} domain.H
// @Failure 404 {object} domain.H
// @Failure 500 {object} domain.H
// reservedPathSegments are path segments that must not be treated as UUIDs (avoid 400 when route is misrouted).
var reservedPathSegments = map[string]bool{
	"account-types": true, "account-type": true, "account-tax": true, "type": true, "tax": true, "code": true,
}

func (h *AOCHandler) GetAOCByID(c *gin.Context) {
	id := c.Param("id")
	if reservedPathSegments[id] {
		c.JSON(http.StatusNotFound, gin.H{"error": "aoc not found"})
		return
	}
	idUUID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	response, err := h.aocService.GetAOCByID(c.Request.Context(), idUUID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "aoc not found"})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	utils.JSONResponse(c, http.StatusOK, "aoc retrieved successfully", response, nil)
}

// GetAOCByCode gets a aoc by code
// GET /api/v1/aoc/code/:code
// @Summary Get a aoc by code
// @Description Get a aoc by code
// @Tags AOC
// @Accept json
// @Produce json
// @Param code path string true "AOC Code"
// @Success 200 {object} domain.AOCResponse
// @Failure 400 {object} domain.H
// @Failure 404 {object} domain.H
// @Failure 500 {object} domain.H
func (h *AOCHandler) GetAOCByCode(c *gin.Context) {
	code := c.Param("code")
	response, err := h.aocService.GetAOCByCode(c.Request.Context(), code)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "aoc not found"})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	utils.JSONResponse(c, http.StatusOK, "aoc retrieved successfully", response, nil)
}

// GetAOCByAccountTypeID gets a aoc by account type id
// GET /api/v1/aoc/account-type/:id
// @Summary Get a aoc by account type id
// @Description Get a aoc by account type id
// @Tags AOC
// @Accept json
// @Produce json
// @Param id path string true "AOC ID"
// @Success 200 {object} domain.AOCResponse
// @Failure 400 {object} domain.H
// @Failure 404 {object} domain.H
// @Failure 500 {object} domain.H
func (h *AOCHandler) GetAOCByAccountTypeID(c *gin.Context) {
	accountTypeId := c.Param("id")
	accountTypeIdInt, ok := util.ToInt(accountTypeId)
	if !ok {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account type id"})
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
	response, err := h.aocService.GetAOCByAccountTypeID(c.Request.Context(), accountTypeIdInt, sortBy, sortOrder)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "aoc not found"})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	utils.JSONResponse(c, http.StatusOK, "aoc retrieved successfully", response, nil)
}

// GetAOCsByAccountType returns all accounts with optional sort. Query: sort=code|name, order=asc|desc
// GET /api/v1/aoc/account-types
// @Summary Get all accounts (optionally sorted)
// @Description Get all accounts. Query params: sort (code|name, default code), order (asc|desc, default asc)
// @Tags AOC
// @Param sort query string false "Sort field: code or name"
// @Param order query string false "Sort order: asc or desc"
// @Success 200 {array} domain.AOCResponse
// @Router /aoc/account-types [get]
func (h *AOCHandler) GetAOCsByAccountType(c *gin.Context) {
	sortBy := c.DefaultQuery("sort", "code")
	if sortBy != "code" && sortBy != "name" {
		sortBy = "code"
	}
	sortOrder := c.DefaultQuery("order", "asc")
	if sortOrder != "asc" && sortOrder != "desc" {
		sortOrder = "asc"
	}
	response, err := h.aocService.GetAOCsByAccountType(c.Request.Context(), sortBy, sortOrder)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	utils.JSONResponse(c, http.StatusOK, "aocs retrieved successfully", response, nil)
}

// UpdateAOC updates a aoc
// PUT /api/v1/aoc/:id
// @Summary Update a aoc
// @Description Update a aoc with the given information
// @Tags AOC
// @Accept json
// @Produce json
// @Param id path string true "AOC ID"
// @Param aoc body domain.AOCRequest true "AOC information"
// @Success 200 {object} domain.AOCResponse
// @Failure 400 {object} domain.H
// @Failure 404 {object} domain.H
// @Failure 500 {object} domain.H
func (h *AOCHandler) UpdateAOC(c *gin.Context) {
	id := c.Param("id")
	idUUID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	var aoc domain.AOCRequest
	if err := util.BindAndValidate(c, &aoc); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	aocRepo := aoc.ToRepo()
	aocRepo.ID = idUUID
	err = h.aocService.UpdateAOC(c.Request.Context(), aocRepo)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	utils.JSONResponse(c, http.StatusOK, "aoc updated successfully", aocRepo.ToResponse(), nil)
}

// DeleteAOC deletes a aoc
// DELETE /api/v1/aoc
// @Summary Delete a aoc
// @Description Delete a aoc
// @Tags AOC
// @Accept json
// @Produce json
// @Param ids body domain.BulkDeleteAOCRequest true "AOC IDs"
// @Success 200 {object} domain.AOCResponse
// @Failure 400 {object} domain.H
// @Failure 404 {object} domain.H
// @Failure 500 {object} domain.H
func (h *AOCHandler) DeleteAOC(c *gin.Context) {
	var req domain.BulkDeleteAOCRequest
	if err := util.BindAndValidate(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request: ids required"})
		return
	}
	if len(req.IDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ids must not be empty"})
		return
	}
	ids := make([]uuid.UUID, 0, len(req.IDs))
	for _, s := range req.IDs {
		id, err := uuid.Parse(s)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id: " + s})
			return
		}
		ids = append(ids, id)
	}
	err := h.aocService.DeleteAOC(c.Request.Context(), ids)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	utils.JSONResponse(c, http.StatusOK, "aoc deleted successfully", gin.H{"deleted": len(ids)}, nil)
}

// BulkUpdateTax updates account_tax_id for multiple accounts. PATCH /aoc/bulk-tax with body { ids, accountTaxId }.
func (h *AOCHandler) BulkUpdateTax(c *gin.Context) {
	var req domain.BulkUpdateTaxRequest
	if err := util.BindAndValidate(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if len(req.IDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ids must not be empty"})
		return
	}
	ids := make([]uuid.UUID, 0, len(req.IDs))
	for _, s := range req.IDs {
		id, err := uuid.Parse(s)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id: " + s})
			return
		}
		ids = append(ids, id)
	}
	if err := h.aocService.BulkUpdateTax(c.Request.Context(), ids, req.AccountTaxID); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	utils.JSONResponse(c, http.StatusOK, "tax updated successfully", gin.H{"updated": len(ids)}, nil)
}

// ArchiveAOC soft-deletes (archives) multiple accounts. PATCH /aoc/archive with body { ids }.
func (h *AOCHandler) ArchiveAOC(c *gin.Context) {
	var req domain.BulkArchiveAOCRequest
	if err := util.BindAndValidate(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request: ids required"})
		return
	}
	if len(req.IDs) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ids must not be empty"})
		return
	}
	ids := make([]uuid.UUID, 0, len(req.IDs))
	for _, s := range req.IDs {
		id, err := uuid.Parse(s)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id: " + s})
			return
		}
		ids = append(ids, id)
	}
	if err := h.aocService.DeleteAOC(c.Request.Context(), ids); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	utils.JSONResponse(c, http.StatusOK, "accounts archived successfully", gin.H{"archived": len(ids)}, nil)
}
