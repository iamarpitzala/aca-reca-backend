package http

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/application/usecase"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
	utils "github.com/iamarpitzala/aca-reca-backend/util"
)

type TransactionHandler struct {
	txnUC *usecase.TransactionService
}

func NewTransactionHandler(txnUC *usecase.TransactionService) *TransactionHandler {
	return &TransactionHandler{txnUC: txnUC}
}

// Create creates a new transaction with ledger lines
// POST /api/v1/transaction
func (h *TransactionHandler) Create(c *gin.Context) {
	var req domain.CreateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": utils.ErrUserNotAuthenticated})
		return
	}
	userIDUUID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": utils.ErrInvalidUserID})
		return
	}

	result, err := h.txnUC.Create(c.Request.Context(), &req, userIDUUID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":     utils.MsgTransactionCreated,
		"transaction": result.Transaction,
		"ledger":      result.Ledger,
	})
}

// GetByID retrieves a transaction by ID
// GET /api/v1/transaction/:id
func (h *TransactionHandler) GetByID(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.ErrInvalidID})
		return
	}

	result, err := h.txnUC.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     utils.MsgTransactionRetrieved,
		"transaction": result.Transaction,
		"ledger":      result.Ledger,
	})
}

// GetByClinicID lists all transactions for a clinic
// GET /api/v1/transaction/clinic/:clinicId
func (h *TransactionHandler) GetByClinicID(c *gin.Context) {
	clinicID, err := uuid.Parse(c.Param("clinicId"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.ErrInvalidClinicID})
		return
	}

	results, err := h.txnUC.GetByClinicID(c.Request.Context(), clinicID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":      utils.MsgTransactionsRetrieved,
		"transactions": results,
		"total":        len(results),
	})
}

// Update updates a DRAFT transaction
// PUT /api/v1/transaction/:id
func (h *TransactionHandler) Update(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.ErrInvalidID})
		return
	}

	var req domain.UpdateTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	result, err := h.txnUC.Update(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     utils.MsgTransactionUpdated,
		"transaction": result.Transaction,
		"ledger":      result.Ledger,
	})
}

// Post changes a DRAFT transaction to POSTED
// POST /api/v1/transaction/:id/post
func (h *TransactionHandler) Post(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.ErrInvalidID})
		return
	}

	result, err := h.txnUC.PostTransaction(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     utils.MsgTransactionPosted,
		"transaction": result.Transaction,
		"ledger":      result.Ledger,
	})
}

// Void marks a transaction as VOIDED
// POST /api/v1/transaction/:id/void
func (h *TransactionHandler) Void(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.ErrInvalidID})
		return
	}

	var req domain.VoidTransactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	txn, err := h.txnUC.VoidTransaction(c.Request.Context(), id, req.Reason)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":     utils.MsgTransactionVoided,
		"transaction": txn,
	})
}

// Delete soft-deletes a transaction
// DELETE /api/v1/transaction/:id
func (h *TransactionHandler) Delete(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.ErrInvalidID})
		return
	}

	if err := h.txnUC.Delete(c.Request.Context(), id); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": utils.MsgTransactionDeleted})
}
