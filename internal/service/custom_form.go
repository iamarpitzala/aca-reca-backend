package service

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/calculation"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
	"github.com/iamarpitzala/aca-reca-backend/internal/repository"
	"github.com/jmoiron/sqlx"
)

type CustomFormService struct {
	db *sqlx.DB
}

func NewCustomFormService(db *sqlx.DB) *CustomFormService {
	return &CustomFormService{db: db}
}

func (s *CustomFormService) Create(ctx context.Context, req *domain.CreateCustomFormRequest, userID uuid.UUID) (*domain.CustomFormResponse, error) {
	clinicID, err := uuid.Parse(req.ClinicID)
	if err != nil {
		return nil, errors.New("invalid clinic ID")
	}
	if _, err := repository.GetClinicByID(ctx, s.db, clinicID); err != nil {
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
		Version:                      1,
		CreatedBy:                    userID,
		CreatedAt:                    now,
		UpdatedAt:                    now,
	}
	if err := repository.CreateCustomForm(ctx, s.db, form); err != nil {
		return nil, err
	}
	return customFormToResponse(form), nil
}

func (s *CustomFormService) GetByID(ctx context.Context, id uuid.UUID) (*domain.CustomFormResponse, error) {
	form, err := repository.GetCustomFormByID(ctx, s.db, id)
	if err != nil {
		return nil, err
	}
	return customFormToResponse(form), nil
}

func (s *CustomFormService) GetByClinicID(ctx context.Context, clinicID uuid.UUID) ([]domain.CustomFormResponse, error) {
	forms, err := repository.GetCustomFormsByClinicID(ctx, s.db, clinicID)
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
	forms, err := repository.GetPublishedCustomFormsByClinicID(ctx, s.db, clinicID)
	if err != nil {
		return nil, err
	}
	out := make([]domain.CustomFormResponse, len(forms))
	for i := range forms {
		out[i] = *customFormToResponse(&forms[i])
	}
	return out, nil
}

func (s *CustomFormService) Update(ctx context.Context, id uuid.UUID, req *domain.UpdateCustomFormRequest) (*domain.CustomFormResponse, error) {
	form, err := repository.GetCustomFormByID(ctx, s.db, id)
	if err != nil {
		return nil, err
	}
	if form.Status == "published" {
		// Published: only allow updating field section if needed; for simplicity we allow full update of fields for section only
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
	}
	form.UpdatedAt = time.Now()
	if err := repository.UpdateCustomForm(ctx, s.db, form); err != nil {
		return nil, err
	}
	return customFormToResponse(form), nil
}

func (s *CustomFormService) Publish(ctx context.Context, id uuid.UUID) (*domain.CustomFormResponse, error) {
	form, err := repository.GetCustomFormByID(ctx, s.db, id)
	if err != nil {
		return nil, err
	}
	if form.Status == "published" {
		return nil, errors.New("form is already published")
	}
	if len(form.Fields) == 0 || string(form.Fields) == "[]" {
		return nil, errors.New("cannot publish a form with no fields")
	}
	if err := repository.PublishCustomForm(ctx, s.db, id); err != nil {
		return nil, err
	}
	form.Status = "published"
	now := time.Now()
	form.PublishedAt = &now
	form.UpdatedAt = now
	return customFormToResponse(form), nil
}

func (s *CustomFormService) Archive(ctx context.Context, id uuid.UUID) (*domain.CustomFormResponse, error) {
	form, err := repository.GetCustomFormByID(ctx, s.db, id)
	if err != nil {
		return nil, err
	}
	if err := repository.ArchiveCustomForm(ctx, s.db, id); err != nil {
		return nil, err
	}
	form.Status = "archived"
	form.UpdatedAt = time.Now()
	return customFormToResponse(form), nil
}

func (s *CustomFormService) Delete(ctx context.Context, id uuid.UUID) error {
	return repository.DeleteCustomForm(ctx, s.db, id)
}

func (s *CustomFormService) Duplicate(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*domain.CustomFormResponse, error) {
	form, err := repository.GetCustomFormByID(ctx, s.db, id)
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
		Version:                      1,
		CreatedBy:                    userID,
		CreatedAt:                    now,
		UpdatedAt:                    now,
	}
	if err := repository.CreateCustomForm(ctx, s.db, newForm); err != nil {
		return nil, err
	}
	return customFormToResponse(newForm), nil
}

func customFormToResponse(f *domain.CustomForm) *domain.CustomFormResponse {
	return &domain.CustomFormResponse{
		ID:                           f.ID.String(),
		ClinicID:                     f.ClinicID.String(),
		Name:                         f.Name,
		Description:                  f.Description,
		CalculationMethod:            f.CalculationMethod,
		FormType:                     f.FormType,
		Status:                       f.Status,
		Fields:                       f.Fields,
		DefaultPaymentResponsibility: f.DefaultPaymentResponsibility,
		ServiceFacilityFeePercent:    f.ServiceFacilityFeePercent,
		Version:                      f.Version,
		CreatedBy:                    f.CreatedBy.String(),
		CreatedAt:                    f.CreatedAt,
		UpdatedAt:                    f.UpdatedAt,
		PublishedAt:                  f.PublishedAt,
	}
}

// Entry methods

func (s *CustomFormService) CreateEntry(ctx context.Context, req *domain.CreateEntryRequest, userID uuid.UUID) (*domain.CustomFormEntryResponse, error) {
	formID, err := uuid.Parse(req.FormID)
	if err != nil {
		return nil, errors.New("invalid form ID")
	}
	clinicID, err := uuid.Parse(req.ClinicID)
	if err != nil {
		return nil, errors.New("invalid clinic ID")
	}
	form, err := repository.GetCustomFormByID(ctx, s.db, formID)
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

	// Run calculations on backend (single source of truth)
	calculations, err := calculation.RunEntryCalculation(
		form.Fields,
		form.FormType,
		form.ServiceFacilityFeePercent,
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
	if err := repository.CreateCustomFormEntry(ctx, s.db, entry); err != nil {
		return nil, err
	}
	return customFormEntryToResponse(entry), nil
}

func (s *CustomFormService) GetEntryByID(ctx context.Context, id uuid.UUID) (*domain.CustomFormEntryResponse, error) {
	entry, err := repository.GetCustomFormEntryByID(ctx, s.db, id)
	if err != nil {
		return nil, err
	}
	return customFormEntryToResponse(entry), nil
}

func (s *CustomFormService) GetEntriesByFormID(ctx context.Context, formID uuid.UUID) ([]domain.CustomFormEntryResponse, error) {
	entries, err := repository.GetCustomFormEntriesByFormID(ctx, s.db, formID)
	if err != nil {
		return nil, err
	}
	out := make([]domain.CustomFormEntryResponse, len(entries))
	for i := range entries {
		out[i] = *customFormEntryToResponse(&entries[i])
	}
	return out, nil
}

func (s *CustomFormService) GetEntriesByClinicID(ctx context.Context, clinicID uuid.UUID) ([]domain.CustomFormEntryResponse, error) {
	entries, err := repository.GetCustomFormEntriesByClinicID(ctx, s.db, clinicID)
	if err != nil {
		return nil, err
	}
	out := make([]domain.CustomFormEntryResponse, len(entries))
	for i := range entries {
		out[i] = *customFormEntryToResponse(&entries[i])
	}
	return out, nil
}

func (s *CustomFormService) UpdateEntry(ctx context.Context, id uuid.UUID, req *domain.UpdateEntryRequest) (*domain.CustomFormEntryResponse, error) {
	entry, err := repository.GetCustomFormEntryByID(ctx, s.db, id)
	if err != nil {
		return nil, err
	}
	form, err := repository.GetCustomFormByID(ctx, s.db, entry.FormID)
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
	calculations, err := calculation.RunEntryCalculation(
		form.Fields,
		form.FormType,
		form.ServiceFacilityFeePercent,
		req.Values,
		deductionsForCalc,
	)
	if err != nil {
		return nil, err
	}
	entry.Calculations = calculations
	entry.UpdatedAt = time.Now()
	if err := repository.UpdateCustomFormEntry(ctx, s.db, entry); err != nil {
		return nil, err
	}
	return customFormEntryToResponse(entry), nil
}

func (s *CustomFormService) DeleteEntry(ctx context.Context, id uuid.UUID) error {
	return repository.DeleteCustomFormEntry(ctx, s.db, id)
}

// PreviewCalculations returns calculations for the given form and values (for live display; no save).
func (s *CustomFormService) PreviewCalculations(ctx context.Context, formID uuid.UUID, valuesJSON, deductionsJSON []byte) ([]byte, error) {
	form, err := repository.GetCustomFormByID(ctx, s.db, formID)
	if err != nil {
		return nil, err
	}
	if len(valuesJSON) == 0 {
		valuesJSON = []byte("[]")
	}
	return calculation.RunEntryCalculation(
		form.Fields,
		form.FormType,
		form.ServiceFacilityFeePercent,
		valuesJSON,
		deductionsJSON,
	)
}

// mergeEntryPaymentResponsibilityIntoDeductions merges the entry's payment responsibility into deductions JSON
// so RunEntryCalculation uses it for Additional Reductions (e.g. Lab Fee GST when "Pay by clinic").
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
