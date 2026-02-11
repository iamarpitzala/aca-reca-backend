package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
)

type UserRepository interface {
	Create(ctx context.Context, user *domain.User) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.User, error)
	GetByEmail(ctx context.Context, email string) (*domain.User, error)
	EmailExists(ctx context.Context, email string) (bool, error)
	Update(ctx context.Context, user *domain.User) error
}

type SessionRepository interface {
	Create(ctx context.Context, session *domain.Session) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Session, error)
	GetByRefreshToken(ctx context.Context, refreshToken string) (*domain.Session, error)
	Update(ctx context.Context, session *domain.Session) error
	Delete(ctx context.Context, id uuid.UUID) error
	GetUserSessions(ctx context.Context, userID uuid.UUID) ([]domain.Session, error)
	Revoke(ctx context.Context, sessionID uuid.UUID) error
}
