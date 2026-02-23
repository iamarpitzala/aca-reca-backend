package port

import (
	"context"

	"github.com/iamarpitzala/aca-reca-backend/internal/domain/auth"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain/user"
)

type UserRepository interface {
	Create(ctx context.Context, user *user.User) error
	GetByID(ctx context.Context, id string) (*user.User, error)
	GetByEmail(ctx context.Context, email string) (*user.User, error)
	EmailExists(ctx context.Context, email string) (bool, error)
	Update(ctx context.Context, user *user.User) error
}

type SessionRepository interface {
	Create(ctx context.Context, session *auth.Session) error
	GetByID(ctx context.Context, id string) (*auth.Session, error)
	GetByRefreshToken(ctx context.Context, refreshToken string) (*auth.Session, error)
	Update(ctx context.Context, session *auth.Session) error
	Delete(ctx context.Context, id string) error
	GetUserSessions(ctx context.Context, userID string) ([]auth.Session, error)
	Revoke(ctx context.Context, sessionID string) error
}
