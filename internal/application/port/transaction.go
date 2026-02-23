package port

import (
	"context"

	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
)

type TransactionRepository interface {
	Create(ctx context.Context, txn *domain.Transaction) error
	CreateLedgerLines(ctx context.Context, lines []domain.TransactionLedger) error
	GetByID(ctx context.Context, id string) (*domain.Transaction, error)
	GetLedgerByTransactionID(ctx context.Context, transactionID string) ([]domain.TransactionLedger, error)
	GetByClinicID(ctx context.Context, clinicID string) ([]domain.Transaction, error)
	Update(ctx context.Context, txn *domain.Transaction) error
	DeleteLedgerByTransactionID(ctx context.Context, transactionID string) error
	VoidTransaction(ctx context.Context, id string, reason string) error
	Delete(ctx context.Context, id string) error
}
