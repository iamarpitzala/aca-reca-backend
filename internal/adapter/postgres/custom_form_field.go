package postgres

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/iamarpitzala/aca-reca-backend/internal/application/port"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain/form"
	"github.com/jmoiron/sqlx"
)

type customFormFieldRepo struct {
	db *sqlx.DB
}

// fieldRow is used for scanning when form_version_id is INTEGER in DB.
type fieldRow struct {
	ID            string     `db:"id"`
	FormVersionID int        `db:"form_version_id"`
	FormID        string     `db:"form_id"`
	StartDate     time.Time  `db:"start_date"`
	EndDate       *time.Time `db:"end_date"`
	Label         string     `db:"label"`
	SectionTypeID int        `db:"section_type_id"`
	Description   *string    `db:"description"`
	IsRequired    bool       `db:"is_required"`
	CoaID         string     `db:"coa_id"`
	Placeholder   *string    `db:"placeholder"`
	MinValue      *float64   `db:"min_value"`
	MaxValue      *float64   `db:"max_value"`
	FieldOrder    int        `db:"field_order"`
	TaxTypeID     *int       `db:"tax_type_id"`
	CreatedAt     time.Time  `db:"created_at"`
	UpdatedAt     time.Time  `db:"updated_at"`
	DeletedAt     *time.Time `db:"deleted_at"`
}

func NewCustomFormFieldRepository(db *sqlx.DB) port.CustomFormFieldRepository {
	return &customFormFieldRepo{db: db}
}

func rowToField(r fieldRow) (*form.Field, error) {
	return &form.Field{
		ID:            r.ID,
		FormVersionID: strconv.Itoa(r.FormVersionID),
		FormID:        r.FormID,
		StartDate:     r.StartDate,
		EndDate:       r.EndDate,
		Label:         r.Label,
		SectionTypeID: r.SectionTypeID,
		Description:   r.Description,
		IsRequired:    r.IsRequired,
		CoaID:         r.CoaID,
		Placeholder:   r.Placeholder,
		MinValue:      r.MinValue,
		MaxValue:      r.MaxValue,
		FieldOrder:    r.FieldOrder,
		TaxTypeID:     r.TaxTypeID,
		CreatedAt:     r.CreatedAt,
		UpdatedAt:     r.UpdatedAt,
		DeletedAt:     r.DeletedAt,
	}, nil
}

func (r *customFormFieldRepo) GetByID(ctx context.Context, id string) (*form.Field, error) {
	q := `
		SELECT id, form_version_id, form_id, start_date, end_date, label, section_type_id,
			description, is_required, coa_id, placeholder, min_value, max_value, field_order,
			tax_type_id, created_at, updated_at, deleted_at
		FROM tbl_custom_form_field
		WHERE id = $1 AND deleted_at IS NULL
	`
	var row fieldRow
	if err := r.db.GetContext(ctx, &row, q, id); err != nil {
		return nil, fmt.Errorf("failed to get custom form field: %w", err)
	}
	return rowToField(row)
}

func (r *customFormFieldRepo) Create(ctx context.Context, field *form.Field) error {
	formVersionID, _ := strconv.Atoi(field.FormVersionID)
	q := `INSERT INTO tbl_custom_form_field (
		id, form_version_id, form_id, start_date, end_date, label, section_type_id,
		description, is_required, coa_id, placeholder, min_value, max_value, field_order,
		tax_type_id, created_at, updated_at
	) VALUES (
		$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15, now(), now()
	)`
	_, err := r.db.ExecContext(ctx, q,
		field.ID, formVersionID, field.FormID, field.StartDate, field.EndDate,
		field.Label, field.SectionTypeID, field.Description, field.IsRequired, field.CoaID,
		field.Placeholder, field.MinValue, field.MaxValue, field.FieldOrder, field.TaxTypeID,
	)
	return err
}

func (r *customFormFieldRepo) CreateBatch(ctx context.Context, fields []*form.Field) error {
	if len(fields) == 0 {
		return nil
	}
	for _, f := range fields {
		if err := r.Create(ctx, f); err != nil {
			return err
		}
	}
	return nil
}

func (r *customFormFieldRepo) Update(ctx context.Context, field *form.Field) error {
	formVersionID, _ := strconv.Atoi(field.FormVersionID)
	q := `UPDATE tbl_custom_form_field SET
		form_version_id = $2, form_id = $3, start_date = $4, end_date = $5, label = $6,
		section_type_id = $7, description = $8, is_required = $9, coa_id = $10,
		placeholder = $11, min_value = $12, max_value = $13, field_order = $14,
		tax_type_id = $15, updated_at = now()
		WHERE id = $1 AND deleted_at IS NULL`
	_, err := r.db.ExecContext(ctx, q,
		field.ID, formVersionID, field.FormID, field.StartDate, field.EndDate,
		field.Label, field.SectionTypeID, field.Description, field.IsRequired, field.CoaID,
		field.Placeholder, field.MinValue, field.MaxValue, field.FieldOrder, field.TaxTypeID,
	)
	return err
}

func (r *customFormFieldRepo) GetByFormID(ctx context.Context, formID string) ([]form.Field, error) {
	q := `
		SELECT id, form_version_id, form_id, start_date, end_date, label, section_type_id,
			description, is_required, coa_id, placeholder, min_value, max_value, field_order,
			tax_type_id, created_at, updated_at, deleted_at
		FROM tbl_custom_form_field
		WHERE form_id = $1 AND deleted_at IS NULL
		ORDER BY field_order ASC
	`
	var rows []fieldRow
	if err := r.db.SelectContext(ctx, &rows, q, formID); err != nil {
		return nil, fmt.Errorf("failed to get custom form fields: %w", err)
	}
	out := make([]form.Field, 0, len(rows))
	for _, row := range rows {
		f, err := rowToField(row)
		if err != nil {
			continue
		}
		out = append(out, *f)
	}
	return out, nil
}

func (r *customFormFieldRepo) GetByFormVersionID(ctx context.Context, formVersionID string) ([]form.Field, error) {
	vid, err := strconv.Atoi(formVersionID)
	if err != nil {
		vid = 0
	}
	q := `
		SELECT id, form_version_id, form_id, start_date, end_date, label, section_type_id,
			description, is_required, coa_id, placeholder, min_value, max_value, field_order,
			tax_type_id, created_at, updated_at, deleted_at
		FROM tbl_custom_form_field
		WHERE form_version_id = $1 AND deleted_at IS NULL
		ORDER BY field_order ASC
	`
	var rows []fieldRow
	if err := r.db.SelectContext(ctx, &rows, q, vid); err != nil {
		return nil, fmt.Errorf("failed to get custom form fields by version: %w", err)
	}
	out := make([]form.Field, 0, len(rows))
	for _, row := range rows {
		f, err := rowToField(row)
		if err != nil {
			continue
		}
		out = append(out, *f)
	}
	return out, nil
}

func (r *customFormFieldRepo) DeleteByFormID(ctx context.Context, formID string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE tbl_custom_form_field SET deleted_at = now() WHERE form_id = $1`, formID)
	return err
}

func (r *customFormFieldRepo) DeleteByFormVersionID(ctx context.Context, formVersionID string) error {
	vid, _ := strconv.Atoi(formVersionID)
	_, err := r.db.ExecContext(ctx, `UPDATE tbl_custom_form_field SET deleted_at = now() WHERE form_version_id = $1`, vid)
	return err
}

func (r *customFormFieldRepo) DeleteByID(ctx context.Context, id string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE tbl_custom_form_field SET deleted_at = now() WHERE id = $1`, id)
	return err
}
