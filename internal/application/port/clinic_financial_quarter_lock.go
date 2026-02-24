package port

import (
	"context"
)

type ClinicFinancialQuarterLockRepository interface {
	CreateLock(ctx context.Context, clinicFinancialYearID, financialQuarterID int) error
}
