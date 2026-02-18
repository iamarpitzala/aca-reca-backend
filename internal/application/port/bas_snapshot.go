package port

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
)

type BASSnapshotRepository interface {
	Create(ctx context.Context, snapshot *domain.BASSnapshot) error
	CreateLines(ctx context.Context, lines []domain.BASSnapshotLine) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.BASSnapshot, error)
	GetLinesBySnapshotID(ctx context.Context, snapshotID uuid.UUID) ([]domain.BASSnapshotLine, error)
	GetByClinicID(ctx context.Context, clinicID uuid.UUID) ([]domain.BASSnapshot, error)
	Update(ctx context.Context, snapshot *domain.BASSnapshot) error
	UpdateStatus(ctx context.Context, id uuid.UUID, status string, finalisedAt *time.Time, finalisedBy *uuid.UUID) error
	UpdateLock(ctx context.Context, id uuid.UUID, lockedAt *time.Time, lockedBy *uuid.UUID) error
	DeleteLinesBySnapshotID(ctx context.Context, snapshotID uuid.UUID) error
	Delete(ctx context.Context, id uuid.UUID) error
	AggregateLedger(ctx context.Context, clinicID uuid.UUID, periodStart, periodEnd time.Time) ([]domain.BASLedgerAggRow, error)
}
