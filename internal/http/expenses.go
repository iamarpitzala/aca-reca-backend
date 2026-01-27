package http

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
	"github.com/iamarpitzala/aca-reca-backend/internal/service"
)

type ExpensesHandler struct {
	expensesService *service.ExpensesService
}

func NewExpensesHandler(expensesService *service.ExpensesService) *ExpensesHandler {
	return &ExpensesHandler{
		expensesService: expensesService,
	}
}

// CreateExpenseType creates a new expense type
// POST /api/v1/expenses/type
// @Summary Create a new expense type
// @Description Create a new expense type with the given information
// @Tags Expenses
// @Accept json
// @Produce json
// @Param expenseType body domain.ExpenseType true "Expense type information"
// @Success 201 {object} domain.ExpenseType
// @Failure 400 {object} domain.H
// @Failure 500 {object} domain.H
// @Router /expenses/type [post]
func (h *ExpensesHandler) CreateExpenseType(c *gin.Context) {
	var expenseType domain.ExpenseType
	if err := c.ShouldBindJSON(&expenseType); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	userIDUUID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user ID"})
		return
	}

	// Set required fields
	expenseType.ID = uuid.New()
	expenseType.CreatedBy = userIDUUID
	expenseType.CreatedAt = time.Now()

	err := h.expensesService.CreateExpenseType(c.Request.Context(), &expenseType)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "expense type created successfully", "expenseType": expenseType})
}

// CreateExpenseCategory creates a new expense category
// POST /api/v1/expenses/category
// @Summary Create a new expense category
// @Description Create a new expense category with the given information
// @Tags Expenses
// @Accept json
// @Produce json
// @Param expenseCategory body domain.ExpenseCategory true "Expense category information"
// @Success 201 {object} domain.ExpenseCategory
// @Failure 400 {object} domain.H
// @Failure 500 {object} domain.H
// @Router /expenses/category [post]
func (h *ExpensesHandler) CreateExpenseCategory(c *gin.Context) {
	var expenseCategory domain.ExpenseCategory
	if err := c.ShouldBindJSON(&expenseCategory); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	userIDUUID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user ID"})
		return
	}

	// Set required fields
	expenseCategory.ID = uuid.New()
	expenseCategory.CreatedBy = userIDUUID
	expenseCategory.CreatedAt = time.Now()

	err := h.expensesService.CreateExpenseCategory(c.Request.Context(), &expenseCategory)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "expense category created successfully", "expenseCategory": expenseCategory})
}

// CreateExpenseCategoryType creates a new expense category type
// POST /api/v1/expenses/category-type
// @Summary Create a new expense category type
// @Description Create a new expense category type with the given information
// @Tags Expenses
// @Accept json
// @Produce json
// @Param expenseCategoryType body domain.ExpenseCategoryType true "Expense category type information"
// @Success 201 {object} domain.ExpenseCategoryType
// @Failure 400 {object} domain.H
// @Failure 500 {object} domain.H
// @Router /expenses/category-type [post]
func (h *ExpensesHandler) CreateExpenseCategoryType(c *gin.Context) {
	var expenseCategoryType domain.ExpenseCategoryType
	if err := c.ShouldBindJSON(&expenseCategoryType); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	userIDUUID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user ID"})
		return
	}

	// Set required fields
	expenseCategoryType.ID = uuid.New()
	expenseCategoryType.CreatedBy = userIDUUID
	expenseCategoryType.CreatedAt = time.Now()

	err := h.expensesService.CreateExpenseCategoryType(c.Request.Context(), &expenseCategoryType)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "expense category type created successfully", "expenseCategoryType": expenseCategoryType})
}

// CreateExpenseEntry creates a new expense entry
// POST /api/v1/expenses/entry
// @Summary Create a new expense entry
// @Description Create a new expense entry with the given information
// @Tags Expenses
// @Accept json
// @Produce json
// @Param expenseEntry body domain.ExpenseEntry true "Expense entry information"
// @Success 201 {object} domain.ExpenseEntry
// @Failure 400 {object} domain.H
// @Failure 500 {object} domain.H
// @Router /expenses/entry [post]
func (h *ExpensesHandler) CreateExpenseEntry(c *gin.Context) {
	var expenseEntry domain.ExpenseEntry
	if err := c.ShouldBindJSON(&expenseEntry); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID from context
	userID, exists := c.Get("user_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "user not authenticated"})
		return
	}
	userIDUUID, ok := userID.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid user ID"})
		return
	}

	// Set required fields
	expenseEntry.ID = uuid.New()
	expenseEntry.CreatedBy = userIDUUID
	expenseEntry.CreatedAt = time.Now()
	expenseEntry.DeletedAt = nil

	err := h.expensesService.CreateExpenseEntry(c.Request.Context(), &expenseEntry)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "expense entry created successfully", "expenseEntry": expenseEntry})
}

// GetExpenseTypeByID retrieves a expense type by ID
// GET /api/v1/expenses/type/:id
// @Summary Retrieve a expense type by ID
// @Description Retrieve a expense type by ID
// @Tags Expenses
// @Accept json
// @Produce json
// @Param id path string true "Expense Type ID"
// @Success 200 {object} domain.ExpenseType
// @Failure 400 {object} domain.H
// @Failure 404 {object} domain.H
// @Failure 500 {object} domain.H
// @Router /expenses/type/{id} [get]
func (h *ExpensesHandler) GetExpenseTypeByID(c *gin.Context) {
	id := c.Param("id")
	expenseType, err := h.expensesService.GetExpenseTypeByID(c.Request.Context(), uuid.MustParse(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "expense type retrieved successfully", "expenseType": expenseType})
}

// GetExpenseCategoryByID retrieves a expense category by ID
// GET /api/v1/expenses/category/:id
// @Summary Retrieve a expense category by ID
// @Description Retrieve a expense category by ID
// @Tags Expenses
// @Accept json
// @Produce json
// @Param id path string true "Expense Category ID"
// @Success 200 {object} domain.ExpenseCategory
// @Failure 400 {object} domain.H
// @Failure 404 {object} domain.H
// @Failure 500 {object} domain.H
// @Router /expenses/category/{id} [get]
func (h *ExpensesHandler) GetExpenseCategoryByID(c *gin.Context) {
	id := c.Param("id")
	expenseCategory, err := h.expensesService.GetExpenseCategoryByID(c.Request.Context(), uuid.MustParse(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "expense category retrieved successfully", "expenseCategory": expenseCategory})
}

// GetExpenseCategoryTypeByID retrieves a expense category type by ID
// GET /api/v1/expenses/category-type/:id
// @Summary Retrieve a expense category type by ID
// @Description Retrieve a expense category type by ID
// @Tags Expenses
// @Accept json
// @Produce json
// @Param id path string true "Expense Category Type ID"
// @Success 200 {object} domain.ExpenseCategoryType
// @Failure 400 {object} domain.H
// @Failure 404 {object} domain.H
// @Failure 500 {object} domain.H
// @Router /expenses/category-type/{id} [get]
func (h *ExpensesHandler) GetExpenseCategoryTypeByID(c *gin.Context) {
	id := c.Param("id")
	expenseCategoryType, err := h.expensesService.GetExpenseCategoryTypeByID(c.Request.Context(), uuid.MustParse(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "expense category type retrieved successfully", "expenseCategoryType": expenseCategoryType})
}

// GetExpenseEntryByID retrieves a expense entry by ID
// GET /api/v1/expenses/entry/:id
// @Summary Retrieve a expense entry by ID
// @Description Retrieve a expense entry by ID
// @Tags Expenses
// @Accept json
// @Produce json
// @Param id path string true "Expense Entry ID"
// @Success 200 {object} domain.ExpenseEntry
// @Failure 400 {object} domain.H
// @Failure 404 {object} domain.H
// @Failure 500 {object} domain.H
// @Router /expenses/entry/{id} [get]
func (h *ExpensesHandler) GetExpenseEntryByID(c *gin.Context) {
	id := c.Param("id")
	expenseEntry, err := h.expensesService.GetExpenseEntryByID(c.Request.Context(), uuid.MustParse(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "expense entry retrieved successfully", "expenseEntry": expenseEntry})
}
