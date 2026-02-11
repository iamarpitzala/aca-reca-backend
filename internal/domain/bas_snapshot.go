package domain

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// BASPeriodType represents the type of BAS period
type BASPeriodType string

const (
	BASPeriodTypeQuarterly BASPeriodType = "QUARTERLY"
	BASPeriodTypeAnnually   BASPeriodType = "ANNUALLY"
)

// BASStatus represents the status of a BAS
type BASStatus string

const (
	BASStatusDraft     BASStatus = "DRAFT"
	BASStatusFinalised BASStatus = "FINALISED"
	BASStatusLocked    BASStatus = "LOCKED"
)

// BASSnapshot represents a finalised BAS record
type BASSnapshot struct {
	ID                    uuid.UUID   `db:"id" json:"id"`
	ClinicID              uuid.UUID   `db:"clinic_id" json:"clinicId"`
	PeriodStart           time.Time   `db:"period_start" json:"periodStart"`
	PeriodEnd             time.Time   `db:"period_end" json:"periodEnd"`
	PeriodType            BASPeriodType `db:"period_type" json:"periodType"`
	G1TotalSales          float64     `db:"g1_total_sales" json:"g1TotalSales"`
	G2ExportSales          float64     `db:"g2_export_sales" json:"g2ExportSales"`
	G3GSTFreeSales         float64     `db:"g3_gst_free_sales" json:"g3GSTFreeSales"`
	G10CapitalPurchases    float64     `db:"g10_capital_purchases" json:"g10CapitalPurchases"`
	G11NonCapitalPurchases float64     `db:"g11_non_capital_purchases" json:"g11NonCapitalPurchases"`
	Label1AGSTOnSales      float64     `db:"label_1a_gst_on_sales" json:"label1AGSTOnSales"`
	Label1BGSTOnPurchases  float64     `db:"label_1b_gst_on_purchases" json:"label1BGSTOnPurchases"`
	NetGSTPayable          float64     `db:"net_gst_payable" json:"netGSTPayable"`
	Status                 BASStatus   `db:"status" json:"status"`
	FinalisedAt            *time.Time `db:"finalised_at" json:"finalisedAt"`
	FinalisedBy            *uuid.UUID `db:"finalised_by" json:"finalisedBy"`
	SnapshotData           json.RawMessage `db:"snapshot_data" json:"snapshotData"` // JSONB in DB
	CreatedAt              time.Time   `db:"created_at" json:"createdAt"`
	UpdatedAt              time.Time   `db:"updated_at" json:"updatedAt"`
	DeletedAt              *time.Time `db:"deleted_at" json:"deletedAt"`
}

// BASSnapshotRequest represents the request to create/finalise a BAS
type BASSnapshotRequest struct {
	PeriodStart           time.Time                `json:"periodStart" validate:"required"`
	PeriodEnd             time.Time                `json:"periodEnd" validate:"required"`
	PeriodType            BASPeriodType            `json:"periodType" validate:"required,oneof=QUARTERLY ANNUALLY"`
	G1TotalSales          float64                  `json:"g1TotalSales"`
	G2ExportSales         float64                  `json:"g2ExportSales"`
	G3GSTFreeSales         float64                  `json:"g3GSTFreeSales"`
	G10CapitalPurchases    float64                 `json:"g10CapitalPurchases"`
	G11NonCapitalPurchases float64                 `json:"g11NonCapitalPurchases"`
	Label1AGSTOnSales      float64                 `json:"label1AGSTOnSales"`
	Label1BGSTOnPurchases   float64                `json:"label1BGSTOnPurchases"`
	NetGSTPayable          float64                 `json:"netGSTPayable"`
	SnapshotData           json.RawMessage         `json:"snapshotData"`
}

// ToBASSnapshot converts a request to a domain model
func (r *BASSnapshotRequest) ToBASSnapshot(clinicID uuid.UUID) *BASSnapshot {
	return &BASSnapshot{
		ClinicID:              clinicID,
		PeriodStart:           r.PeriodStart,
		PeriodEnd:             r.PeriodEnd,
		PeriodType:            r.PeriodType,
		G1TotalSales:          r.G1TotalSales,
		G2ExportSales:         r.G2ExportSales,
		G3GSTFreeSales:        r.G3GSTFreeSales,
		G10CapitalPurchases:   r.G10CapitalPurchases,
		G11NonCapitalPurchases: r.G11NonCapitalPurchases,
		Label1AGSTOnSales:     r.Label1AGSTOnSales,
		Label1BGSTOnPurchases: r.Label1BGSTOnPurchases,
		NetGSTPayable:         r.NetGSTPayable,
		Status:                BASStatusDraft,
		SnapshotData:          r.SnapshotData,
	}
}
