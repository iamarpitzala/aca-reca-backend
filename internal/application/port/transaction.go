package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
)

// TransactionRepository defines persistence for journal entries (posted transactions).
// In accounting terms: transactions are posted ledger entries derived from form entries.
type TransactionRepository interface {
	Create(ctx context.Context, t *domain.TransactionWithLedger) error
	ListByClinicID(ctx context.Context, clinicID uuid.UUID, f *domain.ListTransactionsFilters) ([]domain.TransactionWithLedger, int, error)
	ListByEntryID(ctx context.Context, entryID uuid.UUID) ([]domain.TransactionWithLedger, error)
	DeleteByEntryID(ctx context.Context, entryID uuid.UUID) error
}
