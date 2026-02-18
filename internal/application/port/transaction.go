package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
)

type TransactionRepository interface {
	Create(ctx context.Context, txn *domain.Transaction) error
	CreateLedgerLines(ctx context.Context, lines []domain.TransactionLedger) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Transaction, error)
	GetLedgerByTransactionID(ctx context.Context, transactionID uuid.UUID) ([]domain.TransactionLedger, error)
	GetByClinicID(ctx context.Context, clinicID uuid.UUID) ([]domain.Transaction, error)
	Update(ctx context.Context, txn *domain.Transaction) error
	DeleteLedgerByTransactionID(ctx context.Context, transactionID uuid.UUID) error
	VoidTransaction(ctx context.Context, id uuid.UUID, reason string) error
	Delete(ctx context.Context, id uuid.UUID) error
}
