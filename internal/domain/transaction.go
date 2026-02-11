package domain

import (
	"time"

	"github.com/google/uuid"
)

// Transaction status
const (
	TransactionStatusPosted = "posted"
	TransactionStatusDraft  = "draft"
	TransactionStatusVoided = "voided"
)

// Tax category (matches frontend)
const (
	TaxCategoryGSTFreeIncome   = "gst_free_income"
	TaxCategoryGSTOnIncome     = "gst_on_income"
	TaxCategoryGSTOnExpenses   = "gst_on_expenses"
	TaxCategoryGSTFreeExpenses = "gst_free_expenses"
	TaxCategoryBASExcluded     = "bas_excluded"
)

// TaxNameToCategory maps tbl_account_tax.name to frontend tax category.
func TaxNameToCategory(name string) string {
	switch name {
	case "GST on Income":
		return TaxCategoryGSTOnIncome
	case "GST Free Income":
		return TaxCategoryGSTFreeIncome
	case "GST on Expenses":
		return TaxCategoryGSTOnExpenses
	case "GST Free Expenses":
		return TaxCategoryGSTFreeExpenses
	case "BAS Excluded":
		return TaxCategoryBASExcluded
	default:
		return TaxCategoryBASExcluded
	}
}

// Transaction DB model
type Transaction struct {
	ID             uuid.UUID  `db:"id"`
	ClinicID       uuid.UUID  `db:"clinic_id"`
	SourceEntryID  uuid.UUID  `db:"source_entry_id"`
	SourceFormID   uuid.UUID  `db:"source_form_id"`
	FieldID        *string    `db:"field_id"`
	COAID          uuid.UUID  `db:"coa_id"`
	AccountCode    string     `db:"account_code"`
	AccountName    string     `db:"account_name"`
	TaxCategory    string     `db:"tax_category"`
	TransactionDate time.Time `db:"transaction_date"`
	Reference      string     `db:"reference"`
	Details        string     `db:"details"`
	GrossAmount    float64    `db:"gross_amount"`
	GSTAmount      float64    `db:"gst_amount"`
	NetAmount      float64    `db:"net_amount"`
	Status         string     `db:"status"`
	CreatedAt      time.Time  `db:"created_at"`
	UpdatedAt      time.Time  `db:"updated_at"`
}

// TransactionResponse API response
type TransactionResponse struct {
	ID              string    `json:"id"`
	ClinicID        string    `json:"clinicId"`
	SourceEntryID   string    `json:"sourceEntryId"`
	SourceFormID    string    `json:"sourceFormId"`
	FieldID         *string   `json:"fieldId,omitempty"`
	AccountCode     string    `json:"accountCode"`
	AccountName     string    `json:"accountName"`
	TaxCategory     string    `json:"taxCategory"`
	Date            string    `json:"date"`
	Reference       string    `json:"reference"`
	Details         string    `json:"details"`
	GrossAmount     float64   `json:"grossAmount"`
	GSTAmount       float64   `json:"gstAmount"`
	NetAmount       float64   `json:"netAmount"`
	Status          string    `json:"status"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}

// ListTransactionsFilters for listing
type ListTransactionsFilters struct {
	Search       string
	TaxCategory  string
	Status       string
	DateFrom     string
	DateTo       string
	SortField    string
	SortDirection string
	Page         int
	Limit        int
}

// ListTransactionsResponse paginated response
type ListTransactionsResponse struct {
	Transactions []TransactionResponse `json:"transactions"`
	Total        int                   `json:"total"`
	Page         int                   `json:"page"`
	Limit        int                   `json:"limit"`
	HasMore      bool                  `json:"hasMore"`
}

// FormFieldCOAMappingItem for GET form field COA mapping
type FormFieldCOAMappingItem struct {
	FieldID   string  `json:"fieldId"`
	FieldName string  `json:"fieldName"`
	AccountID *string `json:"accountId,omitempty"`
}

// FormFieldCOAMappingResponse form fields with COA and clinic COA list
type FormFieldCOAMappingResponse struct {
	FormID     string                   `json:"formId"`
	FormName   string                   `json:"formName"`
	Fields     []FormFieldCOAMappingItem `json:"fields"`
	ClinicCOAs []AOCResponse            `json:"clinicCoas"`
}
