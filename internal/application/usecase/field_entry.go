package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/application/port"
	"github.com/iamarpitzala/aca-reca-backend/internal/calculation"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
)

type FieldEntryService struct {
	repo                    port.FieldEntryRepository
	formRepo                port.CustomFormRepository
	fieldRepo               port.CustomFormFieldRepository
	clinicRepo              port.ClinicRepository
	calcRepo                port.CustomFormCalculationRepository
	netDetailsRepo          port.EntryNetDetailsRepository
	grossDetailsRepo        port.EntryGrossDetailsRepository
	grossReductionRepo       port.EntryGrossReductionRepository
	grossReimbursementRepo   port.EntryGrossReimbursementRepository
	financialSettingsRepo   port.ClinicFinancialSettingsRepository
	calcEngine              port.EntryCalculationEngine
}

func NewFieldEntryService(
	repo port.FieldEntryRepository,
	formRepo port.CustomFormRepository,
	fieldRepo port.CustomFormFieldRepository,
	clinicRepo port.ClinicRepository,
	calcRepo port.CustomFormCalculationRepository,
	netDetailsRepo port.EntryNetDetailsRepository,
	grossDetailsRepo port.EntryGrossDetailsRepository,
	grossReductionRepo port.EntryGrossReductionRepository,
	grossReimbursementRepo port.EntryGrossReimbursementRepository,
	financialSettingsRepo port.ClinicFinancialSettingsRepository,
	calcEngine port.EntryCalculationEngine,
) *FieldEntryService {
	return &FieldEntryService{
		repo:                    repo,
		formRepo:                formRepo,
		fieldRepo:               fieldRepo,
		clinicRepo:              clinicRepo,
		calcRepo:                calcRepo,
		netDetailsRepo:          netDetailsRepo,
		grossDetailsRepo:        grossDetailsRepo,
		grossReductionRepo:       grossReductionRepo,
		grossReimbursementRepo:   grossReimbursementRepo,
		financialSettingsRepo:   financialSettingsRepo,
		calcEngine:              calcEngine,
	}
}

func (s *FieldEntryService) Create(ctx context.Context, req *domain.FieldEntryRequest, userID uuid.UUID) (*domain.FieldEntryResponse, error) {
	formID, err := uuid.Parse(req.FormID)
	if err != nil {
		return nil, errors.New("invalid form ID")
	}

	// Verify form exists and get clinic ID for access control
	form, err := s.formRepo.GetByID(ctx, formID)
	if err != nil {
		return nil, errors.New("form not found")
	}

	// Verify clinic exists
	if _, err := s.clinicRepo.GetByID(ctx, form.ClinicID); err != nil {
		return nil, errors.New("clinic not found")
	}

	entry, err := req.ToDBModel(userID)
	if err != nil {
		return nil, err
	}

	// Get the field to retrieve its form_version_id
	customFormFieldID, err := uuid.Parse(req.CustomFormFieldID)
	if err != nil {
		return nil, errors.New("invalid field ID")
	}
	field, err := s.fieldRepo.GetByID(ctx, customFormFieldID)
	if err != nil {
		return nil, errors.New("field not found")
	}
	entry.FormVersionID = field.FormVersionID

	if err := s.repo.Create(ctx, entry); err != nil {
		return nil, err
	}

	return entry.ToResponse(), nil
}

// CreateEntry handles creating a full entry with multiple field values
func (s *FieldEntryService) CreateEntry(ctx context.Context, req *domain.CreateEntryRequest, userID uuid.UUID) (*domain.EntryResponse, error) {
	formID, err := uuid.Parse(req.FormID)
	if err != nil {
		return nil, errors.New("invalid form ID")
	}

	// Step 1: Verify form exists and get clinic ID
	form, err := s.formRepo.GetByID(ctx, formID)
	if err != nil {
		return nil, errors.New("form not found")
	}

	// Step 2: Verify clinic exists
	if _, err := s.clinicRepo.GetByID(ctx, form.ClinicID); err != nil {
		return nil, errors.New("clinic not found")
	}

	// Step 3: Parse field IDs and validate
	fieldIDs := make([]uuid.UUID, 0, len(req.Values))
	fieldIDMap := make(map[string]int) // Map fieldID string to index in req.Values
	for i, val := range req.Values {
		fieldID, err := uuid.Parse(val.FieldID)
		if err != nil {
			return nil, errors.New("invalid field ID: " + val.FieldID)
		}
		fieldIDs = append(fieldIDs, fieldID)
		fieldIDMap[val.FieldID] = i
	}

	// Step 4: Get form version ID from first field (assuming all fields belong to same version)
	// Fetch first field to get formVersionID
	firstField, err := s.fieldRepo.GetByID(ctx, fieldIDs[0])
	if err != nil {
		return nil, errors.New("field not found: " + fieldIDs[0].String())
	}
	formVersionID := firstField.FormVersionID

	// Step 5: Fetch ALL fields for this form version ONCE (optimization)
	fields, err := s.fieldRepo.GetByFormVersionID(ctx, formVersionID)
	if err != nil {
		return nil, errors.New("failed to fetch form fields")
	}

	// Build field map by ID for quick lookups
	fieldMap := make(map[uuid.UUID]domain.CustomFormField)
	for _, field := range fields {
		fieldMap[field.ID] = field
	}

	// Step 6: Create field entries using field map
	entryID := uuid.New()
	now := time.Now()
	fieldEntries := make([]*domain.FieldEntry, 0, len(req.Values))

	for _, fieldID := range fieldIDs {
		field, ok := fieldMap[fieldID]
		if !ok {
			return nil, errors.New("field not found: " + fieldID.String())
		}

		fieldEntry := &domain.FieldEntry{
			ID:                uuid.New(),
			FormID:            formID,
			FormVersionID:     field.FormVersionID,
			CustomFormFieldID: fieldID,
			Value:             req.Values[fieldIDMap[fieldID.String()]].Value,
			CreatedBy:         userID,
			CreatedAt:         now,
			UpdatedAt:         now,
			DeletedAt:         nil,
		}
		fieldEntries = append(fieldEntries, fieldEntry)
	}

	// Step 7: Save all field entries
	for _, entry := range fieldEntries {
		if err := s.repo.Create(ctx, entry); err != nil {
			return nil, err
		}
	}

	// Step 8: Build field value responses using field map
	fieldValueResponses := make([]domain.EntryFieldValueResponse, len(fieldEntries))
	for i, entry := range fieldEntries {
		field := fieldMap[entry.CustomFormFieldID]
		valIndex := fieldIDMap[entry.CustomFormFieldID.String()]
		var manualGSTAmount *float64
		if valIndex < len(req.Values) {
			manualGSTAmount = req.Values[valIndex].ManualGSTAmount
		}

		fieldValueResponses[i] = domain.EntryFieldValueResponse{
			FieldID:         entry.CustomFormFieldID.String(),
			FieldName:       field.Label, // Use field label from DB
			Value:           entry.Value,
			BaseAmount:      nil, // Will be populated by calculation
			GSTAmount:       nil, // Will be populated by calculation
			TotalAmount:     nil, // Will be populated by calculation
			ManualGSTAmount: manualGSTAmount,
		}
	}

	// Step 9: Get form calculation settings ONCE (for both calculation and method-specific logic)
	formCalc, err := s.calcRepo.GetByFormVersionID(ctx, formVersionID)
	if err != nil && err.Error() != "calculation settings not found" {
		// Log error but continue (formCalc can be nil)
		_ = err
	}

	// Step 10: Determine calculation method (form version takes precedence over form level)
	calculationMethod := form.CalculationMethod
	if formCalc != nil && formCalc.CalculationMethod != "" {
		calculationMethod = formCalc.CalculationMethod
	}

	// Step 11: Calculate totals using calculation engine (with proper form calculation settings)
	calculationsJSON, err := s.calculateEntryTotals(ctx, form, formVersionID, fieldValueResponses, req.Deductions, formCalc)
	if err != nil {
		// If calculation fails, use empty calculations but log the error
		calculationsJSON = []byte(`{"fieldTotals":[],"totalBaseAmount":0,"totalGSTAmount":0,"totalAmount":0,"netPayable":0,"netReceivable":0,"basMapping":{"gstOnSales1A":0,"gstCredit1B":0,"totalSalesG1":0,"expensesG11":0}}`)
	} else {
		// Parse calculations to update field value responses with calculated amounts
		var calcResult struct {
			FieldTotals []struct {
				FieldID     string  `json:"fieldId"`
				BaseAmount  float64 `json:"baseAmount"`
				GstAmount   float64 `json:"gstAmount"`
				TotalAmount float64 `json:"totalAmount"`
			} `json:"fieldTotals"`
		}
		if err := json.Unmarshal(calculationsJSON, &calcResult); err == nil {
			// Map calculated amounts to field values
			calcMap := make(map[string]*struct {
				BaseAmount  float64
				GstAmount   float64
				TotalAmount float64
			})
			for i := range calcResult.FieldTotals {
				ft := &calcResult.FieldTotals[i]
				calcMap[ft.FieldID] = &struct {
					BaseAmount  float64
					GstAmount   float64
					TotalAmount float64
				}{
					BaseAmount:  ft.BaseAmount,
					GstAmount:   ft.GstAmount,
					TotalAmount: ft.TotalAmount,
				}
			}
			for i := range fieldValueResponses {
				if calc, ok := calcMap[fieldValueResponses[i].FieldID]; ok {
					fieldValueResponses[i].BaseAmount = &calc.BaseAmount
					fieldValueResponses[i].GSTAmount = &calc.GstAmount
					fieldValueResponses[i].TotalAmount = &calc.TotalAmount
				}
			}
		}
	}

	// Step 12: Get GST settings ONCE (will be reused in method-specific calculations)
	gstRate, gstType := s.getGSTSettings(ctx, form.ClinicID)
	gstSettings := struct {
		Rate float64
		Type string
	}{Rate: gstRate, Type: gstType}

	// Step 13: Calculate and store method-specific details
	if calculationMethod == "NET" {
		if err := s.calculateAndStoreNetDetails(ctx, form.ClinicID, fieldEntries[0].ID, fieldValueResponses, req.Deductions, gstSettings, now); err != nil {
			// Log error but don't fail entry creation
			_ = err
		}
	} else if calculationMethod == "GROSS" {
		// Pass fields, gstSettings, and formCalc to avoid re-fetching
		if err := s.calculateAndStoreGrossDetails(ctx, form.ClinicID, formVersionID, fieldEntries[0].ID, calculationsJSON, fieldValueResponses, req.Deductions, fields, gstSettings, formCalc, now); err != nil {
			// Log error but don't fail entry creation
			_ = err
		}
	}

	return &domain.EntryResponse{
		ID:                    entryID.String(),
		FormID:                req.FormID,
		FormName:              form.Name,
		FormType:              form.FormType,
		ClinicID:              req.ClinicID,
		QuarterID:             req.QuarterID,
		Values:                fieldValueResponses,
		Calculations:          calculationsJSON,
		EntryDate:             req.EntryDate,
		Description:           req.Description,
		Remarks:               req.Remarks,
		PaymentResponsibility: req.PaymentResponsibility,
		Deductions:            req.Deductions,
		CreatedBy:             userID.String(),
		CreatedAt:             now,
		UpdatedAt:             now,
	}, nil
}

func (s *FieldEntryService) GetByID(ctx context.Context, id uuid.UUID) (*domain.EntryResponse, error) {
	entry, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("field entry not found")
	}

	// Get form details
	form, err := s.formRepo.GetByID(ctx, entry.FormID)
	if err != nil {
		return nil, errors.New("form not found")
	}

	// Get all field entries created at the same time for the same form (grouped entry)
	allEntries, err := s.repo.GetByFormID(ctx, entry.FormID)
	if err != nil {
		return nil, err
	}

	// Filter entries created at the same timestamp (within 1 second)
	groupedEntries := make([]domain.FieldEntry, 0)
	targetTime := entry.CreatedAt.Unix()
	for _, e := range allEntries {
		if e.CreatedAt.Unix() == targetTime && e.CreatedBy == entry.CreatedBy {
			groupedEntries = append(groupedEntries, e)
		}
	}

	if len(groupedEntries) == 0 {
		groupedEntries = []domain.FieldEntry{*entry}
	}

	// Get fields for mapping
	fields, _ := s.fieldRepo.GetByFormID(ctx, entry.FormID)
	fieldMap := make(map[uuid.UUID]domain.CustomFormField)
	for _, field := range fields {
		fieldMap[field.ID] = field
	}

	// Build field value responses
	fieldValueResponses := make([]domain.EntryFieldValueResponse, 0, len(groupedEntries))
	for _, e := range groupedEntries {
		field, ok := fieldMap[e.CustomFormFieldID]
		fieldName := ""
		if ok {
			fieldName = field.Label
		}

		fieldValueResponses = append(fieldValueResponses, domain.EntryFieldValueResponse{
			FieldID:     e.CustomFormFieldID.String(),
			FieldName:   fieldName,
			Value:       e.Value,
			BaseAmount:  nil,
			GSTAmount:   nil,
			TotalAmount: nil,
		})
	}

	// Calculate totals
	calculationsJSON, _ := s.calculateEntryTotals(ctx, form, entry.FormVersionID, fieldValueResponses, nil, nil)
	if len(calculationsJSON) == 0 {
		calculationsJSON = []byte(`{"fieldTotals":[],"totalBaseAmount":0,"totalGSTAmount":0,"totalAmount":0,"netPayable":0,"netReceivable":0,"basMapping":{"gstOnSales1A":0,"gstCredit1B":0,"totalSalesG1":0,"expensesG11":0}}`)
	} else {
		// Update field values with calculated amounts
		var calcResult struct {
			FieldTotals []struct {
				FieldID     string  `json:"fieldId"`
				BaseAmount  float64 `json:"baseAmount"`
				GstAmount   float64 `json:"gstAmount"`
				TotalAmount float64 `json:"totalAmount"`
			} `json:"fieldTotals"`
		}
		if err := json.Unmarshal(calculationsJSON, &calcResult); err == nil {
			calcMap := make(map[string]*struct {
				BaseAmount  float64
				GstAmount   float64
				TotalAmount float64
			})
			for i := range calcResult.FieldTotals {
				ft := &calcResult.FieldTotals[i]
				calcMap[ft.FieldID] = &struct {
					BaseAmount  float64
					GstAmount   float64
					TotalAmount float64
				}{BaseAmount: ft.BaseAmount, GstAmount: ft.GstAmount, TotalAmount: ft.TotalAmount}
			}
			for i := range fieldValueResponses {
				if calc, ok := calcMap[fieldValueResponses[i].FieldID]; ok {
					fieldValueResponses[i].BaseAmount = &calc.BaseAmount
					fieldValueResponses[i].GSTAmount = &calc.GstAmount
					fieldValueResponses[i].TotalAmount = &calc.TotalAmount
				}
			}
		}
	}

	entryDate := entry.CreatedAt.Format("2006-01-02")
	if entryDate == "" || entry.CreatedAt.IsZero() {
		entryDate = time.Now().Format("2006-01-02")
	}

	return &domain.EntryResponse{
		ID:                   entry.ID.String(),
		FormID:               entry.FormID.String(),
		FormName:             form.Name,
		FormType:             form.FormType,
		ClinicID:             form.ClinicID.String(),
		QuarterID:            nil,
		Values:               fieldValueResponses,
		Calculations:         calculationsJSON,
		EntryDate:            entryDate,
		Description:          "",
		Remarks:              "",
		PaymentResponsibility: nil,
		Deductions:           nil,
		CreatedBy:            entry.CreatedBy.String(),
		CreatedAt:            entry.CreatedAt,
		UpdatedAt:            entry.UpdatedAt,
	}, nil
}

func (s *FieldEntryService) GetByFormID(ctx context.Context, formID uuid.UUID) ([]domain.EntryResponse, error) {
	// Get form details
	form, err := s.formRepo.GetByID(ctx, formID)
	if err != nil {
		return nil, errors.New("form not found")
	}

	// Get all field entries for this form
	fieldEntries, err := s.repo.GetByFormID(ctx, formID)
	if err != nil {
		return nil, err
	}

	if len(fieldEntries) == 0 {
		return []domain.EntryResponse{}, nil
	}

	// Group field entries by created_at (entries created at the same time belong together)
	// Round to nearest second for grouping
	entryGroups := make(map[int64][]domain.FieldEntry)
	for _, entry := range fieldEntries {
		groupKey := entry.CreatedAt.Unix()
		entryGroups[groupKey] = append(entryGroups[groupKey], entry)
	}

	// Convert grouped entries to EntryResponse
	results := make([]domain.EntryResponse, 0, len(entryGroups))
	for _, entries := range entryGroups {
		if len(entries) == 0 {
			continue
		}

		// Use first entry's metadata for the group
		firstEntry := entries[0]

		// Get fields for this specific form version to ensure correct mapping
		fields, _ := s.fieldRepo.GetByFormVersionID(ctx, firstEntry.FormVersionID)
		fieldMap := make(map[uuid.UUID]domain.CustomFormField)
		for _, field := range fields {
			fieldMap[field.ID] = field
		}

		// Build field value responses
		fieldValueResponses := make([]domain.EntryFieldValueResponse, 0, len(entries))
		for _, entry := range entries {
			field, ok := fieldMap[entry.CustomFormFieldID]
			fieldName := ""
			if ok {
				fieldName = field.Label
			}

			fieldValueResponses = append(fieldValueResponses, domain.EntryFieldValueResponse{
				FieldID:     entry.CustomFormFieldID.String(),
				FieldName:   fieldName,
				Value:       entry.Value,
				BaseAmount:  nil, // TODO: Calculate from GST config
				GSTAmount:   nil, // TODO: Calculate from GST config
				TotalAmount: nil, // TODO: Calculate from GST config
			})
		}

		// Use first entry's ID as the entry ID (or generate a stable ID)
		entryID := firstEntry.ID.String()

		// Calculate totals
		calculationsJSON, calcErr := s.calculateEntryTotals(ctx, form, firstEntry.FormVersionID, fieldValueResponses, nil, nil)
		if calcErr != nil || len(calculationsJSON) == 0 {
			// If calculation fails, use empty calculations
			calculationsJSON = []byte(`{"fieldTotals":[],"totalBaseAmount":0,"totalGSTAmount":0,"totalAmount":0,"netPayable":0,"netReceivable":0,"basMapping":{"gstOnSales1A":0,"gstCredit1B":0,"totalSalesG1":0,"expensesG11":0}}`)
		} else {
			// Update field values with calculated amounts
			var calcResult struct {
				FieldTotals []struct {
					FieldID     string  `json:"fieldId"`
					BaseAmount  float64 `json:"baseAmount"`
					GstAmount   float64 `json:"gstAmount"`
					TotalAmount float64 `json:"totalAmount"`
				} `json:"fieldTotals"`
			}
			if err := json.Unmarshal(calculationsJSON, &calcResult); err == nil {
				calcMap := make(map[string]*struct {
					BaseAmount  float64
					GstAmount   float64
					TotalAmount float64
				})
				for i := range calcResult.FieldTotals {
					ft := &calcResult.FieldTotals[i]
					calcMap[ft.FieldID] = &struct {
						BaseAmount  float64
						GstAmount   float64
						TotalAmount float64
					}{BaseAmount: ft.BaseAmount, GstAmount: ft.GstAmount, TotalAmount: ft.TotalAmount}
				}
				for i := range fieldValueResponses {
					if calc, ok := calcMap[fieldValueResponses[i].FieldID]; ok {
						fieldValueResponses[i].BaseAmount = &calc.BaseAmount
						fieldValueResponses[i].GSTAmount = &calc.GstAmount
						fieldValueResponses[i].TotalAmount = &calc.TotalAmount
					}
				}
			}
		}

		// Format entryDate from created_at (YYYY-MM-DD)
		// Ensure we have a valid date
		entryDate := firstEntry.CreatedAt.Format("2006-01-02")
		if entryDate == "" || firstEntry.CreatedAt.IsZero() {
			entryDate = time.Now().Format("2006-01-02")
		}

		results = append(results, domain.EntryResponse{
			ID:                    entryID,
			FormID:                formID.String(),
			FormName:              form.Name,
			FormType:              form.FormType,
			ClinicID:              form.ClinicID.String(),
			QuarterID:             nil, // TODO: Calculate from entryDate
			Values:                fieldValueResponses,
			Calculations:          calculationsJSON,
			EntryDate:             entryDate,
			Description:           "", // TODO: Store in metadata or separate column
			Remarks:               "", // TODO: Store in metadata or separate column
			PaymentResponsibility: nil,
			Deductions:            nil,
			CreatedBy:             firstEntry.CreatedBy.String(),
			CreatedAt:             firstEntry.CreatedAt,
			UpdatedAt:             firstEntry.UpdatedAt,
		})
	}

	// Sort by created_at descending (newest first)
	for i := 0; i < len(results)-1; i++ {
		for j := i + 1; j < len(results); j++ {
			if results[i].CreatedAt.Before(results[j].CreatedAt) {
				results[i], results[j] = results[j], results[i]
			}
		}
	}

	return results, nil
}

func (s *FieldEntryService) GetByClinicID(ctx context.Context, clinicID uuid.UUID) ([]domain.EntryResponse, error) {
	// Get all field entries for this clinic
	fieldEntries, err := s.repo.GetByClinicID(ctx, clinicID)
	if err != nil {
		return nil, err
	}

	if len(fieldEntries) == 0 {
		return []domain.EntryResponse{}, nil
	}

	// Group by form_id and created_at
	type entryKey struct {
		formID    uuid.UUID
		createdAt int64
	}
	entryGroups := make(map[entryKey][]domain.FieldEntry)
	formIDs := make(map[uuid.UUID]bool)

	for _, entry := range fieldEntries {
		key := entryKey{
			formID:    entry.FormID,
			createdAt: entry.CreatedAt.Unix(),
		}
		entryGroups[key] = append(entryGroups[key], entry)
		formIDs[entry.FormID] = true
	}

	// Get all forms for this clinic
	forms := make(map[uuid.UUID]*domain.CustomForm)
	for formID := range formIDs {
		form, err := s.formRepo.GetByID(ctx, formID)
		if err == nil {
			forms[formID] = form
		}
	}

	// Get all fields for all forms
	fieldsMap := make(map[uuid.UUID]domain.CustomFormField)
	for formID := range formIDs {
		fields, _ := s.fieldRepo.GetByFormID(ctx, formID)
		for _, field := range fields {
			fieldsMap[field.ID] = field
		}
	}

	// Convert grouped entries to EntryResponse
	results := make([]domain.EntryResponse, 0, len(entryGroups))
	for key, entries := range entryGroups {
		if len(entries) == 0 {
			continue
		}

		form, ok := forms[key.formID]
		if !ok {
			continue
		}

		firstEntry := entries[0]

		// Build field value responses
		fieldValueResponses := make([]domain.EntryFieldValueResponse, 0, len(entries))
		for _, entry := range entries {
			field, ok := fieldsMap[entry.CustomFormFieldID]
			fieldName := ""
			if ok {
				fieldName = field.Label
			}

			fieldValueResponses = append(fieldValueResponses, domain.EntryFieldValueResponse{
				FieldID:     entry.CustomFormFieldID.String(),
				FieldName:   fieldName,
				Value:       entry.Value,
				BaseAmount:  nil,
				GSTAmount:   nil,
				TotalAmount: nil,
			})
		}

		entryID := firstEntry.ID.String()
		
		// Calculate totals
		calculationsJSON, _ := s.calculateEntryTotals(ctx, form, firstEntry.FormVersionID, fieldValueResponses, nil, nil)
		if len(calculationsJSON) == 0 {
			calculationsJSON = []byte(`{"fieldTotals":[],"totalBaseAmount":0,"totalGSTAmount":0,"totalAmount":0,"netPayable":0,"netReceivable":0,"basMapping":{"gstOnSales1A":0,"gstCredit1B":0,"totalSalesG1":0,"expensesG11":0}}`)
		} else {
			// Update field values with calculated amounts
			var calcResult struct {
				FieldTotals []struct {
					FieldID     string  `json:"fieldId"`
					BaseAmount  float64 `json:"baseAmount"`
					GstAmount   float64 `json:"gstAmount"`
					TotalAmount float64 `json:"totalAmount"`
				} `json:"fieldTotals"`
			}
			if err := json.Unmarshal(calculationsJSON, &calcResult); err == nil {
				calcMap := make(map[string]*struct {
					BaseAmount  float64
					GstAmount   float64
					TotalAmount float64
				})
				for i := range calcResult.FieldTotals {
					ft := &calcResult.FieldTotals[i]
					calcMap[ft.FieldID] = &struct {
						BaseAmount  float64
						GstAmount   float64
						TotalAmount float64
					}{BaseAmount: ft.BaseAmount, GstAmount: ft.GstAmount, TotalAmount: ft.TotalAmount}
				}
				for i := range fieldValueResponses {
					if calc, ok := calcMap[fieldValueResponses[i].FieldID]; ok {
						fieldValueResponses[i].BaseAmount = &calc.BaseAmount
						fieldValueResponses[i].GSTAmount = &calc.GstAmount
						fieldValueResponses[i].TotalAmount = &calc.TotalAmount
					}
				}
			}
		}

		entryDate := firstEntry.CreatedAt.Format("2006-01-02")
		if entryDate == "" || firstEntry.CreatedAt.IsZero() {
			entryDate = time.Now().Format("2006-01-02")
		}

		results = append(results, domain.EntryResponse{
			ID:                    entryID,
			FormID:                key.formID.String(),
			FormName:              form.Name,
			FormType:              form.FormType,
			ClinicID:              clinicID.String(),
			QuarterID:             nil,
			Values:                fieldValueResponses,
			Calculations:          calculationsJSON,
			EntryDate:             entryDate,
			Description:           "",
			Remarks:               "",
			PaymentResponsibility: nil,
			Deductions:            nil,
			CreatedBy:             firstEntry.CreatedBy.String(),
			CreatedAt:             firstEntry.CreatedAt,
			UpdatedAt:             firstEntry.UpdatedAt,
		})
	}

	// Sort by created_at descending
	for i := 0; i < len(results)-1; i++ {
		for j := i + 1; j < len(results); j++ {
			if results[i].CreatedAt.Before(results[j].CreatedAt) {
				results[i], results[j] = results[j], results[i]
			}
		}
	}

	return results, nil
}

func (s *FieldEntryService) Update(ctx context.Context, id uuid.UUID, value float64) (*domain.FieldEntryResponse, error) {
	entry, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("field entry not found")
	}

	entry.Value = value
	entry.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, entry); err != nil {
		return nil, err
	}

	return entry.ToResponse(), nil
}

func (s *FieldEntryService) Delete(ctx context.Context, id uuid.UUID) error {
	entry, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return errors.New("field entry not found")
	}

	// Verify entry exists before deletion
	_ = entry
	return s.repo.Delete(ctx, id)
}

// GetClinicIDFromEntry retrieves the clinic ID for a field entry via its form
func (s *FieldEntryService) GetClinicIDFromEntry(ctx context.Context, entryID uuid.UUID) (uuid.UUID, error) {
	entry, err := s.repo.GetByID(ctx, entryID)
	if err != nil {
		return uuid.Nil, errors.New("field entry not found")
	}

	form, err := s.formRepo.GetByID(ctx, entry.FormID)
	if err != nil {
		return uuid.Nil, errors.New("form not found")
	}

	return form.ClinicID, nil
}

// GetClinicIDFromForm retrieves the clinic ID for a form
func (s *FieldEntryService) GetClinicIDFromForm(ctx context.Context, formID uuid.UUID) (uuid.UUID, error) {
	form, err := s.formRepo.GetByID(ctx, formID)
	if err != nil {
		return uuid.Nil, errors.New("form not found")
	}
	return form.ClinicID, nil
}

// getGSTSettings fetches GST settings from clinic financial settings
func (s *FieldEntryService) getGSTSettings(ctx context.Context, clinicID uuid.UUID) (gstRate float64, gstType string) {
	gstRate = 10.0 // Default GST rate
	gstType = "exclusive" // Default GST type

	financialSettings, err := s.financialSettingsRepo.GetByClinicID(ctx, clinicID)
	if err == nil && financialSettings != nil {
		gstDefaults, err := financialSettings.GetGSTDefaultsMap()
		if err == nil {
			if rateStr, ok := gstDefaults["defaultGSTRate"]; ok {
				if rate, err := strconv.ParseFloat(rateStr, 64); err == nil {
					gstRate = rate
				}
			}
			if typeStr, ok := gstDefaults["defaultGSTType"]; ok {
				gstType = typeStr
			}
		}
	}
	return gstRate, gstType
}

// calculateEntryTotals computes subtotals, GST, and totals for an entry
func (s *FieldEntryService) calculateEntryTotals(ctx context.Context, form *domain.CustomForm, formVersionID uuid.UUID, values []domain.EntryFieldValueResponse, deductions json.RawMessage, formCalc *domain.CustomFormCalculation) (json.RawMessage, error) {
	// Get form fields for this version
	fields, err := s.fieldRepo.GetByFormVersionID(ctx, formVersionID)
	if err != nil {
		return nil, err
	}

	// Convert CustomFormField to calcField format
	type calcField struct {
		ID             string  `json:"id"`
		Name           string  `json:"name"`
		Type           string  `json:"type"`
		Section        string  `json:"section"`
		IncludeInTotal bool    `json:"includeInTotal"`
		GstConfig      *struct {
			Enabled bool    `json:"enabled"`
			Rate    float64 `json:"rate"`
			Type    string  `json:"type"`
		} `json:"gstConfig"`
		PaymentResp string `json:"paymentResponsibility"`
	}

	type entryValue struct {
		FieldID         string      `json:"fieldId"`
		FieldName       string      `json:"fieldName"`
		Value           interface{} `json:"value"`
		ManualGstAmount *float64   `json:"manualGstAmount"`
	}

	calcFields := make([]calcField, 0, len(fields))
	for _, field := range fields {
		// Parse metadata to get includeInTotal
		var metadata map[string]interface{}
		includeInTotal := true
		if len(field.Metadata) > 0 {
			_ = json.Unmarshal(field.Metadata, &metadata)
			if val, ok := metadata["includeInTotal"].(bool); ok {
				includeInTotal = val
			}
		}

		// Get payment responsibility from metadata
		paymentResp := ""
		if metadata != nil {
			if val, ok := metadata["paymentResponsibility"].(string); ok {
				paymentResp = val
			}
		}

		cf := calcField{
			ID:             field.ID.String(),
			Name:           field.Label,
			Type:           strings.ToLower(field.FieldType),
			Section:        strings.ToLower(field.Section),
			IncludeInTotal: includeInTotal,
			PaymentResp:    paymentResp,
		}

		if field.GSTConfig {
			gstRate := 0.0
			if field.GSTRate != nil {
				gstRate = *field.GSTRate
			}
			gstType := strings.ToLower(field.GSTType)
			if gstType == "" {
				gstType = "exclusive"
			}
			cf.GstConfig = &struct {
				Enabled bool    `json:"enabled"`
				Rate    float64 `json:"rate"`
				Type    string  `json:"type"`
			}{
				Enabled: true,
				Rate:    gstRate,
				Type:    gstType,
			}
		}
		calcFields = append(calcFields, cf)
	}

	// Convert values to entryValue format
	entryValues := make([]entryValue, 0, len(values))
	for _, val := range values {
		ev := entryValue{
			FieldID:   val.FieldID,
			FieldName: val.FieldName,
			Value:     val.Value,
		}
		if val.ManualGSTAmount != nil {
			ev.ManualGstAmount = val.ManualGSTAmount
		}
		entryValues = append(entryValues, ev)
	}

	// Marshal to JSON
	fieldsJSON, err := json.Marshal(calcFields)
	if err != nil {
		return nil, err
	}
	valuesJSON, err := json.Marshal(entryValues)
	if err != nil {
		return nil, err
	}

	// Debug: Log what we're sending to calculation engine (disabled)
	// fmt.Printf("DEBUG: Fields JSON: %s\n", string(fieldsJSON))
	// fmt.Printf("DEBUG: Values JSON: %s\n", string(valuesJSON))

	// Get form calculation settings (service facility fee, outwork, etc.)
	var serviceFacilityFeePercent *float64
	var outworkEnabled bool
	var outworkRatePercent *float64

	// Use form calculation settings if available
	if formCalc != nil {
		serviceFacilityFeePercent = formCalc.ServiceFacilityFeePercent
		outworkEnabled = formCalc.OutworkEnabled
		outworkRatePercent = formCalc.OutworkRatePercent
	}

	// Call calculation engine
	calculationsJSON, err := s.calcEngine.RunEntryCalculation(
		fieldsJSON,
		form.FormType,
		form.CalculationMethod,
		serviceFacilityFeePercent,
		outworkEnabled,
		outworkRatePercent,
		valuesJSON,
		deductions,
	)
	if err != nil {
		return nil, err
	}

	return calculationsJSON, nil
}

// calculateAndStoreNetDetails calculates NET method details and stores them in tbl_entry_net_details
func (s *FieldEntryService) calculateAndStoreNetDetails(
	ctx context.Context,
	clinicID uuid.UUID,
	entryID uuid.UUID,
	fieldValueResponses []domain.EntryFieldValueResponse,
	deductionsJSON json.RawMessage,
	gstSettings struct {
		Rate float64
		Type string
	},
	now time.Time,
) error {
	// Step 1: Parse deductions to get commission percent and super settings
	commissionPercent, superHoldingEnabled, superComponentPercent, err := calculation.ParseNetDeductions(deductionsJSON)
	if err != nil {
		return err
	}

	// If commission percent is 0, skip net details calculation
	if commissionPercent == 0 {
		return nil
	}

	// Step 2: Calculate total payment received from field values
	// Sum all field values that are included in total
	totalPaymentReceived := 0.0
	for _, val := range fieldValueResponses {
		if val.TotalAmount != nil {
			totalPaymentReceived += *val.TotalAmount
		} else {
			totalPaymentReceived += val.Value
		}
	}

	// Step 3: Use GST settings passed as parameter (already fetched)
	gstRate := gstSettings.Rate
	gstType := gstSettings.Type

	// Step 4: Run NET calculation
	netInput := calculation.NetCalculationInput{
		TotalPaymentReceived: totalPaymentReceived,
		CommissionPercent:    commissionPercent,
		SuperHoldingEnabled:  superHoldingEnabled,
		SuperComponentPercent: superComponentPercent,
		GSTRate:              gstRate,
		GSTType:              gstType,
	}
	netOutput := calculation.RunNetCalculation(netInput)

	// Step 5: Create and store net details
	netDetails := &domain.EntryNetDetails{
		ID:                      uuid.New(),
		EntryID:                 entryID,
		CommissionPercent:       netOutput.CommissionPercent,
		Commission:              netOutput.Commission,
		GSTOnCommission:        netOutput.GSTOnCommission,
		TotalPaymentReceived:     netOutput.TotalPaymentReceived,
		SuperHoldingEnabled:      netOutput.SuperHoldingEnabled,
		SuperComponentPercent:    netOutput.SuperComponentPercent,
		CommissionComponent:      netOutput.CommissionComponent,
		SuperComponent:           netOutput.SuperComponent,
		TotalForReconciliation:   netOutput.TotalForReconciliation,
		CreatedAt:                now,
		UpdatedAt:                now,
	}

	return s.netDetailsRepo.Create(ctx, netDetails)
}

// calculateAndStoreGrossDetails calculates GROSS method details and stores them (similar to NET method)
func (s *FieldEntryService) calculateAndStoreGrossDetails(
	ctx context.Context,
	clinicID uuid.UUID,
	formVersionID uuid.UUID,
	entryID uuid.UUID,
	calculationsJSON json.RawMessage,
	fieldValueResponses []domain.EntryFieldValueResponse,
	deductionsJSON json.RawMessage,
	fields []domain.CustomFormField, // Passed from CreateEntry to avoid re-fetching
	gstSettings struct {
		Rate float64
		Type string
	}, // Passed from CreateEntry to avoid re-fetching
	formCalc *domain.CustomFormCalculation, // Passed from CreateEntry to avoid re-fetching
	now time.Time,
) error {
	// Step 1: Parse deductions to get service facility fee percent and outwork settings
	serviceFacilityFeePercentFromDeductions, outworkEnabled, _, err := calculation.ParseGrossDeductions(deductionsJSON)
	if err != nil {
		return err
	}

	// Step 2: Determine service facility fee percent (deductions > formCalc > default)
	var serviceFacilityFeePercent float64 = 60.0 // Default 60%
	if serviceFacilityFeePercentFromDeductions != nil {
		serviceFacilityFeePercent = *serviceFacilityFeePercentFromDeductions
	} else if formCalc != nil && formCalc.ServiceFacilityFeePercent != nil {
		serviceFacilityFeePercent = *formCalc.ServiceFacilityFeePercent
	}

	// Use outwork settings from form calculation if not in deductions
	if !outworkEnabled && formCalc != nil {
		outworkEnabled = formCalc.OutworkEnabled
	}

	// Step 3: Use GST settings passed as parameter (already fetched)
	gstRate := gstSettings.Rate

	// Step 5: Calculate net amount directly from field value responses
	// netAmount = netIncome - netExpenses
	// where netIncome and netExpenses are calculated based on GST type (inclusive/exclusive/manual)
	netAmountResult := calculation.CalculateNetAmountFromFieldValues(fieldValueResponses, fields)

	// Step 6: Parse gross calculation output to get field totals for reductions mapping
	var grossCalcOutput struct {
		FieldTotals []struct {
			FieldID     string  `json:"fieldId"`
			BaseAmount  float64 `json:"baseAmount"`
			GstAmount   float64 `json:"gstAmount"`
			TotalAmount float64 `json:"totalAmount"`
		} `json:"fieldTotals"`
	}
	if err := json.Unmarshal(calculationsJSON, &grossCalcOutput); err != nil {
		return err
	}

	// Create field section map for reductions mapping
	fieldSectionMap := make(map[string]string)
	for _, field := range fields {
		fieldSectionMap[field.ID.String()] = strings.ToUpper(field.Section)
	}

	// Build field totals with section info for reductions mapping
	fieldTotals := make([]calculation.FieldTotal, 0, len(grossCalcOutput.FieldTotals))
	for _, ft := range grossCalcOutput.FieldTotals {
		section := "INCOME" // Default
		if s, ok := fieldSectionMap[ft.FieldID]; ok {
			section = s
		}
		fieldTotals = append(fieldTotals, calculation.FieldTotal{
			FieldID:     ft.FieldID,
			Section:     section,
			BaseAmount:  ft.BaseAmount,
			GstAmount:   ft.GstAmount,
			TotalAmount: ft.TotalAmount,
		})
	}

	// Step 7: Run structured GROSS calculation
	// Use calculated netAmount and incomeExclGST from field values
	grossInput := calculation.GrossCalculationInput{
		IncomeExclGST:             netAmountResult.IncomeExclGST,
		NetAmount:                 netAmountResult.NetAmount,
		ServiceFacilityFeePercent: serviceFacilityFeePercent,
		OutworkEnabled:            outworkEnabled,
		GSTRate:                  gstRate,
		FieldTotals:              fieldTotals,
	}
	grossOutput := calculation.RunGrossCalculationStructured(grossInput)

	// Step 7: Create and store gross details
	grossDetails := &domain.EntryGrossDetails{
		ID:                        uuid.New(),
		EntryID:                   entryID,
		ServiceFacilityFeePercent: grossOutput.ServiceFacilityFeePercent,
		ServiceFeeBase:            grossOutput.ServiceFeeBase,
		GstOnServiceFee:           grossOutput.GstOnServiceFee,
		TotalServiceFee:           grossOutput.TotalServiceFee,
		NetAmount:                 grossOutput.NetAmount,
		CreatedAt:                 now,
		UpdatedAt:                 now,
	}

	// Save gross details
	if err := s.grossDetailsRepo.Create(ctx, grossDetails); err != nil {
		return err
	}

	// Step 8: Map reductions (fields with section = "REDUCTION") - reuse fields already fetched
	reductions := make([]*domain.EntryGrossReduction, 0)
	for _, ft := range fieldTotals {
		if ft.Section == "REDUCTION" {
			fieldID, err := uuid.Parse(ft.FieldID)
			if err != nil {
				continue
			}

			reduction := &domain.EntryGrossReduction{
				ID:             uuid.New(),
				GrossDetailsID: grossDetails.ID,
				EntryID:        entryID,
				FieldID:        fieldID,
				BaseAmount:     ft.BaseAmount,
				GstAmount:      ft.GstAmount,
				TotalAmount:    ft.TotalAmount,
				CreatedAt:      now,
			}
			reductions = append(reductions, reduction)
		}
	}

	// Save reductions if any
	if len(reductions) > 0 {
		if err := s.grossReductionRepo.CreateBatch(ctx, reductions); err != nil {
			return err
		}
	}

	return nil
}
