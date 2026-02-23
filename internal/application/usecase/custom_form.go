package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/application/port"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain/form"
	"github.com/iamarpitzala/aca-reca-backend/util"
)

type CustomFormService struct {
	repo        port.CustomFormRepository
	fieldRepo   port.CustomFormFieldRepository
	versionRepo port.CustomFormVersionRepository
	clinicRepo  port.ClinicRepository
	calcEngine  port.EntryCalculationEngine
}

func NewCustomFormService(repo port.CustomFormRepository, fieldRepo port.CustomFormFieldRepository, versionRepo port.CustomFormVersionRepository, clinicRepo port.ClinicRepository, calcEngine port.EntryCalculationEngine) *CustomFormService {
	return &CustomFormService{
		repo:        repo,
		fieldRepo:   fieldRepo,
		versionRepo: versionRepo,
		clinicRepo:  clinicRepo,
		calcEngine:  calcEngine,
	}
}

func (s *CustomFormService) Create(ctx context.Context, req *form.FormRequest, userID string) (*form.Form, error) {
	if _, err := s.clinicRepo.GetByID(ctx, req.ClinicID); err != nil {
		return nil, errors.New("clinic not found")
	}

	formInstance := &form.Form{}
	formInstance.ToFormDB(req)
	if err := s.repo.Create(ctx, formInstance); err != nil {
		return nil, err
	}

	// Initial version logic: create version 1 and activate it
	formVersion := &form.FormVersion{}
	formVersionReq := &form.FormVersionRequest{
		FormID:    formInstance.ID,
		Version:   1,
		IsActive:  true,
		CreatedBy: userID,
		CreatedAt: time.Now(),
	}
	formVersion.ToFormVersionDB(formVersionReq)
	if err := s.versionRepo.Create(ctx, formVersion); err != nil {
		return nil, err
	}

	return formInstance, nil
}

func (s *CustomFormService) GetByID(ctx context.Context, formID string) (*form.FormResponse, error) {
	f, err := s.repo.GetByID(ctx, formID)
	if err != nil {
		return nil, err
	}

	return f.ToFormResponse(), nil
}

func (s *CustomFormService) GetByClinicID(ctx context.Context, clinicID string) ([]form.FormResponse, error) {
	forms, err := s.repo.GetByClinicID(ctx, clinicID)
	if err != nil {
		return nil, err
	}

	out := make([]form.FormResponse, len(forms))
	for i, f := range forms {
		out[i] = *f.ToFormResponse()
	}
	return out, nil
}

func (s *CustomFormService) GetPublishedByClinicID(ctx context.Context, clinicID string) ([]form.FormResponse, error) {
	forms, err := s.repo.GetPublishedByClinicID(ctx, clinicID)
	if err != nil {
		return nil, err
	}

	out := make([]form.FormResponse, len(forms))
	for i, f := range forms {
		out[i] = *f.ToFormResponse()
	}
	return out, nil
}

func (s *CustomFormService) Update(ctx context.Context, f *form.Form) error {
	return s.repo.Update(ctx, f)
}

func (s *CustomFormService) UpdateByRequest(ctx context.Context, id string, req *form.FormRequest, userID string) (*form.FormResponse, error) {
	// Fetch current form
	formInstance, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	// Update logic: update with current request's content
	formInstance.Name = req.Name
	formInstance.Description = req.Description
	formInstance.Status = req.Status
	formInstance.AccountingMethod = req.AccountingMethod
	formInstance.UpdatedAt = time.Now()
	formInstance.DeletedAt = req.DeletedAt
	formInstance.PublishedAt = req.PublishedAt

	if err := s.repo.Update(ctx, formInstance); err != nil {
		return nil, err
	}

	return formInstance.ToFormResponse(), nil
}

func (s *CustomFormService) Publish(ctx context.Context, id, userID string) (*form.FormResponse, error) {
	f, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if f.Status == string(util.FormStatusPublished) {
		return nil, errors.New("form is already published")
	}

	latestVersion, err := s.versionRepo.GetLatestByFormID(ctx, id)
	if err != nil || latestVersion == nil {
		return nil, errors.New("form must have at least one version to publish")
	}

	// Create a new version (copy fields from latest)
	newVersionNum := latestVersion.Version + 1
	newVersionRecord := &form.FormVersion{}
	newVersionReq := &form.FormVersionRequest{
		FormID:    id,
		Version:   newVersionNum,
		IsActive:  true,
		CreatedBy: userID,
		CreatedAt: time.Now(),
	}
	newVersionRecord.ToFormVersionDB(newVersionReq)
	if err := s.versionRepo.Create(ctx, newVersionRecord); err != nil {
		return nil, err
	}

	// Copy fields if any
	oldFields, err := s.fieldRepo.GetByFormVersionID(ctx, latestVersion.ID)
	if err != nil {
		// If error, treat as no fields
		oldFields = nil
	}
	if len(oldFields) > 0 {
		newFields := make([]*form.Field, len(oldFields))
		for i, old := range oldFields {
			newField := &form.Field{
				ID:            uuid.NewString(),
				FormVersionID: newVersionRecord.ID,
				FormID:        id,
				Label:         old.Label,
				SectionTypeID: old.SectionTypeID,
				IsRequired:    old.IsRequired,
				CoaID:         old.CoaID,
				Placeholder:   old.Placeholder,
				MinValue:      old.MinValue,
				MaxValue:      old.MaxValue,
				FieldOrder:    old.FieldOrder,
			}
			newFields[i] = newField
		}
		if err := s.fieldRepo.CreateBatch(ctx, newFields); err != nil {
			return nil, err
		}
	}

	// Set the new version as active
	if err := s.versionRepo.SetActive(ctx, id, newVersionRecord.ID); err != nil {
		return nil, err
	}

	// Set form as published
	if err := s.repo.Publish(ctx, id); err != nil {
		return nil, err
	}
	f.Status = string(util.FormStatusPublished)
	f.UpdatedAt = time.Now()

	return f.ToFormResponse(), nil
}

func (s *CustomFormService) Unpublish(ctx context.Context, id string) (*form.FormResponse, error) {
	f, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if f.Status != string(util.FormStatusPublished) {
		return nil, errors.New("only published forms can be unpublished")
	}
	if err := s.repo.Unpublish(ctx, id); err != nil {
		return nil, err
	}
	f.Status = string(util.FormStatusDraft)
	f.UpdatedAt = time.Now()

	return f.ToFormResponse(), nil
}

func (s *CustomFormService) Archive(ctx context.Context, id string) (*form.FormResponse, error) {
	f, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}

	if err := s.repo.Archive(ctx, id); err != nil {
		return nil, err
	}
	f.Status = string(util.FormStatusArchived)
	f.UpdatedAt = time.Now()

	return f.ToFormResponse(), nil
}

func (s *CustomFormService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
