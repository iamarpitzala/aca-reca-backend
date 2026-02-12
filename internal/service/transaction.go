package service

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
	"github.com/iamarpitzala/aca-reca-backend/internal/repository"
	"github.com/jmoiron/sqlx"
)

type TransactionService struct {
	db *sqlx.DB
}

func NewTransactionService(db *sqlx.DB) *TransactionService {
	return &TransactionService{db: db}
}

// formFieldForMapping parses form fields JSON for id and accountId
type formFieldForMapping struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	AccountID *string `json:"accountId"`
}

// fieldCalc from entry calculations
type fieldCalc struct {
	FieldID     string  `json:"fieldId"`
	FieldName   string  `json:"fieldName"`
	BaseAmount  float64 `json:"baseAmount"`
	GstAmount   float64 `json:"gstAmount"`
	TotalAmount float64 `json:"totalAmount"`
}

// calculationsPayload from entry.Calculations
type calculationsPayload struct {
	FieldTotals []fieldCalc `json:"fieldTotals"`
}

// GenerateFromEntry creates transaction rows from a custom form entry using form field -> COA mapping.
// Deletes any existing transactions for the entry first (re-post).
func (s *TransactionService) GenerateFromEntry(ctx context.Context, entryID uuid.UUID) ([]domain.TransactionResponse, error) {
	entry, err := repository.GetCustomFormEntryByID(ctx, s.db, entryID)
	if err != nil {
		return nil, err
	}
	form, err := repository.GetCustomFormByID(ctx, s.db, entry.FormID)
	if err != nil {
		return nil, err
	}
	if entry.ClinicID != form.ClinicID {
		return nil, errors.New("entry and form clinic mismatch")
	}

	var fields []formFieldForMapping
	if len(form.Fields) > 0 {
		if err := json.Unmarshal(form.Fields, &fields); err != nil {
			return nil, errors.New("invalid form fields")
		}
	}
	fieldByID := make(map[string]*formFieldForMapping)
	for i := range fields {
		fieldByID[fields[i].ID] = &fields[i]
	}

	var calc calculationsPayload
	if len(entry.Calculations) > 0 {
		if err := json.Unmarshal(entry.Calculations, &calc); err != nil {
			return nil, errors.New("invalid entry calculations")
		}
	}

	// Delete existing transactions for this entry (re-post)
	if err := repository.DeleteTransactionsByEntryID(ctx, s.db, entryID); err != nil {
		return nil, err
	}

	ref := "#" + entryID.String()
	if len(ref) > 12 {
		ref = "#" + ref[len(ref)-8:]
	}
	date := entry.EntryDate
	now := time.Now()
	var out []domain.TransactionResponse

	for _, ft := range calc.FieldTotals {
		if ft.TotalAmount == 0 && ft.BaseAmount == 0 && ft.GstAmount == 0 {
			continue
		}
		field := fieldByID[ft.FieldID]
		var coaIDStr string
		if field != nil && field.AccountID != nil && *field.AccountID != "" {
			coaIDStr = *field.AccountID
		}
		if coaIDStr == "" {
			continue
		}

		coaID, err := uuid.Parse(coaIDStr)
		if err != nil {
			continue
		}
		assigned, err := repository.ClinicCOAAssigned(ctx, s.db, entry.ClinicID, coaID)
		if err != nil || !assigned {
			continue
		}
		aoc, err := repository.GetAOCByID(ctx, s.db, coaID)
		if err != nil || aoc == nil {
			continue
		}
		tax, err := repository.GetAccountTaxByID(ctx, s.db, aoc.AccountTaxID)
		if err != nil || tax == nil {
			continue
		}
		taxCategory := repository.TaxNameToCategory(tax.Name)

		t := &domain.Transaction{
			ID:              uuid.New(),
			ClinicID:        entry.ClinicID,
			SourceEntryID:   entryID,
			SourceFormID:    entry.FormID,
			FieldID:         &ft.FieldID,
			COAID:           coaID,
			AccountCode:     aoc.Code,
			AccountName:     aoc.Name,
			TaxCategory:     taxCategory,
			TransactionDate: date,
			Reference:       ref,
			Details:         form.Name + " - " + ft.FieldName,
			GrossAmount:     ft.TotalAmount,
			GSTAmount:       ft.GstAmount,
			NetAmount:       ft.BaseAmount,
			Status:          domain.TransactionStatusPosted,
			CreatedAt:       now,
			UpdatedAt:       now,
		}
		if err := repository.CreateTransaction(ctx, s.db, t); err != nil {
			return nil, err
		}
		out = append(out, *repository.TransactionToResponse(t))
	}

	return out, nil
}

// List returns transactions for a clinic with filters
func (s *TransactionService) List(ctx context.Context, clinicID uuid.UUID, f *domain.ListTransactionsFilters) (*domain.ListTransactionsResponse, error) {
	if f == nil {
		f = &domain.ListTransactionsFilters{Page: 1, Limit: 50, SortField: "date", SortDirection: "desc"}
	}
	list, total, err := repository.GetTransactionsByClinicID(ctx, s.db, clinicID, f)
	if err != nil {
		return nil, err
	}
	page := f.Page
	if page < 1 {
		page = 1
	}
	limit := f.Limit
	if limit < 1 {
		limit = 50
	}
	resp := make([]domain.TransactionResponse, len(list))
	for i := range list {
		resp[i] = *repository.TransactionToResponse(&list[i])
	}
	return &domain.ListTransactionsResponse{
		Transactions: resp,
		Total:        total,
		Page:         page,
		Limit:        limit,
		HasMore:      page*limit < total,
	}, nil
}

// ListByEntryID returns transactions for an entry
func (s *TransactionService) ListByEntryID(ctx context.Context, entryID uuid.UUID) ([]domain.TransactionResponse, error) {
	list, err := repository.GetTransactionsByEntryID(ctx, s.db, entryID)
	if err != nil {
		return nil, err
	}
	out := make([]domain.TransactionResponse, len(list))
	for i := range list {
		out[i] = *repository.TransactionToResponse(&list[i])
	}
	return out, nil
}

// GetFormFieldCOAMapping returns form fields with their COA mapping and clinic COA list
func (s *TransactionService) GetFormFieldCOAMapping(ctx context.Context, formID uuid.UUID, clinicID uuid.UUID) (*domain.FormFieldCOAMappingResponse, error) {
	form, err := repository.GetCustomFormByID(ctx, s.db, formID)
	if err != nil {
		return nil, err
	}
	if form.ClinicID != clinicID {
		return nil, errors.New("form does not belong to clinic")
	}
	var fields []formFieldForMapping
	if len(form.Fields) > 0 {
		if err := json.Unmarshal(form.Fields, &fields); err != nil {
			return nil, errors.New("invalid form fields")
		}
	}
	items := make([]domain.FormFieldCOAMappingItem, 0, len(fields))
	for _, f := range fields {
		items = append(items, domain.FormFieldCOAMappingItem{
			FieldID:   f.ID,
			FieldName: f.Name,
			AccountID: f.AccountID,
		})
	}
	coas, err := repository.ListClinicCOAs(ctx, s.db, clinicID)
	if err != nil {
		return nil, err
	}
	aocResponses := make([]domain.AOCResponse, len(coas))
	for i := range coas {
		aocResponses[i] = *coas[i].ToResponse()
	}
	return &domain.FormFieldCOAMappingResponse{
		FormID:     formID.String(),
		FormName:   form.Name,
		Fields:     items,
		ClinicCOAs: aocResponses,
	}, nil
}
