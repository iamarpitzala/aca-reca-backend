package domain

import (
	"time"

	"github.com/google/uuid"
)

type Quarter struct {
	ID        uuid.UUID
	Name      string
	StartDate time.Time
	EndDate   time.Time
	CreatedAt time.Time
	UpdatedAt *time.Time
	DeletedAt *time.Time
}
