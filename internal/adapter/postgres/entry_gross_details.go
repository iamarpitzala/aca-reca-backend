package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/application/port"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

type entryGrossDetailsRepo struct {
	db *sqlx.DB
}

func NewEntryGrossDetailsRepository(db *sqlx.DB) port.EntryGrossDetailsRepository {
	return &entryGrossDetailsRepo{db: db}
}

func (r *entryGrossDetailsRepo) Create(ctx context.Context, grossDetails *domain.EntryGrossDetails) error {
	q := `INSERT INTO tbl_entry_gross_details (
		id, source_entry_id, service_facility_fee_percent, service_fee_base,
		gst_on_service_fee, total_service_fee, net_amount, created_at, updated_at
	) VALUES (
		:id, :source_entry_id, :service_facility_fee_percent, :service_fee_base,
		:gst_on_service_fee, :total_service_fee, :net_amount, :created_at, :updated_at
	)`
	_, err := r.db.NamedExecContext(ctx, q, grossDetails)
	return err
}

func (r *entryGrossDetailsRepo) GetByEntryID(ctx context.Context, entryID uuid.UUID) (*domain.EntryGrossDetails, error) {
	q := `
		SELECT
			id, source_entry_id, service_facility_fee_percent, service_fee_base,
			gst_on_service_fee, total_service_fee, net_amount, created_at, updated_at
		FROM tbl_entry_gross_details
		WHERE source_entry_id = $1
	`
	var grossDetails domain.EntryGrossDetails
	err := r.db.GetContext(ctx, &grossDetails, q, entryID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("gross details not found for entry: %w", err)
		}
		return nil, fmt.Errorf("failed to get gross details: %w", err)
	}
	return &grossDetails, nil
}

func (r *entryGrossDetailsRepo) Update(ctx context.Context, grossDetails *domain.EntryGrossDetails) error {
	q := `
		UPDATE tbl_entry_gross_details SET
			service_facility_fee_percent = :service_facility_fee_percent,
			service_fee_base = :service_fee_base,
			gst_on_service_fee = :gst_on_service_fee,
			total_service_fee = :total_service_fee,
			net_amount = :net_amount,
			updated_at = :updated_at
		WHERE source_entry_id = :source_entry_id
	`
	_, err := r.db.NamedExecContext(ctx, q, grossDetails)
	return err
}

func (r *entryGrossDetailsRepo) Delete(ctx context.Context, entryID uuid.UUID) error {
	q := `DELETE FROM tbl_entry_gross_details WHERE source_entry_id = $1`
	_, err := r.db.ExecContext(ctx, q, entryID)
	return err
}
