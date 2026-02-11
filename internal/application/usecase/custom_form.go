package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/application/port"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
)

type CustomFormService struct {
	repo       port.CustomFormRepository
	clinicRepo port.ClinicRepository
}

func NewCustomFormService(repo port.CustomFormRepository, clinicRepo port.ClinicRepository) *CustomFormService {
	return &CustomFormService{repo: repo, clinicRepo: clinicRepo}
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

func (s *CustomFormService) Publish(ctx context.Context, id uuid.UUID) error {
	return s.repo.Publish(ctx, id)
}

func (s *CustomFormService) Archive(ctx context.Context, id uuid.UUID) error {
	return s.repo.Archive(ctx, id)
}

func (s *CustomFormService) Delete(ctx context.Context, id uuid.UUID) error {
	return s.repo.Delete(ctx, id)
}

func (s *CustomFormService) CreateEntry(ctx context.Context, entry *domain.CustomFormEntry) error {
	return s.repo.CreateEntry(ctx, entry)
}

func (s *CustomFormService) GetEntryByID(ctx context.Context, id uuid.UUID) (*domain.CustomFormEntry, error) {
	return s.repo.GetEntryByID(ctx, id)
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

func (s *CustomFormService) UpdateEntry(ctx context.Context, entry *domain.CustomFormEntry) error {
	return s.repo.UpdateEntry(ctx, entry)
}

func (s *CustomFormService) DeleteEntry(ctx context.Context, id uuid.UUID) error {
	return s.repo.DeleteEntry(ctx, id)
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
