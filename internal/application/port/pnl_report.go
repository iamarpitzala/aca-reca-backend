package port

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
)

type PnlReportRepository interface {
	Create(ctx context.Context, report *domain.PnlReport) error
	CreateLines(ctx context.Context, lines []domain.PnlReportLine) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.PnlReport, error)
	GetLinesByReportID(ctx context.Context, reportID uuid.UUID) ([]domain.PnlReportLine, error)
	GetByClinicID(ctx context.Context, clinicID uuid.UUID) ([]domain.PnlReport, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string, finalizedAt *time.Time) error
	DeleteLinesByReportID(ctx context.Context, reportID uuid.UUID) error
	UpdateReport(ctx context.Context, report *domain.PnlReport) error
	Delete(ctx context.Context, id uuid.UUID) error
	AggregateLedger(ctx context.Context, clinicID uuid.UUID, periodStart, periodEnd time.Time) ([]domain.LedgerAggRow, error)
}
