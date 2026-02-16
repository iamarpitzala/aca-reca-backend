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

type entryNetDetailsRepo struct {
	db *sqlx.DB
}

func NewEntryNetDetailsRepository(db *sqlx.DB) port.EntryNetDetailsRepository {
	return &entryNetDetailsRepo{db: db}
}

func (r *entryNetDetailsRepo) Create(ctx context.Context, netDetails *domain.EntryNetDetails) error {
	q := `INSERT INTO tbl_entry_net_details (
		id, source_entry_id, commission_percent, commission, gst_on_commission,
		total_payment_received, super_holding_enabled, super_component_percent,
		commission_component, super_component, total_for_reconciliation,
		created_at, updated_at
	) VALUES (
		:id, :source_entry_id, :commission_percent, :commission, :gst_on_commission,
		:total_payment_received, :super_holding_enabled, :super_component_percent,
		:commission_component, :super_component, :total_for_reconciliation,
		:created_at, :updated_at
	)`
	_, err := r.db.NamedExecContext(ctx, q, netDetails)
	return err
}

func (r *entryNetDetailsRepo) GetByEntryID(ctx context.Context, entryID uuid.UUID) (*domain.EntryNetDetails, error) {
	q := `
		SELECT
			id, source_entry_id, commission_percent, commission, gst_on_commission,
			total_payment_received, super_holding_enabled, super_component_percent,
			commission_component, super_component, total_for_reconciliation,
			created_at, updated_at
		FROM tbl_entry_net_details
		WHERE source_entry_id = $1
	`
	var netDetails domain.EntryNetDetails
	err := r.db.GetContext(ctx, &netDetails, q, entryID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("net details not found for entry: %w", err)
		}
		return nil, fmt.Errorf("failed to get net details: %w", err)
	}
	return &netDetails, nil
}

func (r *entryNetDetailsRepo) Update(ctx context.Context, netDetails *domain.EntryNetDetails) error {
	q := `
		UPDATE tbl_entry_net_details SET
			commission_percent = :commission_percent,
			commission = :commission,
			gst_on_commission = :gst_on_commission,
			total_payment_received = :total_payment_received,
			super_holding_enabled = :super_holding_enabled,
			super_component_percent = :super_component_percent,
			commission_component = :commission_component,
			super_component = :super_component,
			total_for_reconciliation = :total_for_reconciliation,
			updated_at = :updated_at
		WHERE source_entry_id = :source_entry_id
	`
	_, err := r.db.NamedExecContext(ctx, q, netDetails)
	return err
}

func (r *entryNetDetailsRepo) Delete(ctx context.Context, entryID uuid.UUID) error {
	q := `DELETE FROM tbl_entry_net_details WHERE source_entry_id = $1`
	_, err := r.db.ExecContext(ctx, q, entryID)
	return err
}
