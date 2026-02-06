package http

import (
	"database/sql"
	"net/http"
	"strconv"

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
// GET /api/v1/aoc/account-tax/:accountTaxId
// @Summary Get a aoc by account tax id
// @Description Get a aoc by account tax id
// @Tags AOC
// @Accept json
// @Produce json
// @Param accountTaxId path string true "AOC Account Tax ID"
// @Success 200 {object} domain.AOCResponse
// @Failure 400 {object} domain.H
// @Failure 404 {object} domain.H
// @Failure 500 {object} domain.H
func (h *AOCHandler) GetAOCByAccountTaxID(c *gin.Context) {
	accountTaxId := c.Param("accountTaxId")
	accountTaxIdInt, err := strconv.Atoi(accountTaxId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account tax id"})
		return
	}
	response, err := h.aocService.GetAOCByAccountTaxID(c.Request.Context(), accountTaxIdInt)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "aoc not found"})
			return
		}
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	utils.JSONResponse(c, http.StatusOK, "aocs retrieved successfully", response, nil)
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
func (h *AOCHandler) GetAOCByID(c *gin.Context) {
	id := c.Param("id")
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
// GET /api/v1/aoc/account-type/:accountTypeId
// @Summary Get a aoc by account type id
// @Description Get a aoc by account type id
// @Tags AOC
// @Accept json
// @Produce json
// @Param accountTypeId path string true "AOC Account Type ID"
// @Success 200 {object} domain.AOCResponse
// @Failure 400 {object} domain.H
// @Failure 404 {object} domain.H
// @Failure 500 {object} domain.H
func (h *AOCHandler) GetAOCByAccountTypeID(c *gin.Context) {
	accountTypeId := c.Param("accountTypeId")
	accountTypeIdInt, err := strconv.Atoi(accountTypeId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid account type id"})
		return
	}
	response, err := h.aocService.GetAOCByAccountTypeID(c.Request.Context(), accountTypeIdInt)
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
// DELETE /api/v1/aoc/:id
// @Summary Delete a aoc
// @Description Delete a aoc
// @Tags AOC
// @Accept json
// @Produce json
// @Param id path string true "AOC ID"
// @Success 200 {object} domain.AOCResponse
// @Failure 400 {object} domain.H
// @Failure 404 {object} domain.H
// @Failure 500 {object} domain.H
func (h *AOCHandler) DeleteAOC(c *gin.Context) {
	id := c.Param("id")
	idUUID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
		return
	}
	err = h.aocService.DeleteAOC(c.Request.Context(), idUUID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	utils.JSONResponse(c, http.StatusOK, "aoc deleted successfully", nil, nil)
}
