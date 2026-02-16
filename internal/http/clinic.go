package http

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/application/usecase"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
	"github.com/iamarpitzala/aca-reca-backend/util"
)

type ClinicHandler struct {
	clinicUC      *usecase.ClinicService
	userClinicUC  *usecase.UserClinicService
	clinicCOAUC   *usecase.ClinicCOAService
}

func NewClinicHandler(clinicUC *usecase.ClinicService, userClinicUC *usecase.UserClinicService, clinicCOAUC *usecase.ClinicCOAService) *ClinicHandler {
	return &ClinicHandler{
		clinicUC:     clinicUC,
		userClinicUC: userClinicUC,
		clinicCOAUC:  clinicCOAUC,
	}
}

// CreateClinic creates a new clinic and associates the creating user as owner
// POST /api/v1/clinic
// @Summary Create a new clinic
// @Description Create a new clinic with the given information. The creating user is automatically associated as owner.
// @Tags Clinic
// @Accept json
// @Produce json
// @Param clinic body domain.Clinic true "Clinic information"
// @Success 201 {object} domain.H
// @Failure 400 {object} domain.H
// @Failure 500 {object} domain.H
// @Router /clinic [post]
func (h *ClinicHandler) CreateClinic(c *gin.Context) {
	userID, ok := GetAuthUserID(c)
	if !ok {
		return
	}
	var clinic domain.Clinic
	if err := c.ShouldBindJSON(&clinic); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := h.clinicUC.CreateClinic(c.Request.Context(), &clinic)
	if err != nil {
		if errors.Is(err, usecase.ErrDuplicateABN) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	// Auto-associate the creating user as owner so they can access the clinic.
	// If this fails, roll back clinic creation so we never leave a clinic the creator cannot access.
	_, err = h.userClinicUC.AssociateUserWithClinic(c.Request.Context(), userID, clinic.ID, util.RoleOwner)
	if err != nil {
		_ = h.clinicUC.DeleteClinic(c.Request.Context(), clinic.ID)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "clinic was created but we could not link you as owner. Please try again.",
			"detail": err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "clinic created successfully", "clinic_id": clinic.ID, "clinic": clinic})
}

// checkOwnerRequire verifies the authenticated user is the owner of the clinic; if not, responds with 403 and returns false
func (h *ClinicHandler) checkOwnerRequire(c *gin.Context, clinicID uuid.UUID) bool {
	userID, ok := GetAuthUserID(c)
	if !ok {
		return false
	}
	role, err := h.userClinicUC.UserRoleInClinic(c.Request.Context(), userID, clinicID)
	if err != nil || role == "" {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied: owner access required"})
		return false
	}
	if !strings.EqualFold(role, util.RoleOwner) {
		c.JSON(http.StatusForbidden, gin.H{"error": "access denied: only the clinic owner can perform this action"})
		return false
	}
	return true
}

// GetClinic retrieves a clinic by ID (requires user to be associated with the clinic)
// GET /api/v1/clinic/:id
// @Summary Retrieve a clinic by ID
// @Description Retrieve a clinic by ID. User must be associated with the clinic.
// @Tags Clinic
// @Accept json
// @Produce json
// @Param id path string true "Clinic ID"
// @Success 200 {object} domain.H
// @Failure 400 {object} domain.H
// @Failure 403 {object} domain.H
// @Failure 404 {object} domain.H
// @Failure 500 {object} domain.H
// @Router /clinic/{id} [get]
func (h *ClinicHandler) GetClinic(c *gin.Context) {
	id := c.Param("id")
	idUUID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid clinic ID"})
		return
	}
	if !RequireClinicAccess(c, h.userClinicUC, idUUID) {
		return
	}
	clinic, err := h.clinicUC.GetClinicByID(c.Request.Context(), idUUID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	// Include current user's role so frontend can show owner-only actions
	userID, _ := GetAuthUserID(c)
	role, _ := h.userClinicUC.UserRoleInClinic(c.Request.Context(), userID, idUUID)
	roleLower := strings.ToLower(role)
	c.JSON(http.StatusOK, gin.H{"clinic": clinic, "currentUserRole": roleLower})
}

// UpdateClinic updates a clinic by ID (requires user to be associated with the clinic)
// PUT /api/v1/clinic/:id
// @Summary Update a clinic by ID
// @Description Update a clinic by ID. User must be associated with the clinic.
// @Tags Clinic
// @Accept json
// @Produce json
// @Param id path string true "Clinic ID"
// @Param clinic body domain.Clinic true "Clinic information"
// @Success 200 {object} domain.H
// @Failure 400 {object} domain.H
// @Failure 403 {object} domain.H
// @Failure 404 {object} domain.H
// @Failure 500 {object} domain.H
// @Router /clinic/{id} [put]
func (h *ClinicHandler) UpdateClinic(c *gin.Context) {
	id := c.Param("id")
	idUUID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid clinic ID"})
		return
	}
	if !RequireClinicAccess(c, h.userClinicUC, idUUID) {
		return
	}
	if !h.checkOwnerRequire(c, idUUID) {
		return
	}
	var req domain.UpdateClinicRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	clinic, err := h.clinicUC.UpdateClinicPartial(c.Request.Context(), idUUID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, clinic)
}

// DeleteClinic deletes a clinic by ID (requires user to be associated with the clinic)
// DELETE /api/v1/clinic/:id
// @Summary Delete a clinic by ID
// @Description Delete a clinic by ID. User must be associated with the clinic.
// @Tags Clinic
// @Accept json
// @Produce json
// @Param id path string true "Clinic ID"
// @Success 200 {object} domain.H
// @Failure 400 {object} domain.H
// @Failure 403 {object} domain.H
// @Failure 404 {object} domain.H
// @Failure 500 {object} domain.H
// @Router /clinic/{id} [delete]
func (h *ClinicHandler) DeleteClinic(c *gin.Context) {
	id := c.Param("id")
	idUUID, err := uuid.Parse(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid clinic ID"})
		return
	}
	if !RequireClinicAccess(c, h.userClinicUC, idUUID) {
		return
	}
	if !h.checkOwnerRequire(c, idUUID) {
		return
	}
	err = h.clinicUC.DeleteClinic(c.Request.Context(), idUUID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "clinic deleted successfully", "clinic_id": idUUID})
}

// GetAllClinics retrieves clinics the current user has access to
// GET /api/v1/clinic
// @Summary Retrieve user's clinics
// @Description Retrieve all clinics the authenticated user is associated with
// @Tags Clinic
// @Accept json
// @Produce json
// @Success 200 {object} domain.H
// @Failure 401 {object} domain.H
// @Failure 500 {object} domain.H
// @Router /clinic [get]
func (h *ClinicHandler) GetAllClinics(c *gin.Context) {
	userID, ok := GetAuthUserID(c)
	if !ok {
		return
	}
	userClinics, err := h.userClinicUC.GetUserClinics(c.Request.Context(), userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	clinics := make([]domain.Clinic, 0, len(userClinics))
	for _, uc := range userClinics {
		clinics = append(clinics, uc.Clinic)
	}
	c.JSON(http.StatusOK, gin.H{"message": "clinics retrieved successfully", "clinics": clinics})
}

// GetClinicByABNNumber retrieves a clinic by ABN number (requires user to be associated)
// GET /api/v1/clinic/abn/:abnNumber
// @Summary Retrieve a clinic by ABN number
// @Description Retrieve a clinic by ABN number. User must be associated with the clinic.
// @Tags Clinic
// @Accept json
// @Produce json
// @Param abnNumber path string true "ABN number"
// @Success 200 {object} domain.H
// @Failure 400 {object} domain.H
// @Failure 403 {object} domain.H
// @Failure 404 {object} domain.H
// @Failure 500 {object} domain.H
// @Router /clinic/abn/{abnNumber} [get]
func (h *ClinicHandler) GetClinicByABNNumber(c *gin.Context) {
	abnNumber := c.Param("abnNumber")
	clinic, err := h.clinicUC.GetClinicByABNNumber(c.Request.Context(), abnNumber)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	if !RequireClinicAccess(c, h.userClinicUC, clinic.ID) {
		return
	}
	c.JSON(http.StatusOK, clinic)
}

// ListClinicAOCs returns all AOC (chart of accounts) associations for a clinic
// GET /api/v1/clinic/:id/aoc
func (h *ClinicHandler) ListClinicAOCs(c *gin.Context) {
	clinicIDStr := c.Param("id")
	clinicID, err := uuid.Parse(clinicIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid clinic ID"})
		return
	}
	if !RequireClinicAccess(c, h.userClinicUC, clinicID) {
		return
	}
	list, err := h.clinicCOAUC.GetClinicAOCs(c.Request.Context(), clinicID)
	if err != nil {
		if errors.Is(err, usecase.ErrClinicNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": list})
}

// AddClinicAOC associates an AOC with a clinic
// POST /api/v1/clinic/:id/aoc
func (h *ClinicHandler) AddClinicAOC(c *gin.Context) {
	clinicIDStr := c.Param("id")
	clinicID, err := uuid.Parse(clinicIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid clinic ID"})
		return
	}
	if !RequireClinicAccess(c, h.userClinicUC, clinicID) {
		return
	}
	var req domain.ClinicCOARequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request: coaId required"})
		return
	}
	resp, err := h.clinicCOAUC.AddClinicAOC(c.Request.Context(), clinicID, req.COAID)
	if err != nil {
		if errors.Is(err, usecase.ErrClinicCOAExists) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "clinic AOC added successfully", "data": resp})
}

// CreateCOAForClinic creates a new chart-of-accounts entry and assigns it to the clinic.
// POST /api/v1/clinic/:id/coa
func (h *ClinicHandler) CreateCOAForClinic(c *gin.Context) {
	clinicIDStr := c.Param("id")
	clinicID, err := uuid.Parse(clinicIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid clinic ID"})
		return
	}
	if !RequireClinicAccess(c, h.userClinicUC, clinicID) {
		return
	}
	var req domain.CreateCOAForClinicRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request: code, name, accountTypeId, accountTaxId required"})
		return
	}
	resp, err := h.clinicCOAUC.CreateCOAForClinic(c.Request.Context(), clinicID, &req)
	if err != nil {
		if errors.Is(err, usecase.ErrClinicCOAExists) {
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": "chart of accounts entry created and assigned to clinic", "data": resp})
}

// RemoveClinicAOC removes an AOC association from a clinic
// DELETE /api/v1/clinic/:id/aoc/:associationId
func (h *ClinicHandler) RemoveClinicAOC(c *gin.Context) {
	clinicIDStr := c.Param("id")
	clinicID, err := uuid.Parse(clinicIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid clinic ID"})
		return
	}
	if !RequireClinicAccess(c, h.userClinicUC, clinicID) {
		return
	}
	associationIDStr := c.Param("associationId")
	associationID, err := uuid.Parse(associationIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid association ID"})
		return
	}
	existing, err := h.clinicCOAUC.GetClinicAOCByID(c.Request.Context(), associationID)
	if err != nil || existing == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "clinic AOC association not found"})
		return
	}
	if existing.ClinicID != clinicID {
		c.JSON(http.StatusForbidden, gin.H{"error": "association does not belong to this clinic"})
		return
	}
	if err := h.clinicCOAUC.RemoveClinicAOC(c.Request.Context(), associationID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "clinic AOC removed successfully"})
}
