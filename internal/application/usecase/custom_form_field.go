package usecase

import (
	"context"
	"errors"
	"strconv"

	"github.com/iamarpitzala/aca-reca-backend/internal/application/port"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain/form"
)

type CustomFormFieldService struct {
	fieldRepo   port.CustomFormFieldRepository
	sectionRepo port.CustomFormSectionRepository
	versionRepo port.CustomFormVersionRepository
	configRepo  port.CustomFormFieldConfigRepository
	formRepo    port.CustomFormRepository
}

func NewCustomFormFieldService(
	fieldRepo port.CustomFormFieldRepository,
	sectionRepo port.CustomFormSectionRepository,
	versionRepo port.CustomFormVersionRepository,
	configRepo port.CustomFormFieldConfigRepository,
	formRepo port.CustomFormRepository,
) *CustomFormFieldService {
	return &CustomFormFieldService{
		fieldRepo:   fieldRepo,
		sectionRepo: sectionRepo,
		versionRepo: versionRepo,
		configRepo:  configRepo,
		formRepo:    formRepo,
	}
}

func (s *CustomFormFieldService) CreateField(ctx context.Context, req *form.FieldRequest) (*form.FieldResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.sectionRepo.GetByID(ctx, req.SectionID); err != nil {
		return nil, errors.New("section not found")
	}
	f := &form.Field{}
	f.ToFieldDB(req)
	if err := s.fieldRepo.Create(ctx, f); err != nil {
		return nil, err
	}
	return f.ToFieldResponse(), nil
}

func (s *CustomFormFieldService) GetFieldByID(ctx context.Context, id string) (*form.FieldResponse, error) {
	f, err := s.fieldRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return f.ToFieldResponse(), nil
}

func (s *CustomFormFieldService) GetFieldsByFormID(ctx context.Context, formID string) ([]*form.FieldResponse, error) {
	fields, err := s.fieldRepo.GetByFormID(ctx, formID)
	if err != nil {
		return nil, err
	}
	out := make([]*form.FieldResponse, len(fields))
	for i := range fields {
		out[i] = fields[i].ToFieldResponse()
	}
	return out, nil
}

func (s *CustomFormFieldService) GetFieldsByFormVersionID(ctx context.Context, formVersionID string) ([]*form.FieldResponse, error) {
	vid, err := strconv.Atoi(formVersionID)
	if err != nil {
		return nil, errors.New("invalid formVersionId")
	}
	fields, err := s.fieldRepo.GetByFormVersionID(ctx, vid)
	if err != nil {
		return nil, err
	}
	out := make([]*form.FieldResponse, len(fields))
	for i := range fields {
		out[i] = fields[i].ToFieldResponse()
	}
	return out, nil
}

// GetFormIDForField returns the form ID that owns this field (section -> version -> form). For access control.
func (s *CustomFormFieldService) GetFormIDForField(ctx context.Context, fieldID string) (string, error) {
	f, err := s.fieldRepo.GetByID(ctx, fieldID)
	if err != nil {
		return "", errors.New("field not found")
	}
	sec, err := s.sectionRepo.GetByID(ctx, f.SectionID)
	if err != nil || sec == nil {
		return "", errors.New("section not found")
	}
	version, err := s.versionRepo.GetByID(ctx, sec.FormVersionID)
	if err != nil || version == nil {
		return "", errors.New("form version not found")
	}
	return version.FormID, nil
}

// GetFormIDForSection returns the form ID for the form version that owns this section (for access control).
func (s *CustomFormFieldService) GetFormIDForSection(ctx context.Context, sectionID int) (string, error) {
	sec, err := s.sectionRepo.GetByID(ctx, sectionID)
	if err != nil || sec == nil {
		return "", errors.New("section not found")
	}
	version, err := s.versionRepo.GetByID(ctx, sec.FormVersionID)
	if err != nil || version == nil {
		return "", errors.New("form version not found")
	}
	return version.FormID, nil
}

func (s *CustomFormFieldService) UpdateField(ctx context.Context, id string, req *form.FieldRequest) (*form.FieldResponse, error) {
	existing, err := s.fieldRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := req.Validate(); err != nil {
		return nil, err
	}
	existing.ToFieldDB(req)
	existing.ID = id
	if err := s.fieldRepo.Update(ctx, existing); err != nil {
		return nil, err
	}
	return existing.ToFieldResponse(), nil
}

func (s *CustomFormFieldService) DeleteField(ctx context.Context, id string) error {
	return s.fieldRepo.DeleteByID(ctx, id)
}

func (s *CustomFormFieldService) CreateFieldConfig(ctx context.Context, req *form.FieldConfigRequest) (*form.FieldConfigResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
	}
	if _, err := s.fieldRepo.GetByID(ctx, req.FormFieldID); err != nil {
		return nil, errors.New("form field not found")
	}
	c := &form.FieldConfig{}
	c.FromRequest(req)
	if err := s.configRepo.Create(ctx, c); err != nil {
		return nil, err
	}
	return c.ToResponse(), nil
}

func (s *CustomFormFieldService) GetFieldConfigByID(ctx context.Context, id string) (*form.FieldConfigResponse, error) {
	c, err := s.configRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return c.ToResponse(), nil
}

func (s *CustomFormFieldService) GetFieldConfigsByFormFieldID(ctx context.Context, formFieldID string) ([]*form.FieldConfigResponse, error) {
	list, err := s.configRepo.GetByFormFieldID(ctx, formFieldID)
	if err != nil {
		return nil, err
	}
	out := make([]*form.FieldConfigResponse, len(list))
	for i := range list {
		out[i] = list[i].ToResponse()
	}
	return out, nil
}

func (s *CustomFormFieldService) UpdateFieldConfig(ctx context.Context, id string, req *form.FieldConfigRequest) (*form.FieldConfigResponse, error) {
	existing, err := s.configRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if err := req.Validate(); err != nil {
		return nil, err
	}
	existing.FromRequest(req)
	existing.ID = id
	if err := s.configRepo.Update(ctx, existing); err != nil {
		return nil, err
	}
	return existing.ToResponse(), nil
}

func (s *CustomFormFieldService) DeleteFieldConfig(ctx context.Context, id string) error {
	return s.configRepo.DeleteByID(ctx, id)
}
