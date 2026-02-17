package usecase

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/application/port"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
	"github.com/iamarpitzala/aca-reca-backend/util"
)

type CustomFormService struct {
	repo        port.CustomFormRepository
	fieldRepo   port.CustomFormFieldRepository
	versionRepo port.CustomFormVersionRepository
	clinicRepo  port.ClinicRepository
	calcEngine  port.EntryCalculationEngine
}

func NewCustomFormService(
	repo port.CustomFormRepository,
	fieldRepo port.CustomFormFieldRepository,
	versionRepo port.CustomFormVersionRepository,
	clinicRepo port.ClinicRepository,
	calcEngine port.EntryCalculationEngine,
) *CustomFormService {
	return &CustomFormService{
		repo:        repo,
		fieldRepo:   fieldRepo,
		versionRepo: versionRepo,
		clinicRepo:  clinicRepo,
		calcEngine:  calcEngine,
	}
}

func (s *CustomFormService) Create(ctx context.Context, req *domain.CreateCustomFormRequest, userID uuid.UUID) (*domain.CustomFormResponse, error) {
	clinicID, err := uuid.Parse(req.ClinicID)
	if err != nil {
		return nil, errors.New("invalid clinic ID")
	}
	if _, err := s.clinicRepo.GetByID(ctx, clinicID); err != nil {
		return nil, errors.New("clinic not found")
	}
	// Normalise to UPPERCASE so API accepts both lowercase and uppercase
	req.FormType = strings.ToUpper(strings.TrimSpace(req.FormType))
	if req.FormType != util.FormTypeIncome && req.FormType != util.FormTypeExpense && req.FormType != util.FormTypeBoth {
		return nil, errors.New("form type must be INCOME, EXPENSE, or BOTH")
	}
	defaultPayment := req.DefaultPaymentResponsibility
	if defaultPayment == nil || *defaultPayment == "" {
		v := util.PaymentResponsibilityOwner
		defaultPayment = &v
	}
	form, err := req.ToDBModel(userID)
	if err != nil {
		return nil, err
	}
	form.CalculationMethod = req.CalculationMethod
	form.FormType = req.FormType
	form.DefaultPaymentResponsibility = defaultPayment
	if err := s.repo.Create(ctx, form); err != nil {
		return nil, err
	}

	// Create initial form version (version 1)
	formVersion := &domain.CustomFormVersion{
		ID:        uuid.New(),
		FormID:    form.ID,
		Version:   1,
		IsActive:  true,
		CreatedBy: userID,
		CreatedAt: time.Now(),
	}
	if err := s.versionRepo.Create(ctx, formVersion); err != nil {
		return nil, err
	}

	// Save fields if provided (deduplicate by field_key to prevent duplicate entries)
	if len(req.Fields) > 0 {
		fields := make([]*domain.CustomFormField, 0, len(req.Fields))
		seenKeys := make(map[string]bool)
		for _, fieldInput := range req.Fields {
			key := strings.TrimSpace(strings.ToLower(fieldInput.Name))
			if key == "" || seenKeys[key] {
				continue
			}
			seenKeys[key] = true
			field, err := fieldInput.ToDBModel(form.ID, formVersion.ID, userID)
			if err != nil {
				return nil, err
			}
			fields = append(fields, field)
		}
		if len(fields) > 0 {
			if err := s.fieldRepo.CreateBatch(ctx, fields); err != nil {
				return nil, err
			}
		}
	}

	// Fetch fields for response
	fieldList, _ := s.fieldRepo.GetByFormVersionID(ctx, formVersion.ID)
	fieldResponses := make([]domain.CustomFormFieldResponse, len(fieldList))
	for i := range fieldList {
		fieldResponses[i] = *fieldList[i].ToResponse()
	}

	return form.ToResponse(fieldResponses), nil
}

func (s *CustomFormService) GetByID(ctx context.Context, id uuid.UUID) (*domain.CustomFormResponse, error) {
	form, err := s.getFormByID(ctx, id)
	if err != nil {
		return nil, err
	}

	// Get latest version
	version, err := s.versionRepo.GetLatestByFormID(ctx, id)
	if err != nil {
		// If no version exists, return form without fields
		return form.ToResponse([]domain.CustomFormFieldResponse{}), nil
	}

	// Fetch fields for the latest version
	fieldList, _ := s.fieldRepo.GetByFormVersionID(ctx, version.ID)
	fieldResponses := make([]domain.CustomFormFieldResponse, len(fieldList))
	for i := range fieldList {
		fieldResponses[i] = *fieldList[i].ToResponse()
	}

	return form.ToResponse(fieldResponses), nil
}

func (s *CustomFormService) getFormByID(ctx context.Context, id uuid.UUID) (*domain.CustomForm, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *CustomFormService) GetByClinicID(ctx context.Context, clinicID uuid.UUID) ([]domain.CustomFormResponse, error) {
	forms, err := s.repo.GetByClinicID(ctx, clinicID)
	if err != nil {
		return nil, err
	}
	out := make([]domain.CustomFormResponse, len(forms))
	for i := range forms {
		// Get latest version and fields for each form
		version, _ := s.versionRepo.GetLatestByFormID(ctx, forms[i].ID)
		var fieldResponses []domain.CustomFormFieldResponse
		if version != nil {
			fieldList, _ := s.fieldRepo.GetByFormVersionID(ctx, version.ID)
			fieldResponses = make([]domain.CustomFormFieldResponse, len(fieldList))
			for j := range fieldList {
				fieldResponses[j] = *fieldList[j].ToResponse()
			}
		}
		out[i] = *forms[i].ToResponse(fieldResponses)
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
		// Get active version and fields for each published form
		version, _ := s.versionRepo.GetLatestByFormID(ctx, forms[i].ID)
		var fieldResponses []domain.CustomFormFieldResponse
		if version != nil {
			fieldList, _ := s.fieldRepo.GetByFormVersionID(ctx, version.ID)
			fieldResponses = make([]domain.CustomFormFieldResponse, len(fieldList))
			for j := range fieldList {
				fieldResponses[j] = *fieldList[j].ToResponse()
			}
		}
		out[i] = *forms[i].ToResponse(fieldResponses)
	}
	return out, nil
}

func (s *CustomFormService) Update(ctx context.Context, form *domain.CustomForm) error {
	return s.repo.Update(ctx, form)
}

func (s *CustomFormService) UpdateByRequest(ctx context.Context, id uuid.UUID, req *domain.UpdateCustomFormRequest, userID uuid.UUID) (*domain.CustomFormResponse, error) {
	form, err := s.getFormByID(ctx, id)
	if err != nil {
		return nil, err
	}
	applyUpdateToForm(form, req)
	form.UpdatedAt = time.Now()
	if err := s.repo.Update(ctx, form); err != nil {
		return nil, err
	}

	// Update fields if provided
	if len(req.Fields) > 0 {
		// Get or create latest version
		version, err := s.versionRepo.GetLatestByFormID(ctx, id)
		if err != nil {
			// Create new version if none exists
			version = &domain.CustomFormVersion{
				ID:        uuid.New(),
				FormID:    id,
				Version:   1,
				IsActive:  true,
				CreatedBy: userID,
				CreatedAt: time.Now(),
			}
			if err := s.versionRepo.Create(ctx, version); err != nil {
				return nil, err
			}
		}

		// Delete existing fields for this version
		if err := s.fieldRepo.DeleteByFormVersionID(ctx, version.ID); err != nil {
			return nil, err
		}

		// Create new fields (deduplicate by field_key to prevent duplicate entries)
		fields := make([]*domain.CustomFormField, 0, len(req.Fields))
		seenKeys := make(map[string]bool)
		for _, fieldInput := range req.Fields {
			key := strings.TrimSpace(strings.ToLower(fieldInput.Name))
			if key == "" || seenKeys[key] {
				continue
			}
			seenKeys[key] = true
			field, err := fieldInput.ToDBModel(id, version.ID, userID)
			if err != nil {
				return nil, err
			}
			fields = append(fields, field)
		}
		if len(fields) > 0 {
			if err := s.fieldRepo.CreateBatch(ctx, fields); err != nil {
				return nil, err
			}
		}
	}

	// Fetch fields for response
	version, _ := s.versionRepo.GetLatestByFormID(ctx, id)
	var fieldResponses []domain.CustomFormFieldResponse
	if version != nil {
		fieldList, _ := s.fieldRepo.GetByFormVersionID(ctx, version.ID)
		fieldResponses = make([]domain.CustomFormFieldResponse, len(fieldList))
		for i := range fieldList {
			fieldResponses[i] = *fieldList[i].ToResponse()
		}
	}

	return form.ToResponse(fieldResponses), nil
}

func applyUpdateToForm(form *domain.CustomForm, req *domain.UpdateCustomFormRequest) {
	if req.Name != nil {
		form.Name = *req.Name
	}
	if req.Description != nil {
		form.Description = *req.Description
	}
	if req.DefaultPaymentResponsibility != nil {
		form.DefaultPaymentResponsibility = req.DefaultPaymentResponsibility
	}
	if req.CalculationMethod != nil {
		form.CalculationMethod = *req.CalculationMethod
	}
	if req.FormType != nil {
		form.FormType = *req.FormType
	}
}

func (s *CustomFormService) Publish(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*domain.CustomFormResponse, error) {
	form, err := s.getFormByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if form.Status == util.FormStatusPublished {
		return nil, errors.New("form is already published")
	}

	// Get latest version
	latestVersion, err := s.versionRepo.GetLatestByFormID(ctx, id)
	if err != nil {
		return nil, errors.New("form must have at least one version to publish")
	}

	// Create new published version
	versions, _ := s.versionRepo.GetByFormID(ctx, id)
	newVersionNum := len(versions) + 1
	newVersion := &domain.CustomFormVersion{
		ID:        uuid.New(),
		FormID:    id,
		Version:   newVersionNum,
		IsActive:  true,
		CreatedBy: userID,
		CreatedAt: time.Now(),
	}
	if err := s.versionRepo.Create(ctx, newVersion); err != nil {
		return nil, err
	}

	// Copy fields from latest version to new published version
	oldFields, err := s.fieldRepo.GetByFormVersionID(ctx, latestVersion.ID)
	if err == nil && len(oldFields) > 0 {
		newFields := make([]*domain.CustomFormField, len(oldFields))
		for i := range oldFields {
			newFields[i] = &domain.CustomFormField{
				ID:            uuid.New(),
				FormVersionID: newVersion.ID,
				FormID:        id,
				FieldKey:      oldFields[i].FieldKey,
				Label:         oldFields[i].Label,
				FieldType:     oldFields[i].FieldType,
				Section:       oldFields[i].Section,
				IsRequired:    oldFields[i].IsRequired,
				CoaID:         oldFields[i].CoaID,
				Placeholder:   oldFields[i].Placeholder,
				MinValue:      oldFields[i].MinValue,
				MaxValue:      oldFields[i].MaxValue,
				FieldOrder:    oldFields[i].FieldOrder,
				GSTConfig:     oldFields[i].GSTConfig,
				GSTRate:       oldFields[i].GSTRate,
				GSTType:       oldFields[i].GSTType,
				Metadata:      oldFields[i].Metadata,
			}
		}
		if err := s.fieldRepo.CreateBatch(ctx, newFields); err != nil {
			return nil, err
		}
	}

	// Set new version as active
	if err := s.versionRepo.SetActive(ctx, id, newVersion.ID); err != nil {
		return nil, err
	}

	// Publish the form
	if err := s.repo.Publish(ctx, id); err != nil {
		return nil, err
	}
	form.Status = util.FormStatusPublished
	form.UpdatedAt = time.Now()

	// Fetch fields for response
	fieldList, _ := s.fieldRepo.GetByFormVersionID(ctx, newVersion.ID)
	fieldResponses := make([]domain.CustomFormFieldResponse, len(fieldList))
	for i := range fieldList {
		fieldResponses[i] = *fieldList[i].ToResponse()
	}

	return form.ToResponse(fieldResponses), nil
}

func (s *CustomFormService) Unpublish(ctx context.Context, id uuid.UUID) (*domain.CustomFormResponse, error) {
	form, err := s.getFormByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if form.Status != util.FormStatusPublished {
		return nil, errors.New("only published forms can be unpublished")
	}
	if err := s.repo.Unpublish(ctx, id); err != nil {
		return nil, err
	}
	form.Status = util.FormStatusDraft
	form.UpdatedAt = time.Now()

	// Fetch fields for response
	version, _ := s.versionRepo.GetLatestByFormID(ctx, id)
	var fieldResponses []domain.CustomFormFieldResponse
	if version != nil {
		fieldList, _ := s.fieldRepo.GetByFormVersionID(ctx, version.ID)
		fieldResponses = make([]domain.CustomFormFieldResponse, len(fieldList))
		for i := range fieldList {
			fieldResponses[i] = *fieldList[i].ToResponse()
		}
	}
	return form.ToResponse(fieldResponses), nil
}

func (s *CustomFormService) Archive(ctx context.Context, id uuid.UUID) (*domain.CustomFormResponse, error) {
	form, err := s.getFormByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := s.repo.Archive(ctx, id); err != nil {
		return nil, err
	}
	form.Status = util.FormStatusArchived
	form.UpdatedAt = time.Now()

	// Fetch fields for response
	version, _ := s.versionRepo.GetLatestByFormID(ctx, id)
	var fieldResponses []domain.CustomFormFieldResponse
	if version != nil {
		fieldList, _ := s.fieldRepo.GetByFormVersionID(ctx, version.ID)
		fieldResponses = make([]domain.CustomFormFieldResponse, len(fieldList))
		for i := range fieldList {
			fieldResponses[i] = *fieldList[i].ToResponse()
		}
	}
	return form.ToResponse(fieldResponses), nil
}

func (s *CustomFormService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *CustomFormService) Duplicate(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*domain.CustomFormResponse, error) {
	form, err := s.getFormByID(ctx, id)
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
		Status:                       util.FormStatusDraft,
		DefaultPaymentResponsibility: form.DefaultPaymentResponsibility,
		CreatedBy:                    userID,
		CreatedAt:                    now,
		UpdatedAt:                    now,
	}
	if err := s.repo.Create(ctx, newForm); err != nil {
		return nil, err
	}

	// Create initial version for duplicated form
	formVersion := &domain.CustomFormVersion{
		ID:        uuid.New(),
		FormID:    newForm.ID,
		Version:   1,
		IsActive:  true,
		CreatedBy: userID,
		CreatedAt: now,
	}
	if err := s.versionRepo.Create(ctx, formVersion); err != nil {
		return nil, err
	}

	// Copy fields from original form
	originalVersion, _ := s.versionRepo.GetLatestByFormID(ctx, id)
	if originalVersion != nil {
		originalFields, _ := s.fieldRepo.GetByFormVersionID(ctx, originalVersion.ID)
		if len(originalFields) > 0 {
			newFields := make([]*domain.CustomFormField, len(originalFields))
			for i := range originalFields {
				newFields[i] = &domain.CustomFormField{
					ID:            uuid.New(),
					FormVersionID: formVersion.ID,
					FormID:        newForm.ID,
					FieldKey:      originalFields[i].FieldKey,
					Label:         originalFields[i].Label,
					FieldType:     originalFields[i].FieldType,
					Section:       originalFields[i].Section,
					IsRequired:    originalFields[i].IsRequired,
					CoaID:         originalFields[i].CoaID,
					Placeholder:   originalFields[i].Placeholder,
					MinValue:      originalFields[i].MinValue,
					MaxValue:      originalFields[i].MaxValue,
					FieldOrder:    originalFields[i].FieldOrder,
					GSTConfig:     originalFields[i].GSTConfig,
					GSTRate:       originalFields[i].GSTRate,
					GSTType:       originalFields[i].GSTType,
					Metadata:      originalFields[i].Metadata,
				}
			}
			if err := s.fieldRepo.CreateBatch(ctx, newFields); err != nil {
				return nil, err
			}
		}
	}

	// Fetch fields for response
	fieldList, _ := s.fieldRepo.GetByFormVersionID(ctx, formVersion.ID)
	fieldResponses := make([]domain.CustomFormFieldResponse, len(fieldList))
	for i := range fieldList {
		fieldResponses[i] = *fieldList[i].ToResponse()
	}

	return newForm.ToResponse(fieldResponses), nil
}
