package domain

import (
	"time"

	"github.com/google/uuid"
)

// BASSnapshot represents tbl_bas_snapshot (header)
type BASSnapshot struct {
	ID         uuid.UUID  `db:"id" json:"id"`
	ClinicID   uuid.UUID  `db:"clinic_id" json:"clinicId"`
	QuarterID  *uuid.UUID `db:"quarter_id" json:"quarterId,omitempty"`
	PeriodStart time.Time `db:"period_start" json:"periodStart"`
	PeriodEnd   time.Time `db:"period_end" json:"periodEnd"`
	PeriodType  string    `db:"period_type" json:"periodType"`

	// GST on Sales (G1–G9)
	G1TotalSales            float64 `db:"g1_total_sales" json:"g1TotalSales"`
	G2ExportSales           float64 `db:"g2_export_sales" json:"g2ExportSales"`
	G3GSTFreeSales          float64 `db:"g3_gst_free_sales" json:"g3GSTFreeSales"`
	G4InputTaxedSales       float64 `db:"g4_input_taxed_sales" json:"g4InputTaxedSales"`
	G5G2G3G4                float64 `db:"g5_g2_g3_g4" json:"g5G2G3G4"`
	G6TotalTaxableSales     float64 `db:"g6_total_taxable_sales" json:"g6TotalTaxableSales"`
	G7Adjustments           float64 `db:"g7_adjustments" json:"g7Adjustments"`
	G8TotalTaxableSupplies  float64 `db:"g8_total_taxable_supplies" json:"g8TotalTaxableSupplies"`
	G9GSTOnSales            float64 `db:"g9_gst_on_sales" json:"g9GSTOnSales"`

	// GST on Purchases (G10–G20)
	G10CapitalPurchases          float64 `db:"g10_capital_purchases" json:"g10CapitalPurchases"`
	G11NonCapitalPurchases       float64 `db:"g11_non_capital_purchases" json:"g11NonCapitalPurchases"`
	G12G10G11                    float64 `db:"g12_g10_g11" json:"g12G10G11"`
	G13PurchasesForInputTaxed    float64 `db:"g13_purchases_for_input_taxed_sales" json:"g13PurchasesForInputTaxedSales"`
	G14PurchasesWithoutGST       float64 `db:"g14_purchases_without_gst" json:"g14PurchasesWithoutGST"`
	G15PrivateUse                float64 `db:"g15_private_use" json:"g15PrivateUse"`
	G16G13G14G15                 float64 `db:"g16_g13_g14_g15" json:"g16G13G14G15"`
	G17TotalCreditablePurchases  float64 `db:"g17_total_creditable_purchases" json:"g17TotalCreditablePurchases"`
	G18Adjustments               float64 `db:"g18_adjustments" json:"g18Adjustments"`
	G19TotalCreditableAcquisitions float64 `db:"g19_total_creditable_acquisitions" json:"g19TotalCreditableAcquisitions"`
	G20GSTOnPurchases            float64 `db:"g20_gst_on_purchases" json:"g20GSTOnPurchases"`

	// Summary Labels
	Label1AGSTOnSales     float64 `db:"label_1a_gst_on_sales" json:"label1AGSTOnSales"`
	Label1BGSTOnPurchases float64 `db:"label_1b_gst_on_purchases" json:"label1BGSTOnPurchases"`

	// PAYG Withholding (W1–W4)
	W1TotalSalaryWages    float64 `db:"w1_total_salary_wages" json:"w1TotalSalaryWages"`
	W2AmountsWithheld     float64 `db:"w2_amounts_withheld" json:"w2AmountsWithheld"`
	W3OtherAmountsWithheld float64 `db:"w3_other_amounts_withheld" json:"w3OtherAmountsWithheld"`
	W4TotalWithheld       float64 `db:"w4_total_withheld" json:"w4TotalWithheld"`

	// PAYG Instalments (T1–T4)
	T1InstalmentIncome float64 `db:"t1_instalment_income" json:"t1InstalmentIncome"`
	T2InstalmentRate   float64 `db:"t2_instalment_rate" json:"t2InstalmentRate"`
	T3NewVariedRate    float64 `db:"t3_new_varied_rate" json:"t3NewVariedRate"`
	T4InstalmentAmount float64 `db:"t4_instalment_amount" json:"t4InstalmentAmount"`

	// Net Amounts
	NetGSTPayable    float64 `db:"net_gst_payable" json:"netGSTPayable"`
	TotalAmountOwing float64 `db:"total_amount_owing" json:"totalAmountOwing"`

	// Lifecycle
	Status      string     `db:"status" json:"status"`
	GeneratedBy uuid.UUID  `db:"generated_by" json:"generatedBy"`
	FinalisedAt *time.Time `db:"finalised_at" json:"finalisedAt,omitempty"`
	FinalisedBy *uuid.UUID `db:"finalised_by" json:"finalisedBy,omitempty"`
	LockedAt    *time.Time `db:"locked_at" json:"lockedAt,omitempty"`
	LockedBy    *uuid.UUID `db:"locked_by" json:"lockedBy,omitempty"`
	Notes       *string    `db:"notes" json:"notes,omitempty"`

	// Flexible JSONB for raw computation / audit data
	SnapshotData interface{} `db:"snapshot_data" json:"snapshotData,omitempty"`

	CreatedAt time.Time  `db:"created_at" json:"createdAt"`
	UpdatedAt time.Time  `db:"updated_at" json:"updatedAt"`
	DeletedAt *time.Time `db:"deleted_at" json:"deletedAt,omitempty"`
}

// BASSnapshotLine represents tbl_bas_snapshot_line (line item)
type BASSnapshotLine struct {
	ID              uuid.UUID `db:"id" json:"id"`
	BASSnapshotID   uuid.UUID `db:"bas_snapshot_id" json:"basSnapshotId"`
	COAID           uuid.UUID `db:"coa_id" json:"coaId"`
	AccountCode     string    `db:"account_code" json:"accountCode"`
	AccountName     string    `db:"account_name" json:"accountName"`
	AccountTaxName  string    `db:"account_tax_name" json:"accountTaxName"`
	BASLabel        string    `db:"bas_label" json:"basLabel"`
	BaseAmount      float64   `db:"base_amount" json:"baseAmount"`
	GSTAmount       float64   `db:"gst_amount" json:"gstAmount"`
	TotalAmount     float64   `db:"total_amount" json:"totalAmount"`
	TransactionCount int      `db:"transaction_count" json:"transactionCount"`
	CreatedAt       time.Time `db:"created_at" json:"createdAt"`
}

// BASSnapshotWithLines combines header + line items for API responses
type BASSnapshotWithLines struct {
	Snapshot BASSnapshot       `json:"snapshot"`
	Lines    []BASSnapshotLine `json:"lines"`
}

// ---- Request DTOs ----

// GenerateBASRequest is the request body for generating a BAS snapshot
type GenerateBASRequest struct {
	PeriodStart string `json:"periodStart" validate:"required"`
	PeriodEnd   string `json:"periodEnd" validate:"required"`
	PeriodType  string `json:"periodType" validate:"required,oneof=QUARTERLY ANNUALLY"`
}

// UpdateBASSnapshotRequest is the request body for updating a BAS snapshot
type UpdateBASSnapshotRequest struct {
	PeriodStart *string  `json:"periodStart,omitempty"`
	PeriodEnd   *string  `json:"periodEnd,omitempty"`
	PeriodType  *string  `json:"periodType,omitempty"`
	G1TotalSales            *float64 `json:"g1TotalSales,omitempty"`
	G2ExportSales           *float64 `json:"g2ExportSales,omitempty"`
	G3GSTFreeSales          *float64 `json:"g3GSTFreeSales,omitempty"`
	G10CapitalPurchases     *float64 `json:"g10CapitalPurchases,omitempty"`
	G11NonCapitalPurchases  *float64 `json:"g11NonCapitalPurchases,omitempty"`
	Label1AGSTOnSales       *float64 `json:"label1AGSTOnSales,omitempty"`
	Label1BGSTOnPurchases   *float64 `json:"label1BGSTOnPurchases,omitempty"`
	NetGSTPayable           *float64 `json:"netGSTPayable,omitempty"`
	Notes                   *string  `json:"notes,omitempty"`
	SnapshotData            interface{} `json:"snapshotData,omitempty"`
}

// ---- Response DTOs ----

// BASSnapshotListResponse wraps a list of snapshots for API output
type BASSnapshotListResponse struct {
	Snapshots []BASSnapshot `json:"snapshots"`
	Total     int           `json:"total"`
}

// BASLedgerAggRow is used internally to aggregate ledger data for BAS generation
type BASLedgerAggRow struct {
	COAID            uuid.UUID `db:"coa_id"`
	AccountCode      string    `db:"account_code"`
	AccountName      string    `db:"account_name"`
	AccountType      string    `db:"account_type"`
	AccountTaxName   string    `db:"account_tax_name"`
	DebitTotal       float64   `db:"debit_total"`
	CreditTotal      float64   `db:"credit_total"`
	NetAmount        float64   `db:"net_amount"`
	GSTAmount        float64   `db:"gst_amount"`
	TransactionCount int       `db:"transaction_count"`
}
