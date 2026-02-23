package port

import (
	"context"
	"time"

	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
)

type PnlReportRepository interface {
	Create(ctx context.Context, report *domain.PnlReport) error
	CreateLines(ctx context.Context, lines []domain.PnlReportLine) error
	GetByID(ctx context.Context, id string) (*domain.PnlReport, error)
	GetLinesByReportID(ctx context.Context, reportID string) ([]domain.PnlReportLine, error)
	GetByClinicID(ctx context.Context, clinicID string) ([]domain.PnlReport, error)
	UpdateStatus(ctx context.Context, id string, status string, finalizedAt *time.Time) error
	DeleteLinesByReportID(ctx context.Context, reportID string) error
	UpdateReport(ctx context.Context, report *domain.PnlReport) error
	Delete(ctx context.Context, id string) error
	AggregateLedger(ctx context.Context, clinicID string, periodStart, periodEnd time.Time) ([]domain.LedgerAggRow, error)
}
