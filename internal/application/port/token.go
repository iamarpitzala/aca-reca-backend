package port

import (
	"time"

	"github.com/iamarpitzala/aca-reca-backend/internal/domain/auth"
)

// TokenProvider generates and validates JWT tokens (used by auth use case).
type TokenProvider interface {
	GenerateTokenPair(userID string, email string, sessionID string) (*auth.TokenPair, error)
	ValidateToken(tokenString string) (*auth.TokenClaims, error)
	RefreshTokenTTL() time.Duration
}
