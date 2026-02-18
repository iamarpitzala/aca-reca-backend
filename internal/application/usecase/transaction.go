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
	ErrTransactionNotFound     = errors.New("transaction not found")
	ErrTransactionAlreadyVoided = errors.New("transaction is already voided")
	ErrTransactionNotDraft     = errors.New("only DRAFT transactions can be updated")
	ErrLedgerImbalanced        = errors.New("ledger debits and credits must balance")
	ErrInvalidTransactionDate  = errors.New("invalid transaction date format, use YYYY-MM-DD")
	ErrLedgerRequired          = errors.New("at least one ledger line is required")
)

type TransactionService struct {
	repo       port.TransactionRepository
	clinicRepo port.ClinicRepository
}

func NewTransactionService(
	repo port.TransactionRepository,
	clinicRepo port.ClinicRepository,
) *TransactionService {
	return &TransactionService{
		repo:       repo,
		clinicRepo: clinicRepo,
	}
}

// Create creates a new transaction with ledger lines
func (s *TransactionService) Create(
	ctx context.Context,
	req *domain.CreateTransactionRequest,
	userID uuid.UUID,
) (*domain.TransactionWithLedger, error) {
	// Verify clinic exists
	_, err := s.clinicRepo.GetByID(ctx, req.ClinicID)
	if err != nil {
		return nil, err
	}

	// Parse transaction date
	txnDate, err := time.Parse(util.DateFormatDate, req.TransactionDate)
	if err != nil {
		return nil, ErrInvalidTransactionDate
	}

	// Validate ledger lines
	if len(req.Ledger) == 0 {
		return nil, ErrLedgerRequired
	}
	if err := validateLedgerBalance(req.Ledger); err != nil {
		return nil, err
	}

	// Determine status
	status := util.TransactionStatusDraft
	if req.Status == util.TransactionStatusPosted {
		status = util.TransactionStatusPosted
	}

	now := time.Now()
	var postedAt *time.Time
	if status == util.TransactionStatusPosted {
		postedAt = &now
	}

	txn := &domain.Transaction{
		ID:              uuid.New(),
		ClinicID:        req.ClinicID,
		SourceEntryID:   req.SourceEntryID,
		ReferenceNumber: req.ReferenceNumber,
		Description:     req.Description,
		TransactionDate: txnDate,
		Status:          status,
		CreatedBy:       userID,
		PostedAt:        postedAt,
		CreatedAt:       now,
		UpdatedAt:       now,
	}

	if err := s.repo.Create(ctx, txn); err != nil {
		return nil, err
	}

	// Build and save ledger lines
	lines := buildLedgerLines(txn.ID, txnDate, now, req.Ledger)
	if err := s.repo.CreateLedgerLines(ctx, lines); err != nil {
		return nil, err
	}

	return &domain.TransactionWithLedger{
		Transaction: *txn,
		Ledger:      lines,
	}, nil
}

// GetByID retrieves a transaction with its ledger lines
func (s *TransactionService) GetByID(ctx context.Context, id uuid.UUID) (*domain.TransactionWithLedger, error) {
	txn, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if txn == nil {
		return nil, ErrTransactionNotFound
	}

	ledger, err := s.repo.GetLedgerByTransactionID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &domain.TransactionWithLedger{
		Transaction: *txn,
		Ledger:      ledger,
	}, nil
}

// GetByClinicID lists all transactions for a clinic
func (s *TransactionService) GetByClinicID(ctx context.Context, clinicID uuid.UUID) ([]domain.TransactionWithLedger, error) {
	txns, err := s.repo.GetByClinicID(ctx, clinicID)
	if err != nil {
		return nil, err
	}

	result := make([]domain.TransactionWithLedger, 0, len(txns))
	for _, txn := range txns {
		ledger, err := s.repo.GetLedgerByTransactionID(ctx, txn.ID)
		if err != nil {
			return nil, err
		}
		result = append(result, domain.TransactionWithLedger{
			Transaction: txn,
			Ledger:      ledger,
		})
	}
	return result, nil
}

// Update updates a DRAFT transaction (header + ledger lines)
func (s *TransactionService) Update(
	ctx context.Context,
	id uuid.UUID,
	req *domain.UpdateTransactionRequest,
) (*domain.TransactionWithLedger, error) {
	txn, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if txn == nil {
		return nil, ErrTransactionNotFound
	}
	if txn.Status != util.TransactionStatusDraft {
		return nil, ErrTransactionNotDraft
	}

	now := time.Now()

	if req.ReferenceNumber != nil {
		txn.ReferenceNumber = req.ReferenceNumber
	}
	if req.Description != nil {
		txn.Description = req.Description
	}
	if req.TransactionDate != nil {
		txnDate, err := time.Parse(util.DateFormatDate, *req.TransactionDate)
		if err != nil {
			return nil, ErrInvalidTransactionDate
		}
		txn.TransactionDate = txnDate
	}
	txn.UpdatedAt = now

	if err := s.repo.Update(ctx, txn); err != nil {
		return nil, err
	}

	// Replace ledger lines if provided
	var ledger []domain.TransactionLedger
	if len(req.Ledger) > 0 {
		if err := validateLedgerBalance(req.Ledger); err != nil {
			return nil, err
		}
		if err := s.repo.DeleteLedgerByTransactionID(ctx, id); err != nil {
			return nil, err
		}
		ledger = buildLedgerLines(id, txn.TransactionDate, now, req.Ledger)
		if err := s.repo.CreateLedgerLines(ctx, ledger); err != nil {
			return nil, err
		}
	} else {
		ledger, err = s.repo.GetLedgerByTransactionID(ctx, id)
		if err != nil {
			return nil, err
		}
	}

	return &domain.TransactionWithLedger{
		Transaction: *txn,
		Ledger:      ledger,
	}, nil
}

// PostTransaction changes a DRAFT transaction to POSTED
func (s *TransactionService) PostTransaction(ctx context.Context, id uuid.UUID) (*domain.TransactionWithLedger, error) {
	txn, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if txn == nil {
		return nil, ErrTransactionNotFound
	}
	if txn.Status != util.TransactionStatusDraft {
		return nil, ErrTransactionNotDraft
	}

	now := time.Now()
	txn.Status = util.TransactionStatusPosted
	txn.PostedAt = &now
	txn.UpdatedAt = now

	if err := s.repo.Update(ctx, txn); err != nil {
		return nil, err
	}

	ledger, err := s.repo.GetLedgerByTransactionID(ctx, id)
	if err != nil {
		return nil, err
	}

	return &domain.TransactionWithLedger{
		Transaction: *txn,
		Ledger:      ledger,
	}, nil
}

// VoidTransaction marks a transaction as VOIDED
func (s *TransactionService) VoidTransaction(ctx context.Context, id uuid.UUID, reason string) (*domain.Transaction, error) {
	txn, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if txn == nil {
		return nil, ErrTransactionNotFound
	}
	if txn.Status == util.TransactionStatusVoided {
		return nil, ErrTransactionAlreadyVoided
	}

	if err := s.repo.VoidTransaction(ctx, id, reason); err != nil {
		return nil, err
	}

	// Re-fetch to return updated state
	txn, err = s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	return txn, nil
}

// Delete soft-deletes a transaction
func (s *TransactionService) Delete(ctx context.Context, id uuid.UUID) error {
	txn, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if txn == nil {
		return ErrTransactionNotFound
	}
	return s.repo.Delete(ctx, id)
}

// ---- helpers ----

func validateLedgerBalance(lines []domain.LedgerLineRequest) error {
	var totalDebit, totalCredit float64
	for _, l := range lines {
		switch l.EntryType {
		case "DEBIT":
			totalDebit += l.Amount
		case "CREDIT":
			totalCredit += l.Amount
		}
	}
	diff := totalDebit - totalCredit
	if diff < -0.01 || diff > 0.01 {
		return ErrLedgerImbalanced
	}
	return nil
}

func buildLedgerLines(
	transactionID uuid.UUID,
	txnDate time.Time,
	now time.Time,
	lines []domain.LedgerLineRequest,
) []domain.TransactionLedger {
	result := make([]domain.TransactionLedger, 0, len(lines))
	for _, l := range lines {
		netAmount := l.NetAmount
		if netAmount == 0 {
			if l.EntryType == "DEBIT" {
				netAmount = l.Amount
			} else {
				netAmount = -l.Amount
			}
		}
		result = append(result, domain.TransactionLedger{
			ID:              uuid.New(),
			TransactionID:   transactionID,
			COAID:           l.COAID,
			EntryType:       l.EntryType,
			Amount:          l.Amount,
			GSTAmount:       l.GSTAmount,
			NetAmount:       netAmount,
			TransactionDate: txnDate,
			Description:     l.Description,
			CreatedAt:       now,
		})
	}
	return result
}
