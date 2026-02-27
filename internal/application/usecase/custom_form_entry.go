package usecase

import (
	"context"
	"errors"
	"strconv"

	"github.com/iamarpitzala/aca-reca-backend/internal/application/port"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain/form"
)

type CustomFormEntryService struct {
	entryRepo  port.CustomFormEntryRepository
	versionRepo port.CustomFormVersionRepository
	formRepo   port.CustomFormRepository
}

func NewCustomFormEntryService(
	entryRepo port.CustomFormEntryRepository,
	versionRepo port.CustomFormVersionRepository,
	formRepo port.CustomFormRepository,
) *CustomFormEntryService {
	return &CustomFormEntryService{
		entryRepo:  entryRepo,
		versionRepo: versionRepo,
		formRepo:   formRepo,
	}
}

func (s *CustomFormEntryService) Create(ctx context.Context, req *form.EntryCreateRequest, userID string) (*form.EntryResponse, error) {
	formVersionID, err := strconv.Atoi(req.FormVersionID)
	if err != nil {
		return nil, errors.New("invalid formVersionId")
	}
	version, err := s.versionRepo.GetByID(ctx, formVersionID)
	if err != nil || version == nil {
		return nil, errors.New("form version not found")
	}
	fm, err := s.formRepo.GetByID(ctx, version.FormID)
	if err != nil || fm == nil {
		return nil, errors.New("form not found")
	}

	var submittedBy *string
	if userID != "" {
		submittedBy = &userID
	}
	entry := form.NewEntry(formVersionID, submittedBy)
	values := make([]form.EntryValue, 0, len(req.Values))
	for _, v := range req.Values {
		values = append(values, form.EntryValueForCreate(entry.ID, v.FieldID, v.Value, v.GSTAmount))
	}
	if err := s.entryRepo.Create(ctx, entry, values); err != nil {
		return nil, err
	}
	return entry.ToEntryResponse(values), nil
}

func (s *CustomFormEntryService) GetByID(ctx context.Context, id string) (*form.EntryResponse, error) {
	entry, values, err := s.entryRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.New("entry not found")
	}
	return entry.ToEntryResponse(values), nil
}

func (s *CustomFormEntryService) GetByFormVersionID(ctx context.Context, formVersionID int) ([]*form.EntryResponse, error) {
	entries, err := s.entryRepo.GetByFormVersionID(ctx, formVersionID)
	if err != nil {
		return nil, err
	}
	out := make([]*form.EntryResponse, 0, len(entries))
	for _, e := range entries {
		_, values, _ := s.entryRepo.GetByID(ctx, e.ID)
		if values == nil {
			values = []form.EntryValue{}
		}
		out = append(out, e.ToEntryResponse(values))
	}
	return out, nil
}

func (s *CustomFormEntryService) GetByFormID(ctx context.Context, formID string) ([]*form.EntryResponse, error) {
	entries, err := s.entryRepo.GetByFormID(ctx, formID)
	if err != nil {
		return nil, err
	}
	out := make([]*form.EntryResponse, 0, len(entries))
	for _, e := range entries {
		_, values, _ := s.entryRepo.GetByID(ctx, e.ID)
		if values == nil {
			values = []form.EntryValue{}
		}
		out = append(out, e.ToEntryResponse(values))
	}
	return out, nil
}

func (s *CustomFormEntryService) GetByClinicID(ctx context.Context, clinicID string) ([]*form.EntryResponse, error) {
	entries, err := s.entryRepo.GetByClinicID(ctx, clinicID)
	if err != nil {
		return nil, err
	}
	out := make([]*form.EntryResponse, 0, len(entries))
	for _, e := range entries {
		_, values, _ := s.entryRepo.GetByID(ctx, e.ID)
		if values == nil {
			values = []form.EntryValue{}
		}
		out = append(out, e.ToEntryResponse(values))
	}
	return out, nil
}

func (s *CustomFormEntryService) Update(ctx context.Context, entryID string, req *form.EntryUpdateRequest) (*form.EntryResponse, error) {
	entry, _, err := s.entryRepo.GetByID(ctx, entryID)
	if err != nil {
		return nil, errors.New("entry not found")
	}
	newValues := make([]form.EntryValue, 0, len(req.Values))
	for _, v := range req.Values {
		newValues = append(newValues, form.EntryValueForCreate(entryID, v.FieldID, v.Value, v.GSTAmount))
	}
	if err := s.entryRepo.Update(ctx, entryID, newValues); err != nil {
		return nil, err
	}
	entry, updatedValues, _ := s.entryRepo.GetByID(ctx, entryID)
	return entry.ToEntryResponse(updatedValues), nil
}

func (s *CustomFormEntryService) Delete(ctx context.Context, id string) error {
	_, _, err := s.entryRepo.GetByID(ctx, id)
	if err != nil {
		return errors.New("entry not found")
	}
	return s.entryRepo.Delete(ctx, id)
}

// GetClinicIDForEntry returns the clinic ID for the form that owns this entry (for access control).
func (s *CustomFormEntryService) GetClinicIDForEntry(ctx context.Context, entryID string) (string, error) {
	entry, _, err := s.entryRepo.GetByID(ctx, entryID)
	if err != nil {
		return "", errors.New("entry not found")
	}
	return s.GetClinicIDForFormVersion(ctx, strconv.Itoa(entry.FormVersionID))
}

// GetClinicIDForFormVersion returns the clinic ID for the form that owns this version (for access control).
func (s *CustomFormEntryService) GetClinicIDForFormVersion(ctx context.Context, formVersionID string) (string, error) {
	vid, err := strconv.Atoi(formVersionID)
	if err != nil {
		return "", errors.New("invalid formVersionId")
	}
	version, err := s.versionRepo.GetByID(ctx, vid)
	if err != nil || version == nil {
		return "", errors.New("form version not found")
	}
	fm, err := s.formRepo.GetByID(ctx, version.FormID)
	if err != nil || fm == nil {
		return "", errors.New("form not found")
	}
	return fm.ClinicID, nil
}

// GetFormByIDForAccess returns the form for access control (handler needs clinic ID).
func (s *CustomFormEntryService) GetFormByIDForAccess(ctx context.Context, formID string) (*form.Form, error) {
	return s.formRepo.GetByID(ctx, formID)
}
