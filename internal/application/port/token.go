package port

import (
	"time"

	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
)

// TokenProvider generates and validates JWT tokens (used by auth use case).
type TokenProvider interface {
	GenerateTokenPair(userID uuid.UUID, email string, sessionID uuid.UUID) (*domain.TokenPair, error)
	ValidateToken(tokenString string) (*domain.TokenClaims, error)
	RefreshTokenTTL() time.Duration
}
