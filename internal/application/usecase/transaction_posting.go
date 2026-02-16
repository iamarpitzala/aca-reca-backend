package usecase

import (
	"context"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/application/port"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
	"github.com/iamarpitzala/aca-reca-backend/util"
)

// TransactionPostingService handles posting form entries to the general ledger (journal entries).
// Implements accounting best practice: source documents (entries) → posted transactions (double-entry).
type TransactionPostingService struct {
	entryRepo     port.FieldEntryRepository
	formRepo      port.CustomFormRepository
	fieldRepo     port.CustomFormFieldRepository
	versionRepo   port.CustomFormVersionRepository
	txRepo        port.TransactionRepository
	clinicCOARepo port.ClinicCOARepository
	aocRepo       port.AOCRepository
	calcEngine    port.EntryCalculationEngine
}

func NewTransactionPostingService(
	entryRepo port.FieldEntryRepository,
	formRepo port.CustomFormRepository,
	fieldRepo port.CustomFormFieldRepository,
	versionRepo port.CustomFormVersionRepository,
	txRepo port.TransactionRepository,
	clinicCOARepo port.ClinicCOARepository,
	aocRepo port.AOCRepository,
	calcEngine port.EntryCalculationEngine,
) *TransactionPostingService {
	return &TransactionPostingService{
		entryRepo:     entryRepo,
		formRepo:      formRepo,
		fieldRepo:     fieldRepo,
		versionRepo:   versionRepo,
		txRepo:        txRepo,
		clinicCOARepo: clinicCOARepo,
		aocRepo:       aocRepo,
		calcEngine:    calcEngine,
	}
}

type formFieldForMapping struct {
	ID                   string  `json:"id"`
	Name                 string  `json:"name"`
	AccountID            *string `json:"accountId"`
	AmountInterpretation string  `json:"amountInterpretation"` // "net" | "tax_only" - what the field value represents
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
	// Get the first field entry to get form ID and version ID
	entry, err := s.entryRepo.GetByID(ctx, entryID)
	if err != nil {
		return nil, err
	}
	form, err := s.formRepo.GetByID(ctx, entry.FormID)
	if err != nil {
		return nil, err
	}

	// Get all field entries for this entry (grouped by created_at)
	allEntries, err := s.entryRepo.GetByFormID(ctx, entry.FormID)
	if err != nil {
		return nil, err
	}

	// Filter entries created at the same timestamp (within 1 second) and same creator
	groupedEntries := make([]domain.FieldEntry, 0)
	targetTime := entry.CreatedAt.Unix()
	for _, e := range allEntries {
		if e.CreatedAt.Unix() == targetTime && e.CreatedBy == entry.CreatedBy {
			groupedEntries = append(groupedEntries, e)
		}
	}

	if len(groupedEntries) == 0 {
		return []domain.TransactionResponse{}, nil
	}

	// Get form version and fields
	formVersionID := groupedEntries[0].FormVersionID
	formFields, err := s.fieldRepo.GetByFormVersionID(ctx, formVersionID)
	if err != nil {
		return nil, err
	}

	// Build field mapping with account IDs
	fields := make([]formFieldForMapping, 0, len(formFields))
	fieldByID := make(map[string]*formFieldForMapping)
	for _, f := range formFields {
		// Parse metadata to get accountId and amountInterpretation
		var metadata map[string]interface{}
		var accountID *string
		amountInterpretation := "net"
		if len(f.Metadata) > 0 {
			_ = json.Unmarshal(f.Metadata, &metadata)
			if val, ok := metadata["accountId"].(string); ok && val != "" {
				accountID = &val
			}
			if val, ok := metadata["amountInterpretation"].(string); ok {
				amountInterpretation = strings.ToLower(strings.TrimSpace(val))
			}
		}
		// If no accountId in metadata, use coa_id from field
		if accountID == nil {
			coaIDStr := f.CoaID.String()
			accountID = &coaIDStr
		}

		field := formFieldForMapping{
			ID:                   f.ID.String(),
			Name:                 f.Label,
			AccountID:            accountID,
			AmountInterpretation: amountInterpretation,
		}
		fields = append(fields, field)
		fieldByID[field.ID] = &fields[len(fields)-1]
	}

	// Build entry values for calculation
	type entryValue struct {
		FieldID         string      `json:"fieldId"`
		FieldName       string      `json:"fieldName"`
		Value           interface{} `json:"value"`
		ManualGstAmount *float64    `json:"manualGstAmount"`
	}
	entryValues := make([]entryValue, 0, len(groupedEntries))
	for _, e := range groupedEntries {
		fieldName := ""
		for _, f := range formFields {
			if f.ID == e.CustomFormFieldID {
				fieldName = f.Label
				break
			}
		}
		entryValues = append(entryValues, entryValue{
			FieldID:   e.CustomFormFieldID.String(),
			FieldName: fieldName,
			Value:     e.Value,
		})
	}

	// Convert form fields to calcField format for calculation engine
	type calcField struct {
		ID             string `json:"id"`
		Name           string `json:"name"`
		Type           string `json:"type"`
		Section        string `json:"section"`
		IncludeInTotal bool   `json:"includeInTotal"`
		GstConfig      *struct {
			Enabled bool    `json:"enabled"`
			Rate    float64 `json:"rate"`
			Type    string  `json:"type"`
		} `json:"gstConfig"`
		PaymentResp string `json:"paymentResponsibility"`
	}

	calcFields := make([]calcField, 0, len(formFields))
	for _, f := range formFields {
		var metadata map[string]interface{}
		includeInTotal := true
		if len(f.Metadata) > 0 {
			_ = json.Unmarshal(f.Metadata, &metadata)
			if val, ok := metadata["includeInTotal"].(bool); ok {
				includeInTotal = val
			}
		}

		cf := calcField{
			ID:             f.ID.String(),
			Name:           f.Label,
			Type:           strings.ToLower(f.FieldType),
			Section:        strings.ToLower(f.Section),
			IncludeInTotal: includeInTotal,
		}

		if f.GSTConfig {
			gstRate := 0.0
			if f.GSTRate != nil {
				gstRate = *f.GSTRate
			}
			gstType := strings.ToLower(f.GSTType)
			if gstType == "" {
				gstType = "exclusive"
			}
			cf.GstConfig = &struct {
				Enabled bool    `json:"enabled"`
				Rate    float64 `json:"rate"`
				Type    string  `json:"type"`
			}{
				Enabled: true,
				Rate:    gstRate,
				Type:    gstType,
			}
		}
		calcFields = append(calcFields, cf)
	}

	// Marshal to JSON for calculation engine
	fieldsJSON, err := json.Marshal(calcFields)
	if err != nil {
		return nil, err
	}
	valuesJSON, err := json.Marshal(entryValues)
	if err != nil {
		return nil, err
	}

	// Calculate totals using calculation engine
	calculationsJSON, err := s.calcEngine.RunEntryCalculation(
		fieldsJSON,
		form.FormType,
		form.CalculationMethod,
		nil,   // serviceFacilityFeePercent
		false, // outworkEnabled
		nil,   // outworkRatePercent
		valuesJSON,
		nil, // deductions
	)
	if err != nil {
		return nil, err
	}

	// Parse calculations
	var calc calculationsPayload
	if err := json.Unmarshal(calculationsJSON, &calc); err != nil {
		return nil, err
	}

	// Delete existing transactions for this entry
	if err := s.txRepo.DeleteByEntryID(ctx, entryID); err != nil {
		return nil, err
	}

	ref := "#" + entryID.String()
	if len(ref) > 12 {
		ref = "#" + ref[len(ref)-8:]
	}
	date := entry.CreatedAt
	now := time.Now()
	var out []domain.TransactionResponse

	// Create transactions for each field total
	for _, ft := range calc.FieldTotals {
		if ft.TotalAmount == 0 && ft.BaseAmount == 0 && ft.GstAmount == 0 {
			continue
		}
		field := fieldByID[ft.FieldID]
		if field == nil {
			continue
		}

		var coaIDStr string
		if field.AccountID != nil && *field.AccountID != "" {
			coaIDStr = *field.AccountID
		}
		if coaIDStr == "" {
			continue
		}

		coaID, err := uuid.Parse(coaIDStr)
		if err != nil {
			continue
		}
		assigned, err := s.clinicCOARepo.Exists(ctx, form.ClinicID, coaID)
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

		grossAmt := ft.TotalAmount
		gstAmt := ft.GstAmount
		netAmt := ft.BaseAmount
		interp := strings.ToLower(strings.TrimSpace(field.AmountInterpretation))
		switch interp {
		case "tax_only":
			grossAmt = ft.GstAmount
			gstAmt = ft.GstAmount
			netAmt = 0
		case "net":
			fallthrough
		default:
			// Net is primary; gross = net + gst (already in field totals)
			grossAmt = ft.TotalAmount
			netAmt = ft.BaseAmount
			gstAmt = ft.GstAmount
		}

		t := &domain.TransactionWithLedger{
			Transaction: domain.Transaction{
				ID:            uuid.New(),
				ClinicID:      form.ClinicID,
				SourceEntryID: entryID,
				CreatedAt:     now,
				UpdatedAt:     now,
			},
			SourceFormID:    entry.FormID,
			FieldID:         &ft.FieldID,
			COAID:           coaID,
			AccountCode:     aoc.Code,
			AccountName:     aoc.Name,
			TaxCategory:     taxCategory,
			TransactionDate: date,
			Reference:       ref,
			Details:         form.Name + " - " + field.Name,
			GrossAmount:     grossAmt,
			GSTAmount:       gstAmt,
			NetAmount:       netAmt,
			Status:          util.TransactionStatusPosted,
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
	form, err := s.formRepo.GetByID(ctx, formID)
	if err != nil {
		return nil, err
	}
	if form.ClinicID != clinicID {
		return nil, errors.New("form does not belong to clinic")
	}

	// Get latest version and its fields
	version, err := s.versionRepo.GetLatestByFormID(ctx, formID)
	if err != nil {
		return &domain.FormFieldCOAMappingResponse{
			FormID:     formID.String(),
			FormName:   form.Name,
			Fields:     []domain.FormFieldCOAMappingItem{},
			ClinicCOAs: []domain.AOCResponse{},
		}, nil
	}

	formFields, err := s.fieldRepo.GetByFormVersionID(ctx, version.ID)
	if err != nil {
		return nil, err
	}

	items := make([]domain.FormFieldCOAMappingItem, 0, len(formFields))
	for _, f := range formFields {
		// Parse metadata to get accountId and amountInterpretation
		var metadata map[string]interface{}
		var accountID *string
		amountInterpretation := "net"
		if len(f.Metadata) > 0 {
			_ = json.Unmarshal(f.Metadata, &metadata)
			if val, ok := metadata["accountId"].(string); ok && val != "" {
				accountID = &val
			}
			if val, ok := metadata["amountInterpretation"].(string); ok {
				amountInterpretation = strings.ToLower(strings.TrimSpace(val))
			}
		}
		// If no accountId in metadata, use coa_id from field
		if accountID == nil {
			coaIDStr := f.CoaID.String()
			accountID = &coaIDStr
		}

		interp := strings.ToLower(strings.TrimSpace(amountInterpretation))
		if interp != "net" && interp != "tax_only" {
			interp = "net"
		}
		items = append(items, domain.FormFieldCOAMappingItem{
			FieldID:              f.ID.String(),
			FieldName:            f.Label,
			AccountID:            accountID,
			AmountInterpretation: interp,
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

func transactionToResponse(t *domain.TransactionWithLedger) *domain.TransactionResponse {
	dateStr := t.TransactionDate.Format("2006-01-02")
	resp := &domain.TransactionResponse{
		ID:            t.ID.String(),
		TransactionID: t.ID.String(),
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
	if t.COAID != uuid.Nil {
		resp.COAID = t.COAID.String()
	}
	return resp
}
