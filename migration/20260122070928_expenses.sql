-- +goose Up
-- +goose StatementBegin

CREATE TABLE tbl_expense_type (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    clinic_id UUID NOT NULL REFERENCES tbl_clinic(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by UUID REFERENCES tbl_user(id) ON DELETE SET NULL,
    deleted_at TIMESTAMP NULL,
    deleted_by UUID REFERENCES tbl_user(id) ON DELETE SET NULL
);

CREATE TABLE tbl_expense_category (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    clinic_id UUID NOT NULL REFERENCES tbl_clinic(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by UUID REFERENCES tbl_user(id) ON DELETE SET NULL,
    deleted_at TIMESTAMP NULL,
    deleted_by UUID REFERENCES tbl_user(id) ON DELETE SET NULL
);

CREATE TABLE tbl_expense_category_type (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    clinic_id UUID NOT NULL REFERENCES tbl_clinic(id) ON DELETE CASCADE,
    type_id UUID NOT NULL REFERENCES tbl_expense_type(id),
    category_id UUID NOT NULL REFERENCES tbl_expense_category(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by UUID REFERENCES tbl_user(id) ON DELETE SET NULL,
    deleted_at TIMESTAMP NULL,
    deleted_by UUID REFERENCES tbl_user(id) ON DELETE SET NULL
);

CREATE TABLE tbl_expense_entry (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    clinic_id UUID NOT NULL REFERENCES tbl_clinic(id) ON DELETE CASCADE,
    category_id UUID NOT NULL REFERENCES tbl_expense_category(id),
    type_id UUID NOT NULL REFERENCES tbl_expense_type(id),
    amount DECIMAL(10,2) NOT NULL,
    gst_rate DECIMAL(5,2),
    is_gst_inclusive BOOLEAN DEFAULT FALSE,
    expense_date DATE NOT NULL,
    supplier_name VARCHAR(255) NOT NULL,
    notes TEXT,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    created_by UUID REFERENCES tbl_user(id) ON DELETE SET NULL,
    deleted_at TIMESTAMP NULL,
    deleted_by UUID REFERENCES tbl_user(id) ON DELETE SET NULL
);

-- -------------------
-- UNIQUE (SOFT DELETE SAFE)
-- -------------------

CREATE UNIQUE INDEX ux_expense_type_clinic_name_active
ON tbl_expense_type (clinic_id, name)
WHERE deleted_at IS NULL;

CREATE UNIQUE INDEX ux_expense_category_clinic_name_active
ON tbl_expense_category (clinic_id, name)
WHERE deleted_at IS NULL;

CREATE UNIQUE INDEX ux_expense_cat_type_active
ON tbl_expense_category_type (clinic_id, type_id, category_id)
WHERE deleted_at IS NULL;

-- -------------------
-- FILTER & JOIN INDEXES
-- -------------------

CREATE INDEX idx_expense_type_active
ON tbl_expense_type (clinic_id)
WHERE deleted_at IS NULL;

CREATE INDEX idx_expense_category_active
ON tbl_expense_category (clinic_id)
WHERE deleted_at IS NULL;

CREATE INDEX idx_expense_cat_type_active_lookup
ON tbl_expense_category_type (clinic_id, type_id, category_id)
WHERE deleted_at IS NULL;

CREATE INDEX idx_expense_entry_active_clinic_date
ON tbl_expense_entry (clinic_id, expense_date)
WHERE deleted_at IS NULL;

CREATE INDEX idx_expense_entry_category
ON tbl_expense_entry (category_id)
WHERE deleted_at IS NULL;

CREATE INDEX idx_expense_entry_type
ON tbl_expense_entry (type_id)
WHERE deleted_at IS NULL;

CREATE INDEX idx_expense_entry_gst
ON tbl_expense_entry (gst_rate)
WHERE gst_rate IS NOT NULL
  AND deleted_at IS NULL;

CREATE INDEX idx_expense_entry_created_by
ON tbl_expense_entry (created_by)
WHERE deleted_at IS NULL;

-- +goose StatementEnd


-- +goose Down
-- +goose StatementBegin
DROP TABLE tbl_expense_entry;
DROP TABLE tbl_expense_category_type;
DROP TABLE tbl_expense_category;
DROP TABLE tbl_expense_type;
-- +goose StatementEnd
