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

const cogsAccountCode = "310" // Cost of Goods Sold

var (
	ErrPnlReportNotFound  = errors.New("P&L report not found")
	ErrPnlReportNotDraft  = errors.New("only DRAFT reports can be modified")
	ErrPnlReportFinalized = errors.New("report is already finalized")
	ErrInvalidPeriod      = errors.New("periodStart must be before periodEnd")
	ErrInvalidPeriodStart = errors.New("invalid periodStart format, use YYYY-MM-DD")
	ErrInvalidPeriodEnd   = errors.New("invalid periodEnd format, use YYYY-MM-DD")
)

type PnlReportService struct {
	repo       port.PnlReportRepository
	clinicRepo port.ClinicRepository
}

func NewPnlReportService(
	repo port.PnlReportRepository,
	clinicRepo port.ClinicRepository,
) *PnlReportService {
	return &PnlReportService{
		repo:       repo,
		clinicRepo: clinicRepo,
	}
}

// Generate creates a new P&L report by aggregating the transaction ledger.
func (s *PnlReportService) Generate(
	ctx context.Context,
	req *domain.GeneratePnlRequest,
	userID string,
) (*domain.PnlReportWithLines, error) {
	// Verify clinic exists
	if _, err := s.clinicRepo.GetByID(ctx, req.ClinicID); err != nil {
		return nil, err
	}

	// Parse period dates
	periodStart, err := time.Parse(util.DateFormatDate, req.PeriodStart)
	if err != nil {
		return nil, ErrInvalidPeriodStart
	}
	periodEnd, err := time.Parse(util.DateFormatDate, req.PeriodEnd)
	if err != nil {
		return nil, ErrInvalidPeriodEnd
	}
	if !periodStart.Before(periodEnd) {
		return nil, ErrInvalidPeriod
	}

	// Aggregate from ledger
	aggRows, err := s.repo.AggregateLedger(ctx, req.ClinicID, periodStart, periodEnd)
	if err != nil {
		return nil, err
	}

	// Build lines and compute totals
	now := time.Now()
	reportID := uuid.NewString()
	lines, totals := buildPnlLines(reportID, aggRows, now)

	report := &domain.PnlReport{
		ID:                reportID,
		ClinicID:          req.ClinicID,
		QuarterID:         req.QuarterID,
		ReportName:        req.ReportName,
		PeriodStart:       periodStart,
		PeriodEnd:         periodEnd,
		TotalRevenue:      totals.totalRevenue,
		TotalCOGS:         totals.totalCOGS,
		GrossProfit:       totals.grossProfit,
		TotalExpenses:     totals.totalExpenses,
		NetProfit:         totals.netProfit,
		TotalGSTCollected: totals.totalGSTCollected,
		TotalGSTPaid:      totals.totalGSTPaid,
		NetGST:            totals.netGST,
		Status:            util.PnlReportStatusDraft,
		GeneratedBy:       userID,
		Notes:             req.Notes,
		CreatedAt:         now,
		UpdatedAt:         now,
	}

	if err := s.repo.Create(ctx, report); err != nil {
		return nil, err
	}
	if err := s.repo.CreateLines(ctx, lines); err != nil {
		return nil, err
	}

	return &domain.PnlReportWithLines{
		Report: *report,
		Lines:  lines,
	}, nil
}

// GetByID retrieves a P&L report with its line items.
func (s *PnlReportService) GetByID(ctx context.Context, id string) (*domain.PnlReportWithLines, error) {
	report, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if report == nil {
		return nil, ErrPnlReportNotFound
	}

	lines, err := s.repo.GetLinesByReportID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &domain.PnlReportWithLines{
		Report: *report,
		Lines:  lines,
	}, nil
}

// GetByClinicID lists all P&L reports for a clinic.
func (s *PnlReportService) GetByClinicID(ctx context.Context, clinicID string) ([]domain.PnlReport, error) {
	return s.repo.GetByClinicID(ctx, clinicID)
}

// Finalize changes a DRAFT report to FINAL, preventing further modifications.
func (s *PnlReportService) Finalize(ctx context.Context, id string) (*domain.PnlReport, error) {
	report, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if report == nil {
		return nil, ErrPnlReportNotFound
	}
	if report.Status == util.PnlReportStatusFinal {
		return nil, ErrPnlReportFinalized
	}
	if report.Status != util.PnlReportStatusDraft {
		return nil, ErrPnlReportNotDraft
	}

	now := time.Now()
	if err := s.repo.UpdateStatus(ctx, id, util.PnlReportStatusFinal, &now); err != nil {
		return nil, err
	}

	report.Status = util.PnlReportStatusFinal
	report.FinalizedAt = &now
	return report, nil
}

// Regenerate re-aggregates the ledger and replaces lines on a DRAFT report.
func (s *PnlReportService) Regenerate(ctx context.Context, id string) (*domain.PnlReportWithLines, error) {
	report, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if report == nil {
		return nil, ErrPnlReportNotFound
	}
	if report.Status != util.PnlReportStatusDraft {
		return nil, ErrPnlReportNotDraft
	}

	// Re-aggregate
	aggRows, err := s.repo.AggregateLedger(ctx, report.ClinicID, report.PeriodStart, report.PeriodEnd)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	lines, totals := buildPnlLines(report.ID, aggRows, now)

	// Delete old lines and insert new
	if err := s.repo.DeleteLinesByReportID(ctx, report.ID); err != nil {
		return nil, err
	}
	if err := s.repo.CreateLines(ctx, lines); err != nil {
		return nil, err
	}

	// Update totals on header
	report.TotalRevenue = totals.totalRevenue
	report.TotalCOGS = totals.totalCOGS
	report.GrossProfit = totals.grossProfit
	report.TotalExpenses = totals.totalExpenses
	report.NetProfit = totals.netProfit
	report.TotalGSTCollected = totals.totalGSTCollected
	report.TotalGSTPaid = totals.totalGSTPaid
	report.NetGST = totals.netGST
	report.UpdatedAt = now

	if err := s.repo.UpdateReport(ctx, report); err != nil {
		return nil, err
	}

	return &domain.PnlReportWithLines{
		Report: *report,
		Lines:  lines,
	}, nil
}

// Delete soft-deletes a P&L report.
func (s *PnlReportService) Delete(ctx context.Context, id string) error {
	report, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if report == nil {
		return ErrPnlReportNotFound
	}
	return s.repo.Delete(ctx, id)
}

// ---- helpers ----

type pnlTotals struct {
	totalRevenue      float64
	totalCOGS         float64
	grossProfit       float64
	totalExpenses     float64
	netProfit         float64
	totalGSTCollected float64
	totalGSTPaid      float64
	netGST            float64
}

// buildPnlLines converts aggregated ledger rows into report lines and computes totals.
func buildPnlLines(reportID string, aggRows []domain.LedgerAggRow, now time.Time) ([]domain.PnlReportLine, pnlTotals) {
	lines := make([]domain.PnlReportLine, 0, len(aggRows))
	var totals pnlTotals

	orderRevenue := 100
	orderCOGS := 200
	orderExpense := 300

	for _, row := range aggRows {
		category := classifyAccount(row.AccountType, row.AccountCode)
		var displayOrder int

		switch category {
		case "REVENUE":
			displayOrder = orderRevenue
			orderRevenue++
			// For revenue accounts credit_total represents revenue earned
			totals.totalRevenue += row.CreditTotal
			totals.totalGSTCollected += row.GSTAmount
		case "COGS":
			displayOrder = orderCOGS
			orderCOGS++
			totals.totalCOGS += row.DebitTotal
			totals.totalGSTPaid += row.GSTAmount
		case "EXPENSE":
			displayOrder = orderExpense
			orderExpense++
			totals.totalExpenses += row.DebitTotal
			totals.totalGSTPaid += row.GSTAmount
		}

		lines = append(lines, domain.PnlReportLine{
			ID:               uuid.NewString(),
			PnlReportID:      reportID,
			COAID:            row.COAID,
			AccountCode:      row.AccountCode,
			AccountName:      row.AccountName,
			LineCategory:     category,
			DebitTotal:       row.DebitTotal,
			CreditTotal:      row.CreditTotal,
			NetAmount:        row.NetAmount,
			GSTAmount:        row.GSTAmount,
			TransactionCount: row.TransactionCount,
			DisplayOrder:     displayOrder,
			CreatedAt:        now,
		})
	}

	totals.grossProfit = totals.totalRevenue - totals.totalCOGS
	totals.netProfit = totals.grossProfit - totals.totalExpenses
	totals.netGST = totals.totalGSTCollected - totals.totalGSTPaid

	return lines, totals
}

// classifyAccount maps account type + code to a P&L line category.
func classifyAccount(accountType, accountCode string) string {
	if accountType == "Revenue" {
		return "REVENUE"
	}
	if accountType == "Expense" && accountCode == cogsAccountCode {
		return "COGS"
	}
	return "EXPENSE"
}
