package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/application/port"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
)

type CustomFormService struct {
	repo       port.CustomFormRepository
	clinicRepo port.ClinicRepository
	calcEngine port.EntryCalculationEngine
}

func NewCustomFormService(repo port.CustomFormRepository, clinicRepo port.ClinicRepository, calcEngine port.EntryCalculationEngine) *CustomFormService {
	return &CustomFormService{repo: repo, clinicRepo: clinicRepo, calcEngine: calcEngine}
}

func (s *CustomFormService) Create(ctx context.Context, req *domain.CreateCustomFormRequest, userID uuid.UUID) (*domain.CustomFormResponse, error) {
	clinicID, err := uuid.Parse(req.ClinicID)
	if err != nil {
		return nil, errors.New("invalid clinic ID")
	}
	if _, err := s.clinicRepo.GetByID(ctx, clinicID); err != nil {
		return nil, errors.New("clinic not found")
	}
	if req.CalculationMethod != "net" && req.CalculationMethod != "gross" {
		return nil, errors.New("calculation method must be net or gross")
	}
	if req.FormType != "income" && req.FormType != "expense" && req.FormType != "both" {
		return nil, errors.New("form type must be income, expense, or both")
	}
	if len(req.Fields) == 0 {
		req.Fields = []byte("[]")
	}
	outworkEnabled := false
	if req.OutworkEnabled != nil {
		outworkEnabled = *req.OutworkEnabled
	}
	now := time.Now()
	form := &domain.CustomForm{
		ID:                           uuid.New(),
		ClinicID:                     clinicID,
		Name:                         req.Name,
		Description:                  req.Description,
		CalculationMethod:            req.CalculationMethod,
		FormType:                     req.FormType,
		Status:                       "draft",
		Fields:                       req.Fields,
		DefaultPaymentResponsibility: req.DefaultPaymentResponsibility,
		ServiceFacilityFeePercent:    req.ServiceFacilityFeePercent,
		OutworkEnabled:               outworkEnabled,
		OutworkRatePercent:           req.OutworkRatePercent,
		Version:                      1,
		CreatedBy:                    userID,
		CreatedAt:                    now,
		UpdatedAt:                    now,
	}
	if err := s.repo.Create(ctx, form); err != nil {
		return nil, err
	}
	return customFormToResponse(form), nil
}

func (s *CustomFormService) GetByID(ctx context.Context, id uuid.UUID) (*domain.CustomFormResponse, error) {
	form, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return customFormToResponse(form), nil
}

func (s *CustomFormService) GetByClinicID(ctx context.Context, clinicID uuid.UUID) ([]domain.CustomFormResponse, error) {
	forms, err := s.repo.GetByClinicID(ctx, clinicID)
	if err != nil {
		return nil, err
	}
	out := make([]domain.CustomFormResponse, len(forms))
	for i := range forms {
		out[i] = *customFormToResponse(&forms[i])
	}
	return out, nil
}

func (s *CustomFormService) GetPublishedByClinicID(ctx context.Context, clinicID uuid.UUID) ([]domain.CustomFormResponse, error) {
	forms, err := s.repo.GetPublishedByClinicID(ctx, clinicID)
	if err != nil {
		return nil, err
	}
	out := make([]domain.CustomFormResponse, len(forms))
	for i := range forms {
		out[i] = *customFormToResponse(&forms[i])
	}
	return out, nil
}

func (s *CustomFormService) Update(ctx context.Context, form *domain.CustomForm) error {
	return s.repo.Update(ctx, form)
}

// UpdateByRequest updates a form from API request (handles draft vs published).
func (s *CustomFormService) UpdateByRequest(ctx context.Context, id uuid.UUID, req *domain.UpdateCustomFormRequest) (*domain.CustomFormResponse, error) {
	form, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if form.Status == "published" {
		if req.Fields != nil {
			form.Fields = req.Fields
		}
	} else {
		if req.Name != nil {
			form.Name = *req.Name
		}
		if req.Description != nil {
			form.Description = *req.Description
		}
		if req.Fields != nil {
			form.Fields = req.Fields
		}
		if req.DefaultPaymentResponsibility != nil {
			form.DefaultPaymentResponsibility = req.DefaultPaymentResponsibility
		}
		if req.ServiceFacilityFeePercent != nil {
			form.ServiceFacilityFeePercent = req.ServiceFacilityFeePercent
		}
		if req.OutworkEnabled != nil {
			form.OutworkEnabled = *req.OutworkEnabled
		}
		if req.OutworkRatePercent != nil {
			form.OutworkRatePercent = req.OutworkRatePercent
		}
	}
	form.UpdatedAt = time.Now()
	if err := s.repo.Update(ctx, form); err != nil {
		return nil, err
	}
	return customFormToResponse(form), nil
}

func (s *CustomFormService) Publish(ctx context.Context, id uuid.UUID) (*domain.CustomFormResponse, error) {
	form, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if form.Status == "published" {
		return nil, errors.New("form is already published")
	}
	if len(form.Fields) == 0 || string(form.Fields) == "[]" {
		return nil, errors.New("cannot publish a form with no fields")
	}
	if err := s.repo.Publish(ctx, id); err != nil {
		return nil, err
	}
	form.Status = "published"
	now := time.Now()
	form.PublishedAt = &now
	form.UpdatedAt = now
	return customFormToResponse(form), nil
}

func (s *CustomFormService) Unpublish(ctx context.Context, id uuid.UUID) (*domain.CustomFormResponse, error) {
	form, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if form.Status != "published" {
		return nil, errors.New("only published forms can be unpublished")
	}
	if err := s.repo.Unpublish(ctx, id); err != nil {
		return nil, err
	}
	form.Status = "draft"
	form.PublishedAt = nil
	form.UpdatedAt = time.Now()
	return customFormToResponse(form), nil
}

func (s *CustomFormService) Archive(ctx context.Context, id uuid.UUID) (*domain.CustomFormResponse, error) {
	form, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := s.repo.Archive(ctx, id); err != nil {
		return nil, err
	}
	form.Status = "archived"
	form.UpdatedAt = time.Now()
	return customFormToResponse(form), nil
}

func (s *CustomFormService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

// Duplicate creates a copy of a form in draft status.
func (s *CustomFormService) Duplicate(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*domain.CustomFormResponse, error) {
	form, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	newForm := &domain.CustomForm{
		ID:                           uuid.New(),
		ClinicID:                     form.ClinicID,
		Name:                         form.Name + " (Copy)",
		Description:                  form.Description,
		CalculationMethod:            form.CalculationMethod,
		FormType:                     form.FormType,
		Status:                       "draft",
		Fields:                       form.Fields,
		DefaultPaymentResponsibility: form.DefaultPaymentResponsibility,
		ServiceFacilityFeePercent:    form.ServiceFacilityFeePercent,
		OutworkEnabled:               form.OutworkEnabled,
		OutworkRatePercent:           form.OutworkRatePercent,
		Version:                      1,
		CreatedBy:                    userID,
		CreatedAt:                    now,
		UpdatedAt:                    now,
	}
	if err := s.repo.Create(ctx, newForm); err != nil {
		return nil, err
	}
	return customFormToResponse(newForm), nil
}

// CreateEntryFromRequest creates an entry from API request, running backend calculations.
func (s *CustomFormService) CreateEntryFromRequest(ctx context.Context, req *domain.CreateEntryRequest, userID uuid.UUID) (*domain.CustomFormEntryResponse, error) {
	formID, err := uuid.Parse(req.FormID)
	if err != nil {
		return nil, errors.New("invalid form ID")
	}
	clinicID, err := uuid.Parse(req.ClinicID)
	if err != nil {
		return nil, errors.New("invalid clinic ID")
	}
	form, err := s.repo.GetByID(ctx, formID)
	if err != nil {
		return nil, err
	}
	if form.ClinicID != clinicID {
		return nil, errors.New("clinic does not own this form")
	}
	if form.Status != "published" {
		return nil, errors.New("cannot create entries for unpublished forms")
	}

	entryDate, err := time.Parse("2006-01-02", req.EntryDate)
	if err != nil {
		entryDate, err = time.Parse(time.RFC3339, req.EntryDate)
		if err != nil {
			return nil, errors.New("entryDate must be ISO date (YYYY-MM-DD) or RFC3339")
		}
	}

	var quarterID *uuid.UUID
	if req.QuarterID != nil && *req.QuarterID != "" {
		q, err := uuid.Parse(*req.QuarterID)
		if err != nil {
			return nil, errors.New("invalid quarter ID")
		}
		quarterID = &q
	}

	if len(req.Values) == 0 {
		req.Values = []byte("[]")
	}
	deductions := req.Deductions
	deductionsForCalc := mergeEntryPaymentResponsibilityIntoDeductions(deductions, req.PaymentResponsibility)
	if len(deductionsForCalc) == 0 {
		deductionsForCalc = nil
	}
	if len(deductions) == 0 {
		deductions = nil
	}

	calculations, err := s.calcEngine.RunEntryCalculation(
		form.Fields,
		form.FormType,
		form.CalculationMethod,
		form.ServiceFacilityFeePercent,
		form.OutworkEnabled,
		form.OutworkRatePercent,
		req.Values,
		deductionsForCalc,
	)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	entry := &domain.CustomFormEntry{
		ID:                    uuid.New(),
		FormID:                formID,
		FormName:              form.Name,
		FormType:              form.FormType,
		ClinicID:              clinicID,
		QuarterID:             quarterID,
		Values:                req.Values,
		Calculations:          calculations,
		EntryDate:             entryDate,
		Description:           req.Description,
		Remarks:               req.Remarks,
		PaymentResponsibility: req.PaymentResponsibility,
		Deductions:            deductions,
		CreatedBy:             userID,
		CreatedAt:             now,
		UpdatedAt:             now,
	}
	if err := s.repo.CreateEntry(ctx, entry); err != nil {
		return nil, err
	}
	return customFormEntryToResponse(entry), nil
}

func mergeEntryPaymentResponsibilityIntoDeductions(deductions json.RawMessage, paymentResponsibility *string) json.RawMessage {
	var m map[string]interface{}
	if len(deductions) > 0 {
		_ = json.Unmarshal(deductions, &m)
	}
	if m == nil {
		m = make(map[string]interface{})
	}
	if paymentResponsibility != nil && *paymentResponsibility != "" {
		m["entryPaymentResponsibility"] = *paymentResponsibility
	}
	out, _ := json.Marshal(m)
	return out
}

func (s *CustomFormService) CreateEntry(ctx context.Context, entry *domain.CustomFormEntry) error {
	return s.repo.CreateEntry(ctx, entry)
}

func (s *CustomFormService) GetEntryByID(ctx context.Context, id uuid.UUID) (*domain.CustomFormEntry, error) {
	return s.repo.GetEntryByID(ctx, id)
}

func (s *CustomFormService) GetEntryResponseByID(ctx context.Context, id uuid.UUID) (*domain.CustomFormEntryResponse, error) {
	entry, err := s.repo.GetEntryByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return customFormEntryToResponse(entry), nil
}

func (s *CustomFormService) GetEntriesByFormID(ctx context.Context, formID uuid.UUID) ([]domain.CustomFormEntry, error) {
	return s.repo.GetEntriesByFormID(ctx, formID)
}

func (s *CustomFormService) GetEntriesByClinicID(ctx context.Context, clinicID uuid.UUID) ([]domain.CustomFormEntry, error) {
	return s.repo.GetEntriesByClinicID(ctx, clinicID)
}

func (s *CustomFormService) GetEntriesByQuarter(ctx context.Context, clinicID, quarterID uuid.UUID) ([]domain.CustomFormEntry, error) {
	return s.repo.GetEntriesByQuarter(ctx, clinicID, quarterID)
}

func (s *CustomFormService) GetEntriesResponseByFormID(ctx context.Context, formID uuid.UUID) ([]domain.CustomFormEntryResponse, error) {
	entries, err := s.repo.GetEntriesByFormID(ctx, formID)
	if err != nil {
		return nil, err
	}
	out := make([]domain.CustomFormEntryResponse, len(entries))
	for i := range entries {
		out[i] = *customFormEntryToResponse(&entries[i])
	}
	return out, nil
}

func (s *CustomFormService) GetEntriesResponseByClinicID(ctx context.Context, clinicID uuid.UUID) ([]domain.CustomFormEntryResponse, error) {
	entries, err := s.repo.GetEntriesByClinicID(ctx, clinicID)
	if err != nil {
		return nil, err
	}
	out := make([]domain.CustomFormEntryResponse, len(entries))
	for i := range entries {
		out[i] = *customFormEntryToResponse(&entries[i])
	}
	return out, nil
}

// UpdateEntryFromRequest updates entry values and recalculates.
func (s *CustomFormService) UpdateEntryFromRequest(ctx context.Context, id uuid.UUID, req *domain.UpdateEntryRequest) (*domain.CustomFormEntryResponse, error) {
	entry, err := s.repo.GetEntryByID(ctx, id)
	if err != nil {
		return nil, err
	}
	form, err := s.repo.GetByID(ctx, entry.FormID)
	if err != nil {
		return nil, err
	}
	entry.Values = req.Values
	deductions := entry.Deductions
	deductionsForCalc := mergeEntryPaymentResponsibilityIntoDeductions(deductions, entry.PaymentResponsibility)
	if len(deductionsForCalc) == 0 {
		deductionsForCalc = nil
	}
	if len(deductions) == 0 {
		deductions = nil
	}
	calculations, err := s.calcEngine.RunEntryCalculation(
		form.Fields,
		form.FormType,
		form.CalculationMethod,
		form.ServiceFacilityFeePercent,
		form.OutworkEnabled,
		form.OutworkRatePercent,
		req.Values,
		deductionsForCalc,
	)
	if err != nil {
		return nil, err
	}
	entry.Calculations = calculations
	entry.UpdatedAt = time.Now()
	if err := s.repo.UpdateEntry(ctx, entry); err != nil {
		return nil, err
	}
	return customFormEntryToResponse(entry), nil
}

func (s *CustomFormService) UpdateEntry(ctx context.Context, entry *domain.CustomFormEntry) error {
	return s.repo.UpdateEntry(ctx, entry)
}

func (s *CustomFormService) DeleteEntry(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteEntry(ctx, id)
}

// PreviewCalculations returns calculations for given form and values (no save).
func (s *CustomFormService) PreviewCalculations(ctx context.Context, formID uuid.UUID, valuesJSON, deductionsJSON []byte) ([]byte, error) {
	form, err := s.repo.GetByID(ctx, formID)
	if err != nil {
		return nil, err
	}
	if len(valuesJSON) == 0 {
		valuesJSON = []byte("[]")
	}
	return s.calcEngine.RunEntryCalculation(
		form.Fields,
		form.FormType,
		form.CalculationMethod,
		form.ServiceFacilityFeePercent,
		form.OutworkEnabled,
		form.OutworkRatePercent,
		valuesJSON,
		deductionsJSON,
	)
}

func customFormToResponse(form *domain.CustomForm) *domain.CustomFormResponse {
	if form == nil {
		return nil
	}
	return &domain.CustomFormResponse{
		ID:                           form.ID.String(),
		ClinicID:                     form.ClinicID.String(),
		Name:                         form.Name,
		Description:                  form.Description,
		CalculationMethod:            form.CalculationMethod,
		FormType:                     form.FormType,
		Status:                       form.Status,
		Fields:                       form.Fields,
		DefaultPaymentResponsibility: form.DefaultPaymentResponsibility,
		ServiceFacilityFeePercent:    form.ServiceFacilityFeePercent,
		OutworkEnabled:               form.OutworkEnabled,
		OutworkRatePercent:           form.OutworkRatePercent,
		Version:                      form.Version,
		PublishedAt:                  form.PublishedAt,
		CreatedBy:                    form.CreatedBy.String(),
		CreatedAt:                    form.CreatedAt,
		UpdatedAt:                    form.UpdatedAt,
	}
}

func customFormEntryToResponse(e *domain.CustomFormEntry) *domain.CustomFormEntryResponse {
	var quarterID *string
	if e.QuarterID != nil {
		s := e.QuarterID.String()
		quarterID = &s
	}
	return &domain.CustomFormEntryResponse{
		ID:                    e.ID.String(),
		FormID:                e.FormID.String(),
		FormName:              e.FormName,
		FormType:              e.FormType,
		ClinicID:              e.ClinicID.String(),
		QuarterID:             quarterID,
		Values:                e.Values,
		Calculations:          e.Calculations,
		EntryDate:             e.EntryDate,
		Description:           e.Description,
		Remarks:               e.Remarks,
		PaymentResponsibility: e.PaymentResponsibility,
		Deductions:            e.Deductions,
		CreatedBy:             e.CreatedBy.String(),
		CreatedAt:             e.CreatedAt,
		UpdatedAt:             e.UpdatedAt,
	}
}
