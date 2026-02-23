package form

import "time"

type TaxType struct {
	ID          string     `db:"id"`
	Name        string     `db:"name"`
	Type        string     `db:"type"`
	Description string     `db:"description"`
	CreatedAt   time.Time  `db:"created_at"`
	UpdatedAt   time.Time  `db:"updated_at"`
	DeletedAt   *time.Time `db:"deleted_at"`
}

func (t *TaxType) ToTaxTypeDB(taxType *TaxType) {
	t.ID = taxType.ID
	t.Name = taxType.Name
	t.Type = taxType.Type
	t.Description = taxType.Description
	t.CreatedAt = taxType.CreatedAt
	t.UpdatedAt = taxType.UpdatedAt
	t.DeletedAt = taxType.DeletedAt
}

type TaxTypeResponse struct {
	ID          string     `json:"id"`
	Name        string     `json:"name"`
	Type        string     `json:"type"`
	Description string     `json:"description"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
	DeletedAt   *time.Time `json:"deletedAt"`
}

func (t *TaxType) ToTaxTypeResponse() *TaxTypeResponse {
	return &TaxTypeResponse{
		ID:          t.ID,
		Name:        t.Name,
		Type:        t.Type,
		Description: t.Description,
		CreatedAt:   t.CreatedAt,
		UpdatedAt:   t.UpdatedAt,
		DeletedAt:   t.DeletedAt,
	}
}
