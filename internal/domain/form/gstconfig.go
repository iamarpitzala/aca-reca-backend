package form

import "github.com/google/uuid"

type GstType string

const (
	GstTypeInclusive GstType = "INCLUSIVE"
	GstTypeExclusive GstType = "EXCLUSIVE"
	GstTypeManual    GstType = "MANUAL"
)

type GstRate float64

const (
	GstRate0  GstRate = 0
	GstRate10 GstRate = 10.0
)

func (g GstType) String() string {
	return string(g)
}

type GstConfigRequest struct {
	ID      *string  `json:"id" validate:"omitempty,required"`
	Enabled bool     `json:"enabled" validate:"required"`
	Rate    GstRate  `json:"rate" validate:"required"`
	Type    GstType  `json:"type" validate:"required,oneof=INCLUSIVE EXCLUSIVE MANUAL"`
	Amount  *float64 `json:"amount" validate:"omitempty,required"`
}

func (g GstConfigRequest) Validate() error {
	if g.Type == "" {
		g.Type = GstTypeExclusive
	}
	return nil
}

type GstConfig struct {
	ID      uuid.UUID `db:"id"`
	Enabled bool      `db:"enabled"`
	Rate    GstRate   `db:"rate"`
	Type    GstType   `db:"type"`
	Amount  *float64  `db:"amount"`
}

func (g *GstConfig) ToGstDB(gstConfig *GstConfigRequest) {
	g.ID = uuid.MustParse(*gstConfig.ID)
	g.Enabled = gstConfig.Enabled
	g.Rate = gstConfig.Rate
	g.Type = gstConfig.Type
	g.Amount = gstConfig.Amount
}

type GstResponse struct {
	ID      string   `json:"id"`
	Enabled bool     `json:"enabled"`
	Rate    GstRate  `json:"rate"`
	Type    GstType  `json:"type"`
	Amount  *float64 `json:"amount"`
}

func (g *GstConfig) ToGstResponse() *GstResponse {
	return &GstResponse{
		ID:      g.ID.String(),
		Enabled: g.Enabled,
		Rate:    g.Rate,
		Type:    g.Type,
		Amount:  g.Amount,
	}
}
