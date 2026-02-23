package port

import (
	"context"

	"github.com/iamarpitzala/aca-reca-backend/internal/domain/coa"
)

// ArrangementRepository defines persistence for tbl_arrangement.
type ArrangementRepository interface {
	Create(ctx context.Context, a *coa.Arrangement) error
	GetByID(ctx context.Context, id, userID string) (*coa.Arrangement, error)
	ListByUserID(ctx context.Context, userID string) ([]coa.Arrangement, error)
	Update(ctx context.Context, a *coa.Arrangement) error
	Delete(ctx context.Context, id string) error
}
