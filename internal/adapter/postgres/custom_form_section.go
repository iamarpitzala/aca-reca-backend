package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"github.com/iamarpitzala/aca-reca-backend/internal/application/port"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain/form"
	"github.com/jmoiron/sqlx"
)

type customFormSectionRepo struct {
	db *sqlx.DB
}

func NewCustomFormSectionRepository(db *sqlx.DB) port.CustomFormSectionRepository {
	return &customFormSectionRepo{db: db}
}

func (r *customFormSectionRepo) Create(ctx context.Context, section *form.Section) error {
	q := `INSERT INTO tbl_custom_form_section (form_version_id, section_type_id, name, description, section_order, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`
	return r.db.QueryRowContext(ctx, q,
		section.FormVersionID, section.SectionTypeID, section.Name, section.Description, section.SectionOrder,
		section.CreatedAt, section.UpdatedAt,
	).Scan(&section.ID)
}

func (r *customFormSectionRepo) GetByID(ctx context.Context, id int) (*form.Section, error) {
	q := `SELECT id, form_version_id, section_type_id, name, description, section_order, created_at, updated_at
		FROM tbl_custom_form_section WHERE id = $1`
	var s form.Section
	if err := r.db.GetContext(ctx, &s, q, id); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, fmt.Errorf("section not found: %w", err)
		}
		return nil, err
	}
	return &s, nil
}

func (r *customFormSectionRepo) GetByFormVersionID(ctx context.Context, formVersionID int) ([]form.Section, error) {
	q := `SELECT id, form_version_id, section_type_id, name, description, section_order, created_at, updated_at
		FROM tbl_custom_form_section WHERE form_version_id = $1 ORDER BY section_order ASC`
	var list []form.Section
	if err := r.db.SelectContext(ctx, &list, q, formVersionID); err != nil {
		return nil, err
	}
	return list, nil
}

func (r *customFormSectionRepo) Update(ctx context.Context, section *form.Section) error {
	q := `UPDATE tbl_custom_form_section SET
		form_version_id = $2, section_type_id = $3, name = $4, description = $5, section_order = $6, updated_at = $7
		WHERE id = $1`
	_, err := r.db.ExecContext(ctx, q,
		section.ID, section.FormVersionID, section.SectionTypeID, section.Name, section.Description, section.SectionOrder, section.UpdatedAt,
	)
	return err
}

func (r *customFormSectionRepo) DeleteByID(ctx context.Context, id int) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM tbl_custom_form_section WHERE id = $1`, id)
	return err
}
