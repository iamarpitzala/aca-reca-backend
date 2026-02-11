package port

import (
	"context"

	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
	"golang.org/x/oauth2"
)

type OAuthProviderRepository interface {
	GetByProviderAndProviderUserID(ctx context.Context, provider, providerUserID string) (*domain.OAuthProvider, error)
	Create(ctx context.Context, provider *domain.OAuthProvider) error
	UpdateOrCreate(ctx context.Context, provider *domain.OAuthProvider, providerName, providerUserID string, userID uuid.UUID, token *oauth2.Token) (updated bool, err error)
}
