package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/iamarpitzala/aca-reca-backend/internal/application/port"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain/form"
)

type FieldEntryService struct {
	repo                 port.FieldEntryRepository
	formRepo             port.CustomFormRepository
	fieldRepo            port.CustomFormFieldRepository
	clinicRepo           port.ClinicRepository
	financialSettingRepo port.ClinicFinancialSettingRepository
	transactionRepo      port.TransactionRepository
}

func NewFieldEntryService(
	repo port.FieldEntryRepository,
	formRepo port.CustomFormRepository,
	fieldRepo port.CustomFormFieldRepository,
	clinicRepo port.ClinicRepository,
	financialSettingRepo port.ClinicFinancialSettingRepository,
	transactionRepo port.TransactionRepository,
) *FieldEntryService {
	return &FieldEntryService{
		repo:                 repo,
		formRepo:             formRepo,
		fieldRepo:            fieldRepo,
		clinicRepo:           clinicRepo,
		financialSettingRepo: financialSettingRepo,
		transactionRepo:      transactionRepo,
	}
}

func (s *FieldEntryService) Create(ctx context.Context, req *form.EntryFieldRequest, userID string) (*form.FieldEntryResponse, error) {
	form, err := s.formRepo.GetByID(ctx, req.FormID)
	if err != nil {
		return nil, errors.New("form not found")
	}
	if _, err := s.clinicRepo.GetByID(ctx, form.ClinicID); err != nil {
		return nil, errors.New("clinic not found")
	}

	entry, err := req.ToFieldEntryDB()
	if err != nil {
		return nil, err
	}
	if err := s.repo.Create(ctx, entry); err != nil {
		return nil, err
	}
	return entry.ToFieldEntryResponse(), nil
}

func (s *FieldEntryService) GetByID(ctx context.Context, id string) (*form.FieldEntryResponse, error) {
	entry, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("field entry not found")
	}

	return entry.ToFieldEntryResponse(), nil
}

func (s *FieldEntryService) GetByFormID(ctx context.Context, formID string) ([]*form.FieldEntryResponse, error) {
	fm, err := s.formRepo.GetByID(ctx, formID)
	if err != nil {
		return nil, errors.New("form not found")
	}

	fieldEntries, err := s.repo.GetByFormID(ctx, fm.ID)
	if err != nil {
		return nil, err
	}
	if len(fieldEntries) == 0 {
		return nil, nil
	}

	fieldEntryResponses := make([]*form.FieldEntryResponse, 0, len(fieldEntries))
	for _, entry := range fieldEntries {
		resp := entry.ToFieldEntryResponse()
		fieldEntryResponses = append(fieldEntryResponses, resp)
	}

	return fieldEntryResponses, nil
}

func (s *FieldEntryService) GetByClinicID(ctx context.Context, clinicID string) ([]*form.FieldEntryResponse, error) {
	fieldEntries, err := s.repo.GetByClinicID(ctx, clinicID)
	if err != nil {
		return nil, err
	}
	if len(fieldEntries) == 0 {
		return []*form.FieldEntryResponse{}, nil
	}

	type entryKey struct {
		formID    string
		createdAt int64
	}
	entryGroups := make(map[entryKey][]form.FieldEntry)
	formIDs := make(map[string]bool)

	for _, entry := range fieldEntries {
		key := entryKey{
			formID:    entry.FormID,
			createdAt: entry.CreatedAt.Unix(),
		}
		entryGroups[key] = append(entryGroups[key], entry)
		formIDs[entry.FormID] = true
	}

	forms := make(map[string]*form.Form)
	for formID := range formIDs {
		form, err := s.formRepo.GetByID(ctx, formID)
		if err == nil {
			forms[formID] = form
		}
	}

	fieldsMap := make(map[string]form.Field)
	for formID := range formIDs {
		fields, _ := s.fieldRepo.GetByFormID(ctx, formID)
		for _, field := range fields {
			fieldsMap[field.ID] = field
		}
	}

	fieldEntryResponses := make([]*form.FieldEntryResponse, 0, len(entryGroups))
	for _, entries := range entryGroups {
		for _, entry := range entries {
			fieldEntryResponses = append(fieldEntryResponses, entry.ToFieldEntryResponse())
		}
	}
	return fieldEntryResponses, nil
}

func (s *FieldEntryService) Update(ctx context.Context, req *form.EntryFieldUpdateRequest) (*form.FieldEntryResponse, error) {
	entry, err := s.repo.GetByID(ctx, req.ID)
	if err != nil {
		return nil, errors.New("field entry not found")
	}

	entry.Value = req.Value
	entry.GSTAmount = req.GSTAmount
	entry.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, entry); err != nil {
		return nil, err
	}

	return entry.ToFieldEntryResponse(), nil
}

func (s *FieldEntryService) Delete(ctx context.Context, id string) error {
	entry, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return errors.New("field entry not found")
	}
	_ = entry
	return s.repo.Delete(ctx, id)
}
