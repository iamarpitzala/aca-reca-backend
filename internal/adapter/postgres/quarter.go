package postgres

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/application/port"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

type quarterRepo struct {
	db *sqlx.DB
}

func NewQuarterRepository(db *sqlx.DB) port.QuarterRepository {
	return &quarterRepo{db: db}
}

func (r *quarterRepo) Create(ctx context.Context, q *domain.Quarter) error {
	query := `INSERT INTO tbl_quarter (id, name, start_date, end_date, created_at, updated_at) VALUES (:id, :name, :start_date, :end_date, :created_at, :updated_at)`
	args := map[string]interface{}{
		"id":         q.ID,
		"name":       q.Name,
		"start_date": q.StartDate,
		"end_date":   q.EndDate,
		"created_at": q.CreatedAt,
		"updated_at": q.UpdatedAt,
	}
	_, err := r.db.NamedExecContext(ctx, query, args)
	return err
}

func (r *quarterRepo) Update(ctx context.Context, q *domain.Quarter) error {
	query := `UPDATE tbl_quarter SET name = :name, start_date = :start_date, end_date = :end_date, created_at = :created_at, updated_at = :updated_at WHERE id = :id`
	args := map[string]interface{}{
		"id":         q.ID,
		"name":       q.Name,
		"start_date": q.StartDate,
		"end_date":   q.EndDate,
		"created_at": q.CreatedAt,
		"updated_at": q.UpdatedAt,
	}
	_, err := r.db.NamedExecContext(ctx, query, args)
	return err
}

func (r *quarterRepo) Delete(ctx context.Context, id uuid.UUID) error {
	query := `UPDATE tbl_quarter SET deleted_at = $1 WHERE id = $2`
	_, err := r.db.ExecContext(ctx, query, time.Now(), id)
	return err
}

func (r *quarterRepo) List(ctx context.Context) ([]domain.Quarter, error) {
	type row struct {
		ID        uuid.UUID  `db:"id"`
		Name      string     `db:"name"`
		StartDate time.Time  `db:"start_date"`
		EndDate   time.Time  `db:"end_date"`
		CreatedAt time.Time  `db:"created_at"`
		UpdatedAt *time.Time `db:"updated_at"`
		DeletedAt *time.Time `db:"deleted_at"`
	}
	query := `SELECT id, name, start_date, end_date, created_at, updated_at, deleted_at FROM tbl_quarter WHERE deleted_at IS NULL ORDER BY created_at DESC`
	var rows []row
	if err := r.db.SelectContext(ctx, &rows, query); err != nil {
		return nil, errors.New("failed to get quarters")
	}
	out := make([]domain.Quarter, len(rows))
	for i, row := range rows {
		out[i] = domain.Quarter{
			ID:        row.ID,
			Name:      row.Name,
			StartDate: row.StartDate,
			EndDate:   row.EndDate,
			CreatedAt: row.CreatedAt,
			UpdatedAt: row.UpdatedAt,
			DeletedAt: row.DeletedAt,
		}
	}
	return out, nil
}

func (r *quarterRepo) GetByID(ctx context.Context, id uuid.UUID) (domain.Quarter, error) {
	query := `SELECT id, name, start_date, end_date, created_at, updated_at, deleted_at FROM tbl_quarter WHERE id = $1 AND deleted_at IS NULL`
	var q domain.Quarter
	if err := r.db.GetContext(ctx, &q, query, id); err != nil {
		return domain.Quarter{}, err
	}
	return q, nil
}
