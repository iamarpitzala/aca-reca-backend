package usecase

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/application/port"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
	"github.com/iamarpitzala/aca-reca-backend/util"
)

var (
	ErrBASSnapshotNotFound   = errors.New("BAS snapshot not found")
	ErrBASSnapshotNotDraft   = errors.New("only DRAFT snapshots can be modified")
	ErrBASAlreadyFinalised   = errors.New("snapshot is already finalised")
	ErrBASAlreadyLocked      = errors.New("snapshot is already locked")
	ErrBASNotFinalised       = errors.New("snapshot must be FINALISED before locking")
	ErrBASInvalidPeriod      = errors.New("periodStart must be before periodEnd")
	ErrBASInvalidPeriodStart = errors.New("invalid periodStart format, use YYYY-MM-DD")
	ErrBASInvalidPeriodEnd   = errors.New("invalid periodEnd format, use YYYY-MM-DD")
	ErrBASInvalidPeriodType  = errors.New("periodType must be QUARTERLY or ANNUALLY")
)

type BASSnapshotService struct {
	repo       port.BASSnapshotRepository
	clinicRepo port.ClinicRepository
}

func NewBASSnapshotService(
	repo port.BASSnapshotRepository,
	clinicRepo port.ClinicRepository,
) *BASSnapshotService {
	return &BASSnapshotService{
		repo:       repo,
		clinicRepo: clinicRepo,
	}
}

// Generate creates a new BAS snapshot by aggregating the transaction ledger.
func (s *BASSnapshotService) Generate(
	ctx context.Context,
	clinicID uuid.UUID,
	req *domain.GenerateBASRequest,
	userID uuid.UUID,
) (*domain.BASSnapshotWithLines, error) {
	if _, err := s.clinicRepo.GetByID(ctx, clinicID); err != nil {
		return nil, err
	}

	periodStart, err := time.Parse(util.DateFormatDate, req.PeriodStart)
	if err != nil {
		return nil, ErrBASInvalidPeriodStart
	}
	periodEnd, err := time.Parse(util.DateFormatDate, req.PeriodEnd)
	if err != nil {
		return nil, ErrBASInvalidPeriodEnd
	}
	if !periodStart.Before(periodEnd) {
		return nil, ErrBASInvalidPeriod
	}
	if req.PeriodType != util.BASPeriodTypeQuarterly && req.PeriodType != util.BASPeriodTypeAnnually {
		return nil, ErrBASInvalidPeriodType
	}

	aggRows, err := s.repo.AggregateLedger(ctx, clinicID, periodStart, periodEnd)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	snapshotID := uuid.New()
	lines, totals := buildBASLines(snapshotID, aggRows, now)

	snapshot := &domain.BASSnapshot{
		ID:          snapshotID,
		ClinicID:    clinicID,
		PeriodStart: periodStart,
		PeriodEnd:   periodEnd,
		PeriodType:  req.PeriodType,

		G1TotalSales:            totals.g1TotalSales,
		G2ExportSales:           totals.g2ExportSales,
		G3GSTFreeSales:          totals.g3GSTFreeSales,
		G5G2G3G4:                totals.g2ExportSales + totals.g3GSTFreeSales,
		G6TotalTaxableSales:     totals.g1TotalSales - (totals.g2ExportSales + totals.g3GSTFreeSales),
		G8TotalTaxableSupplies:  totals.g1TotalSales - (totals.g2ExportSales + totals.g3GSTFreeSales),
		G9GSTOnSales:            totals.gstOnSales,

		G10CapitalPurchases:          totals.g10CapitalPurchases,
		G11NonCapitalPurchases:       totals.g11NonCapitalPurchases,
		G12G10G11:                    totals.g10CapitalPurchases + totals.g11NonCapitalPurchases,
		G14PurchasesWithoutGST:       totals.g14PurchasesWithoutGST,
		G16G13G14G15:                 totals.g14PurchasesWithoutGST,
		G17TotalCreditablePurchases:  (totals.g10CapitalPurchases + totals.g11NonCapitalPurchases) - totals.g14PurchasesWithoutGST,
		G19TotalCreditableAcquisitions: (totals.g10CapitalPurchases + totals.g11NonCapitalPurchases) - totals.g14PurchasesWithoutGST,
		G20GSTOnPurchases:            totals.gstOnPurchases,

		Label1AGSTOnSales:     totals.gstOnSales,
		Label1BGSTOnPurchases: totals.gstOnPurchases,

		NetGSTPayable:    totals.gstOnSales - totals.gstOnPurchases,
		TotalAmountOwing: totals.gstOnSales - totals.gstOnPurchases,

		Status:       util.BASStatusDraft,
		GeneratedBy:  userID,
		SnapshotData: "{}",
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.repo.Create(ctx, snapshot); err != nil {
		return nil, err
	}
	if err := s.repo.CreateLines(ctx, lines); err != nil {
		return nil, err
	}

	return &domain.BASSnapshotWithLines{
		Snapshot: *snapshot,
		Lines:    lines,
	}, nil
}

// GetByID retrieves a BAS snapshot with its line items.
func (s *BASSnapshotService) GetByID(ctx context.Context, id uuid.UUID) (*domain.BASSnapshotWithLines, error) {
	snapshot, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if snapshot == nil {
		return nil, ErrBASSnapshotNotFound
	}

	lines, err := s.repo.GetLinesBySnapshotID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &domain.BASSnapshotWithLines{
		Snapshot: *snapshot,
		Lines:    lines,
	}, nil
}

// GetByClinicID lists all BAS snapshots for a clinic.
func (s *BASSnapshotService) GetByClinicID(ctx context.Context, clinicID uuid.UUID) ([]domain.BASSnapshot, error) {
	return s.repo.GetByClinicID(ctx, clinicID)
}

// Update modifies a DRAFT BAS snapshot's fields.
func (s *BASSnapshotService) Update(ctx context.Context, id uuid.UUID, req *domain.UpdateBASSnapshotRequest) (*domain.BASSnapshot, error) {
	snapshot, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if snapshot == nil {
		return nil, ErrBASSnapshotNotFound
	}
	if snapshot.Status != util.BASStatusDraft {
		return nil, ErrBASSnapshotNotDraft
	}

	if req.PeriodStart != nil {
		ps, err := time.Parse(util.DateFormatDate, *req.PeriodStart)
		if err != nil {
			return nil, ErrBASInvalidPeriodStart
		}
		snapshot.PeriodStart = ps
	}
	if req.PeriodEnd != nil {
		pe, err := time.Parse(util.DateFormatDate, *req.PeriodEnd)
		if err != nil {
			return nil, ErrBASInvalidPeriodEnd
		}
		snapshot.PeriodEnd = pe
	}
	if !snapshot.PeriodStart.Before(snapshot.PeriodEnd) {
		return nil, ErrBASInvalidPeriod
	}
	if req.PeriodType != nil {
		if *req.PeriodType != util.BASPeriodTypeQuarterly && *req.PeriodType != util.BASPeriodTypeAnnually {
			return nil, ErrBASInvalidPeriodType
		}
		snapshot.PeriodType = *req.PeriodType
	}
	if req.G1TotalSales != nil {
		snapshot.G1TotalSales = *req.G1TotalSales
	}
	if req.G2ExportSales != nil {
		snapshot.G2ExportSales = *req.G2ExportSales
	}
	if req.G3GSTFreeSales != nil {
		snapshot.G3GSTFreeSales = *req.G3GSTFreeSales
	}
	if req.G10CapitalPurchases != nil {
		snapshot.G10CapitalPurchases = *req.G10CapitalPurchases
	}
	if req.G11NonCapitalPurchases != nil {
		snapshot.G11NonCapitalPurchases = *req.G11NonCapitalPurchases
	}
	if req.Label1AGSTOnSales != nil {
		snapshot.Label1AGSTOnSales = *req.Label1AGSTOnSales
	}
	if req.Label1BGSTOnPurchases != nil {
		snapshot.Label1BGSTOnPurchases = *req.Label1BGSTOnPurchases
	}
	if req.NetGSTPayable != nil {
		snapshot.NetGSTPayable = *req.NetGSTPayable
	}
	if req.Notes != nil {
		snapshot.Notes = req.Notes
	}
	if req.SnapshotData != nil {
		snapshot.SnapshotData = req.SnapshotData
	}

	// Recompute derived fields
	snapshot.G5G2G3G4 = snapshot.G2ExportSales + snapshot.G3GSTFreeSales + snapshot.G4InputTaxedSales
	snapshot.G6TotalTaxableSales = snapshot.G1TotalSales - snapshot.G5G2G3G4
	snapshot.G8TotalTaxableSupplies = snapshot.G6TotalTaxableSales + snapshot.G7Adjustments
	snapshot.G12G10G11 = snapshot.G10CapitalPurchases + snapshot.G11NonCapitalPurchases
	snapshot.G16G13G14G15 = snapshot.G13PurchasesForInputTaxed + snapshot.G14PurchasesWithoutGST + snapshot.G15PrivateUse
	snapshot.G17TotalCreditablePurchases = snapshot.G12G10G11 - snapshot.G16G13G14G15
	snapshot.G19TotalCreditableAcquisitions = snapshot.G17TotalCreditablePurchases + snapshot.G18Adjustments
	snapshot.TotalAmountOwing = snapshot.NetGSTPayable + snapshot.W4TotalWithheld + snapshot.T4InstalmentAmount

	snapshot.UpdatedAt = time.Now()

	if err := s.repo.Update(ctx, snapshot); err != nil {
		return nil, err
	}
	return snapshot, nil
}

// Finalise changes a DRAFT snapshot to FINALISED.
func (s *BASSnapshotService) Finalise(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*domain.BASSnapshot, error) {
	snapshot, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if snapshot == nil {
		return nil, ErrBASSnapshotNotFound
	}
	if snapshot.Status == util.BASStatusFinalised {
		return nil, ErrBASAlreadyFinalised
	}
	if snapshot.Status != util.BASStatusDraft {
		return nil, ErrBASSnapshotNotDraft
	}

	now := time.Now()
	if err := s.repo.UpdateStatus(ctx, id, util.BASStatusFinalised, &now, &userID); err != nil {
		return nil, err
	}

	snapshot.Status = util.BASStatusFinalised
	snapshot.FinalisedAt = &now
	snapshot.FinalisedBy = &userID
	return snapshot, nil
}

// Lock changes a FINALISED snapshot to LOCKED.
func (s *BASSnapshotService) Lock(ctx context.Context, id uuid.UUID, userID uuid.UUID) (*domain.BASSnapshot, error) {
	snapshot, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if snapshot == nil {
		return nil, ErrBASSnapshotNotFound
	}
	if snapshot.Status == util.BASStatusLocked {
		return nil, ErrBASAlreadyLocked
	}
	if snapshot.Status != util.BASStatusFinalised {
		return nil, ErrBASNotFinalised
	}

	now := time.Now()
	if err := s.repo.UpdateLock(ctx, id, &now, &userID); err != nil {
		return nil, err
	}

	snapshot.Status = util.BASStatusLocked
	snapshot.LockedAt = &now
	snapshot.LockedBy = &userID
	return snapshot, nil
}

// Delete soft-deletes a BAS snapshot.
func (s *BASSnapshotService) Delete(ctx context.Context, id uuid.UUID) error {
	snapshot, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if snapshot == nil {
		return ErrBASSnapshotNotFound
	}
	return s.repo.Delete(ctx, id)
}

// ---- helpers ----

type basTotals struct {
	g1TotalSales           float64
	g2ExportSales          float64
	g3GSTFreeSales         float64
	g10CapitalPurchases    float64
	g11NonCapitalPurchases float64
	g14PurchasesWithoutGST float64
	gstOnSales             float64
	gstOnPurchases         float64
}

// buildBASLines converts aggregated ledger rows into BAS snapshot lines and computes totals.
func buildBASLines(snapshotID uuid.UUID, aggRows []domain.BASLedgerAggRow, now time.Time) ([]domain.BASSnapshotLine, basTotals) {
	lines := make([]domain.BASSnapshotLine, 0, len(aggRows))
	var totals basTotals

	for _, row := range aggRows {
		label := classifyBASLabel(row.AccountType, row.AccountTaxName)

		var baseAmount, gstAmount float64

		switch row.AccountType {
		case "Revenue":
			baseAmount = row.CreditTotal
			gstAmount = row.GSTAmount
		default:
			baseAmount = row.DebitTotal
			gstAmount = row.GSTAmount
		}

		totalAmount := baseAmount + gstAmount

		switch label {
		case "G1":
			totals.g1TotalSales += totalAmount
			totals.gstOnSales += gstAmount
		case "G3":
			totals.g3GSTFreeSales += totalAmount
		case "G10":
			totals.g10CapitalPurchases += totalAmount
			totals.gstOnPurchases += gstAmount
		case "G11":
			totals.g11NonCapitalPurchases += totalAmount
			totals.gstOnPurchases += gstAmount
		case "G14":
			totals.g14PurchasesWithoutGST += baseAmount
		}

		lines = append(lines, domain.BASSnapshotLine{
			ID:               uuid.New(),
			BASSnapshotID:    snapshotID,
			COAID:            row.COAID,
			AccountCode:      row.AccountCode,
			AccountName:      row.AccountName,
			AccountTaxName:   row.AccountTaxName,
			BASLabel:         label,
			BaseAmount:       baseAmount,
			GSTAmount:        gstAmount,
			TotalAmount:      totalAmount,
			TransactionCount: row.TransactionCount,
			CreatedAt:        now,
		})
	}

	return lines, totals
}

// classifyBASLabel maps account type + tax name to BAS label.
func classifyBASLabel(accountType, accountTaxName string) string {
	switch accountTaxName {
	case "GST on Income":
		return "G1"
	case "GST Free Income":
		return "G3"
	case "GST on Expenses":
		if accountType == "Asset" {
			return "G10"
		}
		return "G11"
	case "GST Free Expenses":
		return "G14"
	default:
		return "BAS_EXCLUDED"
	}
}
