package usecase

import (
	"context"

	"github.com/iamarpitzala/aca-reca-backend/internal/application/port"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain/form"
)

type TaxTypeService struct {
	repo port.TaxTypeRepository
}

func NewTaxTypeService(repo port.TaxTypeRepository) *TaxTypeService {
	return &TaxTypeService{repo: repo}
}

func (s *TaxTypeService) GetTaxTypes(ctx context.Context) ([]form.TaxTypeResponse, error) {
	list, err := s.repo.ListAll(ctx)
	if err != nil {
		return nil, err
	}
	out := make([]form.TaxTypeResponse, 0, len(list))
	for i := range list {
		out = append(out, *list[i].ToTaxTypeResponse())
	}
	return out, nil
}
