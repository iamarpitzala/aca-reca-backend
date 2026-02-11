package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
	"github.com/jmoiron/sqlx"
)

// saveNormalizedEntry saves an entry in normalized tables
func (r *customFormRepo) saveNormalizedEntry(ctx context.Context, tx *sqlx.Tx, entry *domain.CustomFormEntry, form *domain.CustomForm) error {
	// Convert JSONB entry to normalized structure
	normalized, err := domain.ConvertJSONBToNormalized(entry, form)
	if err != nil {
		return fmt.Errorf("failed to convert entry to normalized: %w", err)
	}

	// Save header
	headerQ := `INSERT INTO tbl_entry_header (id, form_id, form_name, form_type, calculation_method, clinic_id, quarter_id, entry_date, description, remarks, payment_responsibility, created_by, created_at, updated_at)
		VALUES (:id, :form_id, :form_name, :form_type, :calculation_method, :clinic_id, :quarter_id, :entry_date, :description, :remarks, :payment_responsibility, :created_by, :created_at, :updated_at)`
	if _, err := tx.NamedExecContext(ctx, headerQ, normalized.Header); err != nil {
		return fmt.Errorf("failed to insert entry header: %w", err)
	}

	// Save field values
	if len(normalized.FieldValues) > 0 {
		for _, fv := range normalized.FieldValues {
			fv.ID = uuid.New()
			fvQ := `INSERT INTO tbl_entry_field_value (id, entry_id, field_id, field_name, value, text_value, boolean_value, manual_gst_amount, display_order, created_at)
				VALUES (:id, :entry_id, :field_id, :field_name, :value, :text_value, :boolean_value, :manual_gst_amount, :display_order, :created_at)`
			if _, err := tx.NamedExecContext(ctx, fvQ, fv); err != nil {
				return fmt.Errorf("failed to insert field value: %w", err)
			}
		}
	}

	// Save field calculations
	if len(normalized.FieldCalculations) > 0 {
		for _, fc := range normalized.FieldCalculations {
			fc.ID = uuid.New()
			fcQ := `INSERT INTO tbl_entry_field_calculation (id, entry_id, field_id, field_name, base_amount, gst_amount, total_amount, gst_rate, gst_type, section, payment_responsibility, display_order, created_at)
				VALUES (:id, :entry_id, :field_id, :field_name, :base_amount, :gst_amount, :total_amount, :gst_rate, :gst_type, :section, :payment_responsibility, :display_order, :created_at)`
			if _, err := tx.NamedExecContext(ctx, fcQ, fc); err != nil {
				return fmt.Errorf("failed to insert field calculation: %w", err)
			}
		}
	}

	// Save summary
	if normalized.Summary != nil {
		normalized.Summary.ID = uuid.New()
		summaryQ := `INSERT INTO tbl_entry_summary (id, entry_id, total_base_amount, total_gst_amount, total_amount, net_payable, net_receivable, net_fee, bas_gst_on_sales_1a, bas_gst_credit_1b, bas_total_sales_g1, bas_expenses_g11, created_at, updated_at)
			VALUES (:id, :entry_id, :total_base_amount, :total_gst_amount, :total_amount, :net_payable, :net_receivable, :net_fee, :bas_gst_on_sales_1a, :bas_gst_credit_1b, :bas_total_sales_g1, :bas_expenses_g11, :created_at, :updated_at)`
		if _, err := tx.NamedExecContext(ctx, summaryQ, normalized.Summary); err != nil {
			return fmt.Errorf("failed to insert entry summary: %w", err)
		}
	}

	// Save net details
	if normalized.NetDetails != nil {
		normalized.NetDetails.ID = uuid.New()
		netQ := `INSERT INTO tbl_entry_net_details (id, entry_id, commission_percent, commission, gst_on_commission, total_payment_received, super_holding_enabled, super_component_percent, commission_component, super_component, total_for_reconciliation, created_at, updated_at)
			VALUES (:id, :entry_id, :commission_percent, :commission, :gst_on_commission, :total_payment_received, :super_holding_enabled, :super_component_percent, :commission_component, :super_component, :total_for_reconciliation, :created_at, :updated_at)`
		if _, err := tx.NamedExecContext(ctx, netQ, normalized.NetDetails); err != nil {
			return fmt.Errorf("failed to insert net details: %w", err)
		}
	}

	// Save gross details
	if normalized.GrossDetails != nil {
		normalized.GrossDetails.ID = uuid.New()
		grossQ := `INSERT INTO tbl_entry_gross_details (id, entry_id, service_facility_fee_percent, service_fee_base, gst_on_service_fee, total_service_fee, subtotal_after_deductions, remitted_amount, created_at, updated_at)
			VALUES (:id, :entry_id, :service_facility_fee_percent, :service_fee_base, :gst_on_service_fee, :total_service_fee, :subtotal_after_deductions, :remitted_amount, :created_at, :updated_at)`
		if _, err := tx.NamedExecContext(ctx, grossQ, normalized.GrossDetails); err != nil {
			return fmt.Errorf("failed to insert gross details: %w", err)
		}
	}

	// Save gross reductions
	if len(normalized.GrossReductions) > 0 {
		for _, gr := range normalized.GrossReductions {
			gr.ID = uuid.New()
			grQ := `INSERT INTO tbl_entry_gross_reduction (id, entry_id, field_calculation_id, field_id, field_name, base_amount, gst_amount, total_amount, display_order, created_at)
				VALUES (:id, :entry_id, :field_calculation_id, :field_id, :field_name, :base_amount, :gst_amount, :total_amount, :display_order, :created_at)`
			if _, err := tx.NamedExecContext(ctx, grQ, gr); err != nil {
				return fmt.Errorf("failed to insert gross reduction: %w", err)
			}
		}
	}

	// Save gross reimbursements
	if len(normalized.GrossReimbursements) > 0 {
		for _, gr := range normalized.GrossReimbursements {
			gr.ID = uuid.New()
			grQ := `INSERT INTO tbl_entry_gross_reimbursement (id, entry_id, field_calculation_id, field_id, field_name, base_amount, gst_amount, total_amount, display_order, created_at)
				VALUES (:id, :entry_id, :field_calculation_id, :field_id, :field_name, :base_amount, :gst_amount, :total_amount, :display_order, :created_at)`
			if _, err := tx.NamedExecContext(ctx, grQ, gr); err != nil {
				return fmt.Errorf("failed to insert gross reimbursement: %w", err)
			}
		}
	}

	// Save gross additional reductions
	if len(normalized.GrossAdditionalReductions) > 0 {
		for _, gar := range normalized.GrossAdditionalReductions {
			gar.ID = uuid.New()
			garQ := `INSERT INTO tbl_entry_gross_additional_reduction (id, entry_id, field_calculation_id, field_id, field_name, base_amount, gst_amount, total_amount, display_order, created_at)
				VALUES (:id, :entry_id, :field_calculation_id, :field_id, :field_name, :base_amount, :gst_amount, :total_amount, :display_order, :created_at)`
			if _, err := tx.NamedExecContext(ctx, garQ, gar); err != nil {
				return fmt.Errorf("failed to insert gross additional reduction: %w", err)
			}
		}
	}

	// Save gross reductions summary
	if normalized.GrossReductionsSummary != nil {
		normalized.GrossReductionsSummary.ID = uuid.New()
		grsQ := `INSERT INTO tbl_entry_gross_reductions_summary (id, entry_id, total_reductions, total_reduction_base, total_expense_gst, total_reimbursements, total_additional_reduction, total_additional_reduction_base, total_additional_reduction_gst, created_at, updated_at)
			VALUES (:id, :entry_id, :total_reductions, :total_reduction_base, :total_expense_gst, :total_reimbursements, :total_additional_reduction, :total_additional_reduction_base, :total_additional_reduction_gst, :created_at, :updated_at)`
		if _, err := tx.NamedExecContext(ctx, grsQ, normalized.GrossReductionsSummary); err != nil {
			return fmt.Errorf("failed to insert gross reductions summary: %w", err)
		}
	}

	// Save gross outwork
	if normalized.GrossOutwork != nil {
		normalized.GrossOutwork.ID = uuid.New()
		outworkQ := `INSERT INTO tbl_entry_gross_outwork (id, entry_id, outwork_enabled, outwork_rate_percent, outwork_charge_base, outwork_charge_gst, outwork_charge_total, created_at, updated_at)
			VALUES (:id, :entry_id, :outwork_enabled, :outwork_rate_percent, :outwork_charge_base, :outwork_charge_gst, :outwork_charge_total, :created_at, :updated_at)`
		if _, err := tx.NamedExecContext(ctx, outworkQ, normalized.GrossOutwork); err != nil {
			return fmt.Errorf("failed to insert gross outwork: %w", err)
		}
	}

	// Save deductions
	if normalized.Deductions != nil {
		normalized.Deductions.ID = uuid.New()
		dedQ := `INSERT INTO tbl_entry_deductions (id, entry_id, service_facility_fee_percent, service_fee_override, commission_percent, super_holding_enabled, super_component_percent, outwork_enabled, outwork_rate_percent, entry_payment_responsibility, created_at)
			VALUES (:id, :entry_id, :service_facility_fee_percent, :service_fee_override, :commission_percent, :super_holding_enabled, :super_component_percent, :outwork_enabled, :outwork_rate_percent, :entry_payment_responsibility, :created_at)`
		if _, err := tx.NamedExecContext(ctx, dedQ, normalized.Deductions); err != nil {
			return fmt.Errorf("failed to insert deductions: %w", err)
		}
	}

	return nil
}

// loadNormalizedEntry loads an entry from normalized tables
func (r *customFormRepo) loadNormalizedEntry(ctx context.Context, tx *sqlx.Tx, entryID uuid.UUID) (*domain.NormalizedEntry, error) {
	normalized := &domain.NormalizedEntry{}

	// Load header
	headerQ := `SELECT id, form_id, form_name, form_type, calculation_method, clinic_id, quarter_id, entry_date, description, remarks, payment_responsibility, created_by, created_at, updated_at, deleted_at, original_entry_id
		FROM tbl_entry_header WHERE id = $1 AND deleted_at IS NULL`
	var header domain.EntryHeader
	if err := tx.GetContext(ctx, &header, headerQ, entryID); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("entry not found")
		}
		return nil, fmt.Errorf("failed to load entry header: %w", err)
	}
	normalized.Header = &header

	// Get form to get calculation method (we'll use it later for conversion)
	formQ := `SELECT calculation_method FROM tbl_custom_form WHERE id = $1`
	var calcMethod string
	if err := tx.GetContext(ctx, &calcMethod, formQ, header.FormID); err != nil {
		// If form doesn't exist, still try to load entry data
		calcMethod = header.CalculationMethod
	}

	// Load field values
	valuesQ := `SELECT id, entry_id, field_id, field_name, value, text_value, boolean_value, manual_gst_amount, display_order, created_at
		FROM tbl_entry_field_value WHERE entry_id = $1 ORDER BY display_order`
	var fieldValues []domain.EntryFieldValue
	if err := tx.SelectContext(ctx, &fieldValues, valuesQ, entryID); err != nil {
		return nil, fmt.Errorf("failed to load field values: %w", err)
	}
	normalized.FieldValues = fieldValues

	// Load field calculations
	calcQ := `SELECT id, entry_id, field_id, field_name, base_amount, gst_amount, total_amount, gst_rate, gst_type, section, payment_responsibility, display_order, created_at
		FROM tbl_entry_field_calculation WHERE entry_id = $1 ORDER BY display_order`
	var fieldCalcs []domain.EntryFieldCalculation
	if err := tx.SelectContext(ctx, &fieldCalcs, calcQ, entryID); err != nil {
		return nil, fmt.Errorf("failed to load field calculations: %w", err)
	}
	normalized.FieldCalculations = fieldCalcs

	// Load summary
	summaryQ := `SELECT id, entry_id, total_base_amount, total_gst_amount, total_amount, net_payable, net_receivable, net_fee, bas_gst_on_sales_1a, bas_gst_credit_1b, bas_total_sales_g1, bas_expenses_g11, created_at, updated_at
		FROM tbl_entry_summary WHERE entry_id = $1`
	var summary domain.EntrySummary
	if err := tx.GetContext(ctx, &summary, summaryQ, entryID); err != nil {
		if err != sql.ErrNoRows {
			return nil, fmt.Errorf("failed to load entry summary: %w", err)
		}
	} else {
		normalized.Summary = &summary
	}

	// Load net details
	if calcMethod == "net" {
		netQ := `SELECT id, entry_id, commission_percent, commission, gst_on_commission, total_payment_received, super_holding_enabled, super_component_percent, commission_component, super_component, total_for_reconciliation, created_at, updated_at
			FROM tbl_entry_net_details WHERE entry_id = $1`
		var netDetails domain.EntryNetDetails
		if err := tx.GetContext(ctx, &netDetails, netQ, entryID); err != nil {
			if err != sql.ErrNoRows {
				return nil, fmt.Errorf("failed to load net details: %w", err)
			}
		} else {
			normalized.NetDetails = &netDetails
		}
	}

	// Load gross details
	if calcMethod == "gross" {
		grossQ := `SELECT id, entry_id, service_facility_fee_percent, service_fee_base, gst_on_service_fee, total_service_fee, subtotal_after_deductions, remitted_amount, created_at, updated_at
			FROM tbl_entry_gross_details WHERE entry_id = $1`
		var grossDetails domain.EntryGrossDetails
		if err := tx.GetContext(ctx, &grossDetails, grossQ, entryID); err != nil {
			if err != sql.ErrNoRows {
				return nil, fmt.Errorf("failed to load gross details: %w", err)
			}
		} else {
			normalized.GrossDetails = &grossDetails
		}

		// Load gross reductions
		reductionQ := `SELECT id, entry_id, field_calculation_id, field_id, field_name, base_amount, gst_amount, total_amount, display_order, created_at
			FROM tbl_entry_gross_reduction WHERE entry_id = $1 ORDER BY display_order`
		var reductions []domain.EntryGrossReduction
		if err := tx.SelectContext(ctx, &reductions, reductionQ, entryID); err == nil {
			normalized.GrossReductions = reductions
		}

		// Load gross reimbursements
		reimbQ := `SELECT id, entry_id, field_calculation_id, field_id, field_name, base_amount, gst_amount, total_amount, display_order, created_at
			FROM tbl_entry_gross_reimbursement WHERE entry_id = $1 ORDER BY display_order`
		var reimbursements []domain.EntryGrossReimbursement
		if err := tx.SelectContext(ctx, &reimbursements, reimbQ, entryID); err == nil {
			normalized.GrossReimbursements = reimbursements
		}

		// Load gross additional reductions
		addlRedQ := `SELECT id, entry_id, field_calculation_id, field_id, field_name, base_amount, gst_amount, total_amount, display_order, created_at
			FROM tbl_entry_gross_additional_reduction WHERE entry_id = $1 ORDER BY display_order`
		var addlReductions []domain.EntryGrossAdditionalReduction
		if err := tx.SelectContext(ctx, &addlReductions, addlRedQ, entryID); err == nil {
			normalized.GrossAdditionalReductions = addlReductions
		}

		// Load gross reductions summary
		reductionSummaryQ := `SELECT id, entry_id, total_reductions, total_reduction_base, total_expense_gst, total_reimbursements, total_additional_reduction, total_additional_reduction_base, total_additional_reduction_gst, created_at, updated_at
			FROM tbl_entry_gross_reductions_summary WHERE entry_id = $1`
		var reductionSummary domain.EntryGrossReductionsSummary
		if err := tx.GetContext(ctx, &reductionSummary, reductionSummaryQ, entryID); err == nil {
			normalized.GrossReductionsSummary = &reductionSummary
		}

		// Load gross outwork
		outworkQ := `SELECT id, entry_id, outwork_enabled, outwork_rate_percent, outwork_charge_base, outwork_charge_gst, outwork_charge_total, created_at, updated_at
			FROM tbl_entry_gross_outwork WHERE entry_id = $1`
		var outwork domain.EntryGrossOutwork
		if err := tx.GetContext(ctx, &outwork, outworkQ, entryID); err == nil {
			normalized.GrossOutwork = &outwork
		}
	}

	// Load deductions
	dedQ := `SELECT id, entry_id, service_facility_fee_percent, service_fee_override, commission_percent, super_holding_enabled, super_component_percent, outwork_enabled, outwork_rate_percent, entry_payment_responsibility, created_at
		FROM tbl_entry_deductions WHERE entry_id = $1`
	var deductions domain.EntryDeductions
	if err := tx.GetContext(ctx, &deductions, dedQ, entryID); err == nil {
		normalized.Deductions = &deductions
	}

	return normalized, nil
}

// deleteNormalizedEntry deletes an entry from normalized tables (soft delete)
func (r *customFormRepo) deleteNormalizedEntry(ctx context.Context, tx *sqlx.Tx, entryID uuid.UUID) error {
	_, err := tx.ExecContext(ctx, `UPDATE tbl_entry_header SET deleted_at = $1 WHERE id = $2`, time.Now(), entryID)
	return err
}

// updateNormalizedEntry updates an entry in normalized tables
func (r *customFormRepo) updateNormalizedEntry(ctx context.Context, tx *sqlx.Tx, entry *domain.CustomFormEntry, form *domain.CustomForm) error {
	// Delete existing normalized data (except header)
	entryID := entry.ID
	tx.ExecContext(ctx, `DELETE FROM tbl_entry_field_value WHERE entry_id = $1`, entryID)
	tx.ExecContext(ctx, `DELETE FROM tbl_entry_field_calculation WHERE entry_id = $1`, entryID)
	tx.ExecContext(ctx, `DELETE FROM tbl_entry_summary WHERE entry_id = $1`, entryID)
	tx.ExecContext(ctx, `DELETE FROM tbl_entry_net_details WHERE entry_id = $1`, entryID)
	tx.ExecContext(ctx, `DELETE FROM tbl_entry_gross_details WHERE entry_id = $1`, entryID)
	tx.ExecContext(ctx, `DELETE FROM tbl_entry_gross_reduction WHERE entry_id = $1`, entryID)
	tx.ExecContext(ctx, `DELETE FROM tbl_entry_gross_reimbursement WHERE entry_id = $1`, entryID)
	tx.ExecContext(ctx, `DELETE FROM tbl_entry_gross_additional_reduction WHERE entry_id = $1`, entryID)
	tx.ExecContext(ctx, `DELETE FROM tbl_entry_gross_reductions_summary WHERE entry_id = $1`, entryID)
	tx.ExecContext(ctx, `DELETE FROM tbl_entry_gross_outwork WHERE entry_id = $1`, entryID)
	tx.ExecContext(ctx, `DELETE FROM tbl_entry_deductions WHERE entry_id = $1`, entryID)

	// Update header
	headerQ := `UPDATE tbl_entry_header SET updated_at = $1 WHERE id = $2`
	if _, err := tx.ExecContext(ctx, headerQ, entry.UpdatedAt, entryID); err != nil {
		return fmt.Errorf("failed to update entry header: %w", err)
	}

	// Re-save normalized data
	return r.saveNormalizedEntry(ctx, tx, entry, form)
}
