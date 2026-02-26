package usecase

import (
	"context"
	"errors"

	"github.com/iamarpitzala/aca-reca-backend/internal/application/port"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain/form"
)

type CustomFormFieldService struct {
	fieldRepo  port.CustomFormFieldRepository
	configRepo port.CustomFormFieldConfigRepository
	formRepo   port.CustomFormRepository
}

func NewCustomFormFieldService(
	fieldRepo port.CustomFormFieldRepository,
	configRepo port.CustomFormFieldConfigRepository,
	formRepo port.CustomFormRepository,
) *CustomFormFieldService {
	return &CustomFormFieldService{
		fieldRepo:  fieldRepo,
		configRepo: configRepo,
		formRepo:   formRepo,
	}
}

func (s *CustomFormFieldService) CreateField(ctx context.Context, req *form.FieldRequest) (*form.FieldResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, err
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
	fields, err := s.fieldRepo.GetByFormVersionID(ctx, formVersionID)
	if err != nil {
		return nil, err
	}
	out := make([]*form.FieldResponse, len(fields))
	for i := range fields {
		out[i] = fields[i].ToFieldResponse()
	}
	return out, nil
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
