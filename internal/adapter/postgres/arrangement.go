package postgres

import (
	"context"
	"database/sql"
	"errors"
	"time"

	"github.com/iamarpitzala/aca-reca-backend/internal/application/port"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain/coa"
	"github.com/iamarpitzala/aca-reca-backend/util"
	"github.com/jmoiron/sqlx"
)

type arrangementRepo struct {
	db *sqlx.DB
}

func NewArrangementRepository(db *sqlx.DB) port.ArrangementRepository {
	return &arrangementRepo{db: db}
}

func (r *arrangementRepo) Create(ctx context.Context, a *coa.Arrangement) error {
	query := `INSERT INTO tbl_arrangement (id, user_id, method, name, percentage, amount, description, created_at, updated_at, deleted_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`
	now := time.Now()
	_, err := r.db.ExecContext(ctx, query,
		a.ID, a.UserID, a.Method, a.Name, a.Percentage, a.Amount, a.Description,
		now, now, nil,
	)
	return err
}

func (r *arrangementRepo) GetByID(ctx context.Context, id, userID string) (*coa.Arrangement, error) {
	query := `SELECT id, user_id, method, name, percentage, amount, description, created_at, updated_at, deleted_at
		FROM tbl_arrangement WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL`
	var a coa.Arrangement
	err := r.db.GetContext(ctx, &a, query, id, userID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, errors.New(util.ErrArrangementFailed)
	}
	return &a, nil
}

func (r *arrangementRepo) ListByUserID(ctx context.Context, userID string) ([]coa.Arrangement, error) {
	query := `SELECT id, user_id, method, name, percentage, amount, description, created_at, updated_at, deleted_at
		FROM tbl_arrangement WHERE user_id = $1 AND deleted_at IS NULL ORDER BY name`
	var list []coa.Arrangement
	err := r.db.SelectContext(ctx, &list, query, userID)
	if err != nil {
		return nil, err
	}
	if list == nil {
		list = []coa.Arrangement{}
	}
	return list, nil
}

func (r *arrangementRepo) Update(ctx context.Context, a *coa.Arrangement) error {
	query := `UPDATE tbl_arrangement SET method = $1, name = $2, percentage = $3, amount = $4, description = $5, updated_at = $6
		WHERE id = $7 AND deleted_at IS NULL`
	_, err := r.db.ExecContext(ctx, query,
		a.Method, a.Name, a.Percentage, a.Amount, a.Description, time.Now(), a.ID,
	)
	return err
}

func (r *arrangementRepo) Delete(ctx context.Context, id string) error {
	query := `UPDATE tbl_arrangement SET deleted_at = $1 WHERE id = $2 AND deleted_at IS NULL`
	_, err := r.db.ExecContext(ctx, query, time.Now(), id)
	return err
}
