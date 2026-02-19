package form

import (
	"github.com/google/uuid"
)

type TaxConfigRequest struct {
	ID   *string `json:"id" validate:"omitempty,required"`
	Rate float64 `json:"rate" validate:"required"`
	Name string  `json:"name" validate:"required,min=3,max=255"`
}

type TaxConfig struct {
	ID   uuid.UUID `db:"id"`
	Rate float64   `db:"rate"`
	Name string    `db:"name"`
}

func (t *TaxConfig) ToTaxConfigDB(taxConfig *TaxConfigRequest) {
	if taxConfig.ID != nil && *taxConfig.ID != "" {
		t.ID = uuid.MustParse(*taxConfig.ID)
	}
	t.Rate = taxConfig.Rate
	t.Name = taxConfig.Name
}

type TaxConfigResponse struct {
	ID   uuid.UUID `json:"id"`
	Rate float64   `json:"rate"`
	Name string    `json:"name"`
}

func (t *TaxConfig) ToTaxConfigResponse() *TaxConfigResponse {
	return &TaxConfigResponse{
		ID:   t.ID,
		Rate: t.Rate,
		Name: t.Name,
	}
}
