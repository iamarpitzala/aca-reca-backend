package http

import (
	"errors"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/iamarpitzala/aca-reca-backend/internal/application/usecase"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain/clinic"
	"github.com/iamarpitzala/aca-reca-backend/util"
	utils "github.com/iamarpitzala/aca-reca-backend/util"
)

type ClinicHandler struct {
	clinicUC     *usecase.ClinicService
	userClinicUC *usecase.UserClinicService
	settingsUC   *usecase.ClinicFinancialSettingsService
}

func NewClinicHandler(clinicUC *usecase.ClinicService, userClinicUC *usecase.UserClinicService, settingsUC *usecase.ClinicFinancialSettingsService) *ClinicHandler {
	return &ClinicHandler{
		clinicUC:     clinicUC,
		userClinicUC: userClinicUC,
		settingsUC:   settingsUC,
	}
}

// CreateClinic creates a new clinic and associates the creating user as owner
// POST /api/v1/clinic
// @Summary Create a new clinic
// @Description Create a new clinic with the given information. The creating user is automatically associated as owner.
// @Tags Clinic
// @Accept json
// @Produce json
// @Param clinic body clinic.Clinic true "Clinic information"
// @Success 201 {object} domain.H
// @Failure 400 {object} domain.H
// @Failure 500 {object} domain.H
// @Router /clinic [post]
func (h *ClinicHandler) CreateClinic(c *gin.Context) {
	userID, ok := GetAuthUserID(c)
	if !ok {
		return
	}
	var clinic clinic.Clinic
	if err := c.ShouldBindJSON(&clinic); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	clinic.UserID = userID
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
			"error":  utils.ErrClinicCreatedButLinkFailed,
			"detail": err.Error(),
		})
		return
	}
	if err := h.settingsUC.CreateDefaultFinancialYearForClinic(c.Request.Context(), clinic.ID); err != nil {
		_ = h.clinicUC.DeleteClinic(c.Request.Context(), clinic.ID)
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":  "clinic created but financial year setup failed",
			"detail": err.Error(),
		})
		return
	}
	c.JSON(http.StatusCreated, gin.H{"message": utils.MsgClinicCreatedSuccessfully, "clinic_id": clinic.ID, "clinic": clinic})
}

// checkOwnerRequire verifies the authenticated user is the owner of the clinic; if not, responds with 403 and returns false
func (h *ClinicHandler) checkOwnerRequire(c *gin.Context, clinicID string) bool {
	userID, ok := GetAuthUserID(c)
	if !ok {
		return false
	}
	role, err := h.userClinicUC.UserRoleInClinic(c.Request.Context(), userID, clinicID)
	if err != nil || role == "" {
		c.JSON(http.StatusForbidden, gin.H{"error": utils.ErrAccessDeniedOwnerRequired})
		return false
	}
	if !strings.EqualFold(role, util.RoleOwner) {
		c.JSON(http.StatusForbidden, gin.H{"error": utils.ErrAccessDeniedOwnerOnly})
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
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "clinic id required"})
		return
	}

	clinic, err := h.clinicUC.GetClinicByID(c.Request.Context(), id)
	if err != nil || clinic == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "clinic not found"})
		return
	}
	if !RequireClinicAccess(c, h.userClinicUC, id) {
		return
	}
	// Include current user's role so frontend can show owner-only actions
	userID, _ := GetAuthUserID(c)
	role, _ := h.userClinicUC.UserRoleInClinic(c.Request.Context(), userID, id)
	roleLower := strings.ToLower(role)
	c.JSON(http.StatusOK, gin.H{"clinic": clinic, "currentUserRole": roleLower})
}

// ListClinicCOA returns the chart of accounts linked to the clinic.
// GET /api/v1/clinic/:id/coa
// Returns 200 with data: [] (empty until clinic-COA association is implemented).
func (h *ClinicHandler) ListClinicCOA(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "clinic id required"})
		return
	}
	if !RequireClinicAccess(c, h.userClinicUC, id) {
		return
	}
	// TODO: when tbl_clinic_coa (or equivalent) exists, return clinic COA associations with details
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "clinic coa retrieved", "data": []interface{}{}})
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
	if !RequireClinicAccess(c, h.userClinicUC, id) {
		return
	}
	if !h.checkOwnerRequire(c, id) {
		return
	}
	var req clinic.UpdateClinicRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	clinic, err := h.clinicUC.UpdateClinicPartial(c.Request.Context(), id, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, clinic)
}

// ActivateClinic sets the clinic as active (owner only).
// PATCH /api/v1/clinic/:id/activate
func (h *ClinicHandler) ActivateClinic(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "clinic id required"})
		return
	}
	cl, err := h.clinicUC.GetClinicByID(c.Request.Context(), id)
	if err != nil || cl == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": utils.ErrClinicNotFound})
		return
	}
	if !RequireClinicAccess(c, h.userClinicUC, id) {
		return
	}
	if !h.checkOwnerRequire(c, id) {
		return
	}
	clinic, err := h.clinicUC.SetActiveStatus(c.Request.Context(), id, true)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "clinic activated", "clinic": clinic})
}

// DeactivateClinic sets the clinic as inactive (owner only).
// PATCH /api/v1/clinic/:id/deactivate
func (h *ClinicHandler) DeactivateClinic(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "clinic id required"})
		return
	}

	cl, err := h.clinicUC.GetClinicByID(c.Request.Context(), id)
	if err != nil || cl == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": utils.ErrClinicNotFound})
		return
	}
	if !RequireClinicAccess(c, h.userClinicUC, id) {
		return
	}
	if !h.checkOwnerRequire(c, id) {
		return
	}
	clinic, err := h.clinicUC.SetActiveStatus(c.Request.Context(), id, false)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "clinic deactivated", "clinic": clinic})
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
	if !RequireClinicAccess(c, h.userClinicUC, id) {
		return
	}
	if !h.checkOwnerRequire(c, id) {
		return
	}
	err := h.clinicUC.DeleteClinic(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": utils.MsgClinicDeletedSuccessfully, "clinic_id": id})
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
	clinics := make([]clinic.Clinic, 0, len(userClinics))
	for _, uc := range userClinics {
		clinics = append(clinics, clinic.Clinic{
			ID:          uc.C_ID,
			Name:        uc.C_Name,
			ABNNumber:   uc.C_ABNNumber,
			Address:     uc.C_Address,
			City:        uc.C_City,
			State:       uc.C_State,
			Postcode:    uc.C_Postcode,
			Phone:       uc.C_Phone,
			Email:       uc.C_Email,
			Website:     uc.C_Website,
			LogoURL:     uc.C_LogoURL,
			Description: uc.C_Description,
			ShareType:   uc.C_ShareType,
			MethodType:  uc.C_MethodType,
			ClinicShare: uc.C_ClinicShare,
			OwnerShare:  uc.C_OwnerShare,
			IsActive:    uc.C_IsActive,
			CreatedAt:   uc.C_CreatedAt,
			UpdatedAt:   uc.C_UpdatedAt,
		})
	}
	if clinics == nil {
		clinics = []clinic.Clinic{}
	}
	c.JSON(http.StatusOK, gin.H{"message": utils.MsgClinicsRetrievedSuccessfully, "clinics": clinics})
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
