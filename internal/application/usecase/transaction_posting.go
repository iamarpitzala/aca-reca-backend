package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/application/port"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
)

// TransactionPostingService handles posting form entries to the general ledger (journal entries).
// Implements accounting best practice: source documents (entries) → posted transactions (double-entry).
type TransactionPostingService struct {
	entryRepo     port.CustomFormRepository
	txRepo        port.TransactionRepository
	clinicCOARepo port.ClinicCOARepository
	aocRepo       port.AOCRepository
}

func NewTransactionPostingService(
	entryRepo port.CustomFormRepository,
	txRepo port.TransactionRepository,
	clinicCOARepo port.ClinicCOARepository,
	aocRepo port.AOCRepository,
) *TransactionPostingService {
	return &TransactionPostingService{
		entryRepo:     entryRepo,
		txRepo:        txRepo,
		clinicCOARepo: clinicCOARepo,
		aocRepo:       aocRepo,
	}
}

type formFieldForMapping struct {
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	AccountID *string `json:"accountId"`
}

type fieldCalc struct {
	FieldID     string  `json:"fieldId"`
	FieldName   string  `json:"fieldName"`
	BaseAmount  float64 `json:"baseAmount"`
	GstAmount   float64 `json:"gstAmount"`
	TotalAmount float64 `json:"totalAmount"`
}

type calculationsPayload struct {
	FieldTotals []fieldCalc `json:"fieldTotals"`
}

// PostEntryToLedger creates journal entries (transactions) from a custom form entry.
// Deletes any existing posted transactions for the entry first (re-post).
func (s *TransactionPostingService) PostEntryToLedger(ctx context.Context, entryID uuid.UUID) ([]domain.TransactionResponse, error) {
	entry, err := s.entryRepo.GetEntryByID(ctx, entryID)
	if err != nil {
		return nil, err
	}
	form, err := s.entryRepo.GetByID(ctx, entry.FormID)
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

	if err := s.txRepo.DeleteByEntryID(ctx, entryID); err != nil {
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
		assigned, err := s.clinicCOARepo.Exists(ctx, entry.ClinicID, coaID)
		if err != nil || !assigned {
			continue
		}
		aoc, err := s.aocRepo.GetByID(ctx, coaID)
		if err != nil || aoc == nil {
			continue
		}
		tax, err := s.aocRepo.GetAccountTaxByID(ctx, aoc.AccountTaxID)
		if err != nil || tax == nil {
			continue
		}
		taxCategory := domain.TaxNameToCategory(tax.Name)

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
		if err := s.txRepo.Create(ctx, t); err != nil {
			return nil, err
		}
		out = append(out, *transactionToResponse(t))
	}

	return out, nil
}

// ListJournalEntries returns transactions for a clinic with filters.
func (s *TransactionPostingService) ListJournalEntries(ctx context.Context, clinicID uuid.UUID, f *domain.ListTransactionsFilters) (*domain.ListTransactionsResponse, error) {
	if f == nil {
		f = &domain.ListTransactionsFilters{Page: 1, Limit: 50, SortField: "date", SortDirection: "desc"}
	}
	list, total, err := s.txRepo.ListByClinicID(ctx, clinicID, f)
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
		resp[i] = *transactionToResponse(&list[i])
	}
	return &domain.ListTransactionsResponse{
		Transactions: resp,
		Total:        total,
		Page:         page,
		Limit:        limit,
		HasMore:      page*limit < total,
	}, nil
}

// ListJournalEntriesByEntry returns transactions for a single entry.
func (s *TransactionPostingService) ListJournalEntriesByEntry(ctx context.Context, entryID uuid.UUID) ([]domain.TransactionResponse, error) {
	list, err := s.txRepo.ListByEntryID(ctx, entryID)
	if err != nil {
		return nil, err
	}
	out := make([]domain.TransactionResponse, len(list))
	for i := range list {
		out[i] = *transactionToResponse(&list[i])
	}
	return out, nil
}

// GetFormFieldCOAMapping returns form fields with their COA mapping and clinic COA list.
func (s *TransactionPostingService) GetFormFieldCOAMapping(ctx context.Context, formID, clinicID uuid.UUID) (*domain.FormFieldCOAMappingResponse, error) {
	form, err := s.entryRepo.GetByID(ctx, formID)
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
	coas, err := s.aocRepo.ListAOCsAssignedToClinic(ctx, clinicID)
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

func transactionToResponse(t *domain.Transaction) *domain.TransactionResponse {
	dateStr := t.TransactionDate.Format("2006-01-02")
	return &domain.TransactionResponse{
		ID:            t.ID.String(),
		ClinicID:      t.ClinicID.String(),
		SourceEntryID: t.SourceEntryID.String(),
		SourceFormID:  t.SourceFormID.String(),
		FieldID:       t.FieldID,
		AccountCode:   t.AccountCode,
		AccountName:   t.AccountName,
		TaxCategory:   t.TaxCategory,
		Date:          dateStr,
		Reference:     t.Reference,
		Details:       t.Details,
		GrossAmount:   t.GrossAmount,
		GSTAmount:     t.GSTAmount,
		NetAmount:     t.NetAmount,
		Status:        t.Status,
		CreatedAt:     t.CreatedAt,
		UpdatedAt:     t.UpdatedAt,
	}
}
