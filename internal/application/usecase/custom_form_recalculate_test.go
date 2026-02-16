package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/application/port"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
	"github.com/iamarpitzala/aca-reca-backend/util"
)

// mockRecalcRepo returns a fixed entry and form, and records UpdateEntry calls.
type mockRecalcRepo struct {
	entry *domain.CustomFormEntry
	form  *domain.CustomForm
	// lastUpdateEntry is the entry passed to UpdateEntry
	lastUpdateEntry *domain.CustomFormEntry
	updateErr       error
}

func (m *mockRecalcRepo) Create(ctx context.Context, form *domain.CustomForm) error              { return nil }
func (m *mockRecalcRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.CustomForm, error) {
	if m.form != nil && m.form.ID == id {
		return m.form, nil
	}
	return nil, errors.New("form not found")
}
func (m *mockRecalcRepo) GetByClinicID(ctx context.Context, clinicID uuid.UUID) ([]domain.CustomForm, error) {
	return nil, nil
}
func (m *mockRecalcRepo) GetPublishedByClinicID(ctx context.Context, clinicID uuid.UUID) ([]domain.CustomForm, error) {
	return nil, nil
}
func (m *mockRecalcRepo) Update(ctx context.Context, form *domain.CustomForm) error   { return nil }
func (m *mockRecalcRepo) Publish(ctx context.Context, id uuid.UUID) error              { return nil }
func (m *mockRecalcRepo) Unpublish(ctx context.Context, id uuid.UUID) error           { return nil }
func (m *mockRecalcRepo) Archive(ctx context.Context, id uuid.UUID) error              { return nil }
func (m *mockRecalcRepo) Delete(ctx context.Context, id uuid.UUID) error               { return nil }
func (m *mockRecalcRepo) CreateEntry(ctx context.Context, entry *domain.CustomFormEntry) error { return nil }
func (m *mockRecalcRepo) GetEntryByID(ctx context.Context, id uuid.UUID) (*domain.CustomFormEntry, error) {
	if m.entry != nil && m.entry.ID == id {
		return m.entry, nil
	}
	return nil, errors.New("entry not found")
}
func (m *mockRecalcRepo) GetEntriesByFormID(ctx context.Context, formID uuid.UUID) ([]domain.CustomFormEntry, error) {
	return nil, nil
}
func (m *mockRecalcRepo) GetEntriesByClinicID(ctx context.Context, clinicID uuid.UUID) ([]domain.CustomFormEntry, error) {
	return nil, nil
}
func (m *mockRecalcRepo) GetEntriesByQuarter(ctx context.Context, clinicID, quarterID uuid.UUID) ([]domain.CustomFormEntry, error) {
	return nil, nil
}
func (m *mockRecalcRepo) UpdateEntry(ctx context.Context, entry *domain.CustomFormEntry) error {
	m.lastUpdateEntry = entry
	return m.updateErr
}
func (m *mockRecalcRepo) DeleteEntry(ctx context.Context, id uuid.UUID) error { return nil }

type mockRecalcClinicRepo struct{}

func (m *mockRecalcClinicRepo) Create(ctx context.Context, clinic *domain.Clinic) error { return nil }
func (m *mockRecalcClinicRepo) GetByID(ctx context.Context, id uuid.UUID) (*domain.Clinic, error) {
	return nil, nil
}
func (m *mockRecalcClinicRepo) Update(ctx context.Context, clinic *domain.Clinic) error { return nil }
func (m *mockRecalcClinicRepo) Delete(ctx context.Context, id uuid.UUID) error         { return nil }
func (m *mockRecalcClinicRepo) List(ctx context.Context) ([]domain.Clinic, error)       { return nil, nil }
func (m *mockRecalcClinicRepo) GetByABN(ctx context.Context, abnNumber string) (*domain.Clinic, error) {
	return nil, nil
}

// mockRecalcEngine returns fixed calculations JSON.
type mockRecalcEngine struct {
	out []byte
	err error
}

func (m *mockRecalcEngine) RunEntryCalculation(
	formFieldsJSON []byte,
	formType string,
	formCalculationMethod string,
	serviceFacilityFeePercent *float64,
	outworkEnabled bool,
	outworkRatePercent *float64,
	valuesJSON, deductionsJSON []byte,
) ([]byte, error) {
	return m.out, m.err
}

func TestRecalculateEntry_SetsCalculationsAndCallsUpdate(t *testing.T) {
	entryID := uuid.New()
	formID := uuid.New()
	clinicID := uuid.New()
	userID := uuid.New()
	now := time.Now()

	entry := &domain.CustomFormEntry{
		ID:             entryID,
		FormID:         formID,
		ClinicID:       clinicID,
		Values:         []byte(`[{"fieldId":"f1","value":1000}]`),
		Deductions:     []byte(`{}`),
		Calculations:   []byte(`{}`),
		CreatedBy:      userID,
		CreatedAt:      now,
		UpdatedAt:      now,
	}

	form := &domain.CustomForm{
		ID:                 formID,
		ClinicID:            clinicID,
		Name:                "Test",
		CalculationMethod:   util.MethodTypeGross,
		FormType:            util.FormTypeIncome,
		Fields:              []byte(`[{"id":"f1","name":"Fee","type":"number","section":"income","includeInTotal":true,"gstConfig":{"enabled":true,"rate":10,"type":"exclusive"}}]`),
		OutworkEnabled:      false,
	}

	calcOut := map[string]interface{}{
		"totalBaseAmount": 1000,
		"totalAmount":     1100,
		"serviceFeeBase":  500,
		"totalServiceFee": 550,
	}
	calcJSON, _ := json.Marshal(calcOut)

	repo := &mockRecalcRepo{entry: entry, form: form}
	clinicRepo := &mockRecalcClinicRepo{}
	engine := &mockRecalcEngine{out: calcJSON}

	svc := NewCustomFormService(repo, clinicRepo, engine)
	ctx := context.Background()

	resp, err := svc.RecalculateEntry(ctx, entryID)
	if err != nil {
		t.Fatalf("RecalculateEntry: %v", err)
	}
	if resp == nil {
		t.Fatal("expected non-nil response")
	}
	if repo.lastUpdateEntry == nil {
		t.Fatal("UpdateEntry was not called")
	}
	if len(repo.lastUpdateEntry.Calculations) == 0 {
		t.Error("expected entry.Calculations to be set before UpdateEntry")
	}
	var updated map[string]interface{}
	if err := json.Unmarshal(repo.lastUpdateEntry.Calculations, &updated); err != nil {
		t.Fatalf("unmarshal updated calculations: %v", err)
	}
	if updated["serviceFeeBase"].(float64) != 500 {
		t.Errorf("calculations.serviceFeeBase want 500, got %v", updated["serviceFeeBase"])
	}
}

func TestRecalculateEntry_EntryNotFound(t *testing.T) {
	repo := &mockRecalcRepo{entry: nil, form: nil}
	svc := NewCustomFormService(repo, &mockRecalcClinicRepo{}, &mockRecalcEngine{})
	ctx := context.Background()

	_, err := svc.RecalculateEntry(ctx, uuid.New())
	if err == nil {
		t.Error("expected error when entry not found")
	}
}

func TestRecalculateEntry_CalculationError(t *testing.T) {
	entryID := uuid.New()
	formID := uuid.New()
	entry := &domain.CustomFormEntry{ID: entryID, FormID: formID, Values: []byte("[]")}
	form := &domain.CustomForm{ID: formID, Fields: []byte("[]"), CalculationMethod: util.MethodTypeGross, FormType: util.FormTypeIncome}

	repo := &mockRecalcRepo{entry: entry, form: form}
	engine := &mockRecalcEngine{err: errors.New("calc failed")}
	svc := NewCustomFormService(repo, &mockRecalcClinicRepo{}, engine)
	ctx := context.Background()

	_, err := svc.RecalculateEntry(ctx, entryID)
	if err == nil {
		t.Error("expected error when calculation fails")
	}
	if repo.lastUpdateEntry != nil {
		t.Error("UpdateEntry should not be called when calculation fails")
	}
}

// Ensure mockRecalcRepo implements port.CustomFormRepository
var _ port.CustomFormRepository = (*mockRecalcRepo)(nil)
