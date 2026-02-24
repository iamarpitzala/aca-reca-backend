package form

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// FormulaOperator is the allowed operator for formula fields.
const (
	FormulaOpAdd = "+"
	FormulaOpSub = "-"
	FormulaOpMul = "*"
	FormulaOpDiv = "/"
)

// SourceRole for formula source (PRIMARY, SECONDARY).
const (
	SourceRolePrimary   = "PRIMARY"
	SourceRoleSecondary = "SECONDARY"
)

type FieldConfigRequest struct {
	ID           *string `json:"id" validate:"omitempty,required"`
	FormFieldID  string  `json:"formFieldId" validate:"required"`
	TaxTypeID     int     `json:"taxTypeId" validate:"required"`
	ArrangementID *string `json:"arrangementId" validate:"omitempty"`
	IsFormula     bool    `json:"isFormula"`
	Operator      *string `json:"operator" validate:"omitempty,oneof=+ - * /"`
}

func (r *FieldConfigRequest) Validate() error {
	if r.FormFieldID == "" {
		return errors.New("formFieldId is required")
	}
	if r.IsFormula && (r.Operator == nil || *r.Operator == "") {
		return errors.New("operator is required when isFormula is true")
	}
	if r.Operator != nil {
		op := *r.Operator
		if op != FormulaOpAdd && op != FormulaOpSub && op != FormulaOpMul && op != FormulaOpDiv {
			return errors.New("operator must be one of +, -, *, /")
		}
	}
	return nil
}

type FieldConfig struct {
	ID            string     `db:"id"`
	FormFieldID   string     `db:"form_field_id"`
	TaxTypeID     int        `db:"tax_type_id"`
	ArrangementID *string    `db:"arrangement_id"`
	IsFormula     bool       `db:"is_formula"`
	Operator      *string    `db:"operator"`
	CreatedAt     time.Time  `db:"created_at"`
	UpdatedAt     time.Time  `db:"updated_at"`
	DeletedAt     *time.Time `db:"deleted_at"`
}

func (c *FieldConfig) FromRequest(req *FieldConfigRequest) {
	if req.ID != nil && *req.ID != "" {
		c.ID = *req.ID
	} else {
		c.ID = uuid.New().String()
	}
	c.FormFieldID = req.FormFieldID
	c.TaxTypeID = req.TaxTypeID
	c.ArrangementID = req.ArrangementID
	c.IsFormula = req.IsFormula
	c.Operator = req.Operator
}

type FieldConfigResponse struct {
	ID            string     `json:"id"`
	FormFieldID   string     `json:"formFieldId"`
	TaxTypeID     int        `json:"taxTypeId"`
	ArrangementID *string    `json:"arrangementId,omitempty"`
	IsFormula     bool       `json:"isFormula"`
	Operator      *string    `json:"operator,omitempty"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
	DeletedAt     *time.Time `json:"deletedAt,omitempty"`
}

func (c *FieldConfig) ToResponse() *FieldConfigResponse {
	return &FieldConfigResponse{
		ID:            c.ID,
		FormFieldID:   c.FormFieldID,
		TaxTypeID:     c.TaxTypeID,
		ArrangementID: c.ArrangementID,
		IsFormula:     c.IsFormula,
		Operator:      c.Operator,
		CreatedAt:     c.CreatedAt,
		UpdatedAt:     c.UpdatedAt,
		DeletedAt:     c.DeletedAt,
	}
}

type FormulaSourceRequest struct {
	ID            *string `json:"id" validate:"omitempty,required"`
	FieldConfigID string  `json:"fieldConfigId" validate:"required"`
	SourceFieldID string  `json:"sourceFieldId" validate:"required"`
	SourceRole    string  `json:"sourceRole" validate:"required,oneof=PRIMARY SECONDARY"`
	SourceOrder   int     `json:"sourceOrder" validate:"required,min=0"`
}

func (r *FormulaSourceRequest) Validate() error {
	if r.FieldConfigID == "" {
		return errors.New("fieldConfigId is required")
	}
	if r.SourceFieldID == "" {
		return errors.New("sourceFieldId is required")
	}
	if r.SourceRole != SourceRolePrimary && r.SourceRole != SourceRoleSecondary {
		return errors.New("sourceRole must be PRIMARY or SECONDARY")
	}
	return nil
}

type FormulaSource struct {
	ID            string     `db:"id"`
	FieldConfigID string     `db:"field_config_id"`
	SourceFieldID string     `db:"source_field_id"`
	SourceRole    string     `db:"source_role"`
	SourceOrder   int        `db:"source_order"`
	CreatedAt     time.Time  `db:"created_at"`
	UpdatedAt     time.Time  `db:"updated_at"`
	DeletedAt     *time.Time `db:"deleted_at"`
}

func (s *FormulaSource) FromRequest(req *FormulaSourceRequest) {
	if req.ID != nil && *req.ID != "" {
		s.ID = *req.ID
	} else {
		s.ID = uuid.New().String()
	}
	s.FieldConfigID = req.FieldConfigID
	s.SourceFieldID = req.SourceFieldID
	s.SourceRole = req.SourceRole
	s.SourceOrder = req.SourceOrder
}

type FormulaSourceResponse struct {
	ID            string     `json:"id"`
	FieldConfigID string     `json:"fieldConfigId"`
	SourceFieldID string     `json:"sourceFieldId"`
	SourceRole    string     `json:"sourceRole"`
	SourceOrder   int        `json:"sourceOrder"`
	CreatedAt     time.Time  `json:"createdAt"`
	UpdatedAt     time.Time  `json:"updatedAt"`
	DeletedAt     *time.Time `json:"deletedAt,omitempty"`
}

func (s *FormulaSource) ToResponse() *FormulaSourceResponse {
	return &FormulaSourceResponse{
		ID:            s.ID,
		FieldConfigID: s.FieldConfigID,
		SourceFieldID: s.SourceFieldID,
		SourceRole:    s.SourceRole,
		SourceOrder:   s.SourceOrder,
		CreatedAt:     s.CreatedAt,
		UpdatedAt:     s.UpdatedAt,
		DeletedAt:     s.DeletedAt,
	}
}
