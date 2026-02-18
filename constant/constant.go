package constants

import "github.com/iamarpitzala/aca-reca-backend/util"

// Re-export from util for backward compatibility. Prefer using util directly.
const (
	INCLUSIVE = util.GSTTypeInclusive
	EXCLUSIVE = util.GSTTypeExclusive
	MANUAL    = util.GSTTypeManual

	PERCENTAGE = util.ShareTypePercentage
	FIXED      = util.ShareTypeFixed

	PAID_BY_CLINIC = "PAID_BY_CLINIC"
	PAID_BY_OWNER  = "PAID_BY_OWNER"
)
