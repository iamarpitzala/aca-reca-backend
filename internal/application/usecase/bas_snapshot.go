package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/application/port"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
)

var ErrBASSnapshotNotFound = errors.New("BAS snapshot not found")
var ErrBASAlreadyFinalised = errors.New("BAS is already finalised and cannot be modified")
var ErrBASPeriodLocked = errors.New("BAS period is locked and cannot be modified")

type BASSnapshotService struct {
	repo port.BASSnapshotRepository
	clinicRepo port.ClinicRepository
}

func NewBASSnapshotService(repo port.BASSnapshotRepository, clinicRepo port.ClinicRepository) *BASSnapshotService {
	return &BASSnapshotService{
		repo:      repo,
		clinicRepo: clinicRepo,
	}
}

// Create creates a new BAS snapshot (draft)
func (s *BASSnapshotService) Create(ctx context.Context, clinicID uuid.UUID, req *domain.BASSnapshotRequest) (*domain.BASSnapshot, error) {
	// Verify clinic exists
	_, err := s.clinicRepo.GetByID(ctx, clinicID)
	if err != nil {
		return nil, err
	}
	
	// Check if a finalised BAS already exists for this period
	existing, err := s.repo.GetByClinicIDAndPeriod(ctx, clinicID, req.PeriodStart, req.PeriodEnd)
	if err == nil && existing != nil && (existing.Status == domain.BASStatusFinalised || existing.Status == domain.BASStatusLocked) {
		return nil, ErrBASAlreadyFinalised
	}
	
	snapshot := req.ToBASSnapshot(clinicID)
	if err := s.repo.Create(ctx, snapshot); err != nil {
		return nil, err
	}
	return snapshot, nil
}

// GetByID retrieves a BAS snapshot by ID
func (s *BASSnapshotService) GetByID(ctx context.Context, id uuid.UUID) (*domain.BASSnapshot, error) {
	return s.repo.GetByID(ctx, id)
}

// GetByClinicIDAndPeriod retrieves a BAS snapshot for a specific clinic and period
func (s *BASSnapshotService) GetByClinicIDAndPeriod(ctx context.Context, clinicID uuid.UUID, periodStart, periodEnd time.Time) (*domain.BASSnapshot, error) {
	return s.repo.GetByClinicIDAndPeriod(ctx, clinicID, periodStart, periodEnd)
}

// GetByClinicID retrieves all BAS snapshots for a clinic
func (s *BASSnapshotService) GetByClinicID(ctx context.Context, clinicID uuid.UUID) ([]domain.BASSnapshot, error) {
	return s.repo.GetByClinicID(ctx, clinicID)
}

// Finalise finalises a BAS snapshot (locks the period)
func (s *BASSnapshotService) Finalise(ctx context.Context, id uuid.UUID, finalisedBy uuid.UUID) (*domain.BASSnapshot, error) {
	snapshot, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	if snapshot.Status == domain.BASStatusFinalised || snapshot.Status == domain.BASStatusLocked {
		return nil, ErrBASAlreadyFinalised
	}
	
	now := time.Now()
	snapshot.Status = domain.BASStatusFinalised
	snapshot.FinalisedAt = &now
	snapshot.FinalisedBy = &finalisedBy
	
	if err := s.repo.Update(ctx, snapshot); err != nil {
		return nil, err
	}
	return snapshot, nil
}

// Lock locks a BAS snapshot (prevents any further changes)
func (s *BASSnapshotService) Lock(ctx context.Context, id uuid.UUID) (*domain.BASSnapshot, error) {
	snapshot, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	if snapshot.Status == domain.BASStatusLocked {
		return snapshot, nil // Already locked
	}
	
	snapshot.Status = domain.BASStatusLocked
	
	if err := s.repo.Update(ctx, snapshot); err != nil {
		return nil, err
	}
	return snapshot, nil
}

// Update updates a BAS snapshot (only if not finalised/locked)
func (s *BASSnapshotService) Update(ctx context.Context, id uuid.UUID, req *domain.BASSnapshotRequest) (*domain.BASSnapshot, error) {
	snapshot, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	
	if snapshot.Status == domain.BASStatusFinalised || snapshot.Status == domain.BASStatusLocked {
		return nil, ErrBASAlreadyFinalised
	}
	
	// Update fields
	snapshot.G1TotalSales = req.G1TotalSales
	snapshot.G2ExportSales = req.G2ExportSales
	snapshot.G3GSTFreeSales = req.G3GSTFreeSales
	snapshot.G10CapitalPurchases = req.G10CapitalPurchases
	snapshot.G11NonCapitalPurchases = req.G11NonCapitalPurchases
	snapshot.Label1AGSTOnSales = req.Label1AGSTOnSales
	snapshot.Label1BGSTOnPurchases = req.Label1BGSTOnPurchases
	snapshot.NetGSTPayable = req.NetGSTPayable
	snapshot.SnapshotData = req.SnapshotData
	
	if err := s.repo.Update(ctx, snapshot); err != nil {
		return nil, err
	}
	return snapshot, nil
}

// GetConsolidatedGSTSummary retrieves consolidated GST data across multiple clinics (read-only)
func (s *BASSnapshotService) GetConsolidatedGSTSummary(ctx context.Context, clinicIDs []uuid.UUID, periodStart, periodEnd time.Time) ([]domain.BASSnapshot, error) {
	// Verify all clinics exist
	for _, clinicID := range clinicIDs {
		_, err := s.clinicRepo.GetByID(ctx, clinicID)
		if err != nil {
			return nil, errors.New("one or more clinics not found")
		}
	}
	
	// Get only finalised BAS snapshots
	return s.repo.GetFinalisedByClinicIDs(ctx, clinicIDs, periodStart, periodEnd)
}
