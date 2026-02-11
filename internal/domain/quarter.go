package domain

import (
	"time"

	"github.com/google/uuid"
)

// Quarter represents a stored quarter (deprecated - use CalculatedQuarter from utils instead)
// Kept for backward compatibility with existing database records
type Quarter struct {
	ID        uuid.UUID
	Name      string
	StartDate time.Time
	EndDate   time.Time
	CreatedAt time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}
