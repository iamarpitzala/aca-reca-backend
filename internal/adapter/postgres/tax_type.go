package postgres

import (
	"context"

	"github.com/iamarpitzala/aca-reca-backend/internal/application/port"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain/form"
	"github.com/jmoiron/sqlx"
)

type taxTypeRepo struct {
	db *sqlx.DB
}

func NewTaxTypeRepository(db *sqlx.DB) port.TaxTypeRepository {
	return &taxTypeRepo{db: db}
}

func (r *taxTypeRepo) ListAll(ctx context.Context) ([]form.TaxType, error) {
	query := `SELECT id, name, type, description, created_at, updated_at, deleted_at
		FROM tbl_tax_type WHERE deleted_at IS NULL ORDER BY id`
	var list []form.TaxType
	if err := r.db.SelectContext(ctx, &list, query); err != nil {
		return nil, err
	}
	if list == nil {
		list = []form.TaxType{}
	}
	return list, nil
}
