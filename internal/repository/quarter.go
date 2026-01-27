package repository

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

func CreateQuarter(ctx context.Context, db *sqlx.DB, quart *domain.Quarter) error {
	query := `INSERT INTO tbl_quarter (id, name, start_date, end_date, created_at, updated_at)
		VALUES (:id, :name, :start_date, :start_date, :created_at, :updated_at)`

	args := map[string]interface{}{
		"id":         quart.ID,
		"name":       quart.Name,
		"start_date": quart.StartDate,
		"end_date":   quart.EndDate,
		"created_at": quart.CreatedAt,
		"updated_at": quart.UpdatedAt,
	}

	_, err := db.NamedExecContext(ctx, query, args)
	if err != nil {
		return err
	}
	return nil
}

func UpdateQuarter(ctx context.Context, db *sqlx.DB, quart *domain.Quarter) error {
	query := `UPDATE tbl_quarter SET name = :name, start_date = :start_date, end_date = :end_date, 
		created_at = :created_at, updated_at = :updated_at WHERE id = :id`

	args := map[string]interface{}{
		"id":         quart.ID,
		"name":       quart.Name,
		"start_date": quart.StartDate,
		"end_date":   quart.EndDate,
		"created_at": quart.CreatedAt,
		"updated_at": quart.UpdatedAt,
	}

	_, err := db.NamedExecContext(ctx, query, args)
	if err != nil {
		return err
	}
	return nil
}

func DeleteQuarter(ctx context.Context, db *sqlx.DB, id uuid.UUID) error {
	query := `UPDATE tbl_quarter SET deleted_at = $1 WHERE id = $2`
	_, err := db.ExecContext(ctx, query, time.Now(), id)
	if err != nil {
		return err
	}
	return nil
}

func ListQuarter(ctx context.Context, db *sqlx.DB) ([]domain.Quarter, error) {
	type formRow struct {
		ID        uuid.UUID `db:"id"`
		Name      string    `db:"name"`
		StartDate time.Time `db:"start_date"`
		EndDate   time.Time `db:"end_date"`

		CreatedAt time.Time  `db:"created_at"`
		UpdatedAt *time.Time `db:"updated_at"`
		DeletedAt *time.Time `db:"deleted_at"`
	}

	query := `SELECT id, name, start_date, end_date, created_at, updated_at, deleted_at
		FROM tbl_quarter WHERE deleted_at IS NULL ORDER BY created_at DESC`

	var rows []formRow
	err := db.SelectContext(ctx, &rows, query)
	if err != nil {
		return nil, errors.New("failed to get financial forms")
	}

	forms := make([]domain.Quarter, len(rows))
	for i, row := range rows {
		forms[i] = domain.Quarter{
			ID:        row.ID,
			Name:      row.Name,
			StartDate: row.StartDate,
			EndDate:   row.EndDate,
			CreatedAt: row.CreatedAt,
			UpdatedAt: row.UpdatedAt,
			DeletedAt: row.DeletedAt,
		}
	}

	return forms, nil
}

func GetQuarterByID(ctx context.Context, db *sqlx.DB, id uuid.UUID) (domain.Quarter, error) {
	query := `
		SELECT id, name, start_date, end_date, created_at, updated_at, deleted_at
		FROM tbl_quarter
		WHERE id = $1 AND deleted_at IS NULL
	`

	var quart domain.Quarter
	if err := db.GetContext(ctx, &quart, query, id); err != nil {
		return domain.Quarter{}, err
	}

	return quart, nil
}
