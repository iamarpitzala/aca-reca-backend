package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/iamarpitzala/aca-reca-backend/internal/application/port"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain/form"
	"github.com/iamarpitzala/aca-reca-backend/util"
)

type CustomFormService struct {
	repo        port.CustomFormRepository
	fieldRepo   port.CustomFormFieldRepository
	sectionRepo port.CustomFormSectionRepository
	versionRepo port.CustomFormVersionRepository
	clinicRepo  port.ClinicRepository
}

func NewCustomFormService(
	repo port.CustomFormRepository,
	fieldRepo port.CustomFormFieldRepository,
	sectionRepo port.CustomFormSectionRepository,
	versionRepo port.CustomFormVersionRepository,
	clinicRepo port.ClinicRepository,
) *CustomFormService {
	return &CustomFormService{
		repo:        repo,
		fieldRepo:   fieldRepo,
		sectionRepo: sectionRepo,
		versionRepo: versionRepo,
		clinicRepo:  clinicRepo,
	}
}

func (s *CustomFormService) Create(ctx context.Context, req *form.FormRequest, userID string) (*form.FormResponse, error) {
	clinic, err := s.clinicRepo.GetByID(ctx, req.ClinicID)
	if err != nil || clinic == nil {
		return nil, errors.New("clinic not found")
	}
	// Use clinic's method type (NET/GROSS) as the form's calculation method
	if clinic.MethodType == "NET" || clinic.MethodType == "GROSS" {
		req.CalculationMethod = clinic.MethodType
	} else {
		req.CalculationMethod = "NET"
	}

	formInstance := &form.Form{}
	req.CreatedBy = userID
	req.CreatedAt = time.Now()
	formInstance.ToFormDB(req)
	if err := s.repo.Create(ctx, formInstance); err != nil {
		return nil, err
	}

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

	// Create sections and fields when provided (so "Service and facility" / "Other cost" etc. are stored)
	if len(req.Sections) > 0 {
		for order, sec := range req.Sections {
			sectionTypeID := sectionTypeCodeToID(sec.SectionTypeCode)
			if sectionTypeID == 0 {
				continue
			}
			secName := sec.Name
			if secName == "" {
				secName = sec.SectionTypeCode
			}
			newSec := &form.Section{
				FormVersionID: formVersion.ID,
				SectionTypeID: sectionTypeID,
				Name:          secName,
				Description:   strPtrIfNonEmpty(sec.Description),
				SectionOrder:  order,
				CreatedAt:     time.Now(),
				UpdatedAt:     time.Now(),
			}
			if err := s.sectionRepo.Create(ctx, newSec); err != nil {
				return nil, err
			}
			for _, f := range sec.Fields {
				if f.Label == "" {
					continue
				}
				payRespID := paymentResponsibilityCodeToID(f.PaymentResponsibility)
				newField := &form.Field{
					ID:                     form.NewFieldID(),
					SectionID:              newSec.ID,
					Label:                  f.Label,
					PaymentResponsibilityID: payRespID,
					CreatedAt:              time.Now(),
					UpdatedAt:              time.Now(),
				}
				if err := s.fieldRepo.Create(ctx, newField); err != nil {
					return nil, err
				}
			}
		}
	}

	var activeID *int
	activeID = &formVersion.ID
	return formInstance.ToFormResponse(activeID), nil
}

// sectionTypeCodeToID maps tbl_section_type code to id (migration order: 1=COLLECTION, 2=COST, 3=SERVICE_FACILITY, 4=OTHER_COST).
func sectionTypeCodeToID(code string) int {
	switch code {
	case "COLLECTION":
		return 1
	case "COST":
		return 2
	case "SERVICE_FACILITY":
		return 3
	case "OTHER_COST":
		return 4
	default:
		return 0
	}
}

// sectionTypeIDToCode maps tbl_section_type id to code (for GET form response).
func sectionTypeIDToCode(id int) string {
	switch id {
	case 1:
		return "COLLECTION"
	case 2:
		return "COST"
	case 3:
		return "SERVICE_FACILITY"
	case 4:
		return "OTHER_COST"
	default:
		return ""
	}
}

// paymentResponsibilityCodeToID maps tbl_payment_responsibility code to id (1=OWNER, 2=CLINIC).
func paymentResponsibilityCodeToID(code string) *int {
	var id int
	switch code {
	case "OWNER":
		id = 1
	case "CLINIC":
		id = 2
	default:
		return nil
	}
	return &id
}

func strPtrIfNonEmpty(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

func (s *CustomFormService) getActiveVersionID(ctx context.Context, formID string) *int {
	v, err := s.versionRepo.GetActiveByFormID(ctx, formID)
	if err != nil || v == nil {
		return nil
	}
	id := v.ID
	return &id
}

func (s *CustomFormService) GetByID(ctx context.Context, formID string) (*form.FormResponse, error) {
	f, err := s.repo.GetByID(ctx, formID)
	if err != nil {
		return nil, err
	}
	activeVersionID := s.getActiveVersionID(ctx, formID)
	resp := f.ToFormResponse(activeVersionID)
	if activeVersionID != nil {
		sections, _ := s.sectionRepo.GetByFormVersionID(ctx, *activeVersionID)
		resp.Sections = make([]form.SectionWithFieldsResponse, 0, len(sections))
		for _, sec := range sections {
			fields, _ := s.fieldRepo.GetBySectionID(ctx, sec.ID)
			fieldResponses := make([]form.FieldResponse, 0, len(fields))
			for _, fl := range fields {
				fieldResponses = append(fieldResponses, *fl.ToFieldResponse())
			}
			desc := ""
			if sec.Description != nil {
				desc = *sec.Description
			}
			resp.Sections = append(resp.Sections, form.SectionWithFieldsResponse{
				SectionTypeCode: sectionTypeIDToCode(sec.SectionTypeID),
				Name:            sec.Name,
				Description:     desc,
				SectionOrder:    sec.SectionOrder,
				Fields:          fieldResponses,
			})
		}
	}
	return resp, nil
}

func (s *CustomFormService) GetByClinicID(ctx context.Context, clinicID string) ([]form.FormResponse, error) {
	forms, err := s.repo.GetByClinicID(ctx, clinicID)
	if err != nil {
		return nil, err
	}
	out := make([]form.FormResponse, len(forms))
	for i, f := range forms {
		out[i] = *f.ToFormResponse(s.getActiveVersionID(ctx, f.ID))
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
		out[i] = *f.ToFormResponse(s.getActiveVersionID(ctx, f.ID))
	}
	return out, nil
}

func (s *CustomFormService) Update(ctx context.Context, f *form.Form) error {
	return s.repo.Update(ctx, f)
}

func (s *CustomFormService) UpdateByRequest(ctx context.Context, id string, req *form.FormRequest, userID string) (*form.FormResponse, error) {
	formInstance, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	formInstance.Name = req.Name
	formInstance.Description = req.Description
	formInstance.Status = req.Status
	formInstance.CalculationMethod = req.CalculationMethod
	formInstance.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, formInstance); err != nil {
		return nil, err
	}
	return formInstance.ToFormResponse(s.getActiveVersionID(ctx, id)), nil
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

	newVersionNum := latestVersion.Version + 1
	newVersionRecord := &form.FormVersion{}
	newVersionRecord.ToFormVersionDB(&form.FormVersionRequest{
		FormID:    id,
		Version:   newVersionNum,
		IsActive:  false, // insert inactive; SetActive() below ensures only one active per form (uniq_active_form_version)
		CreatedBy: userID,
		CreatedAt: time.Now(),
	})
	if err := s.versionRepo.Create(ctx, newVersionRecord); err != nil {
		return nil, err
	}

	// Copy sections from latest version
	oldSections, _ := s.sectionRepo.GetByFormVersionID(ctx, latestVersion.ID)
	sectionIDMap := make(map[int]int)
	for _, oldSec := range oldSections {
		newSec := &form.Section{
			FormVersionID: newVersionRecord.ID,
			SectionTypeID:  oldSec.SectionTypeID,
			Name:          oldSec.Name,
			Description:   oldSec.Description,
			SectionOrder:  oldSec.SectionOrder,
			CreatedAt:     time.Now(),
			UpdatedAt:     time.Now(),
		}
		if err := s.sectionRepo.Create(ctx, newSec); err != nil {
			return nil, err
		}
		sectionIDMap[oldSec.ID] = newSec.ID
	}

	// Copy fields (by section)
	for oldSecID, newSecID := range sectionIDMap {
		oldFields, _ := s.fieldRepo.GetBySectionID(ctx, oldSecID)
		for _, old := range oldFields {
			newField := &form.Field{
				ID:                     form.NewFieldID(),
				SectionID:              newSecID,
				Label:                  old.Label,
				PaymentResponsibilityID: old.PaymentResponsibilityID,
				CreatedAt:              time.Now(),
				UpdatedAt:              time.Now(),
			}
			if err := s.fieldRepo.Create(ctx, newField); err != nil {
				return nil, err
			}
		}
	}

	if err := s.versionRepo.SetActive(ctx, id, newVersionRecord.ID); err != nil {
		return nil, err
	}
	if err := s.repo.Publish(ctx, id); err != nil {
		return nil, err
	}
	f.Status = string(util.FormStatusPublished)
	f.UpdatedAt = time.Now()
	return f.ToFormResponse(&newVersionRecord.ID), nil
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
	return f.ToFormResponse(s.getActiveVersionID(ctx, id)), nil
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
	return f.ToFormResponse(s.getActiveVersionID(ctx, id)), nil
}

func (s *CustomFormService) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
