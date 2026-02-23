package port

import (
	"context"

	"github.com/iamarpitzala/aca-reca-backend/internal/domain/auth"
	"golang.org/x/oauth2"
)

type OAuthProviderRepository interface {
	GetByProviderAndProviderUserID(ctx context.Context, provider, providerUserID string) (*auth.OAuthProvider, error)
	Create(ctx context.Context, provider *auth.OAuthProvider) error
	Update(ctx context.Context, provider *auth.OAuthProvider) error
	UpdateOrCreate(ctx context.Context, provider *auth.OAuthProvider, providerName, providerUserID string, userID string, token *oauth2.Token) (existed bool, err error)
}
