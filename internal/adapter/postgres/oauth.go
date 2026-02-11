package postgres

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/application/port"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
	"github.com/jmoiron/sqlx"
	"golang.org/x/oauth2"
)

type oauthProviderRepo struct {
	db *sqlx.DB
}

func NewOAuthProviderRepository(db *sqlx.DB) port.OAuthProviderRepository {
	return &oauthProviderRepo{db: db}
}

func (r *oauthProviderRepo) GetByProviderAndProviderUserID(ctx context.Context, provider, providerUserID string) (*domain.OAuthProvider, error) {
	var oauthProvider domain.OAuthProvider
	err := r.db.GetContext(ctx, &oauthProvider,
		`SELECT id, user_id, provider, provider_user_id, provider_email, access_token, refresh_token, token_expires_at, created_at, updated_at FROM tbl_auth_provider WHERE provider = $1 AND provider_user_id = $2 AND deleted_at IS NULL`,
		provider, providerUserID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &oauthProvider, nil
}

func (r *oauthProviderRepo) Create(ctx context.Context, provider *domain.OAuthProvider) error {
	query := `INSERT INTO tbl_auth_provider (id, user_id, provider, provider_user_id, provider_email, access_token, refresh_token, token_expires_at, created_at, updated_at)
		VALUES (:id, :user_id, :provider, :provider_user_id, :provider_email, :access_token, :refresh_token, :token_expires_at, :created_at, :updated_at)`
	_, err := r.db.NamedExecContext(ctx, query, provider)
	return err
}

func (r *oauthProviderRepo) UpdateOrCreate(ctx context.Context, provider *domain.OAuthProvider, providerName, providerUserID string, userID uuid.UUID, token *oauth2.Token) (bool, error) {
	err := r.db.GetContext(ctx, provider,
		`SELECT id, user_id, provider, provider_user_id, provider_email, access_token, refresh_token, token_expires_at, created_at, updated_at FROM tbl_auth_provider WHERE provider = $1 AND provider_user_id = $2 AND deleted_at IS NULL`,
		providerName, providerUserID)
	if err == nil {
		provider.UserID = userID
		provider.AccessToken = token.AccessToken
		if token.RefreshToken != "" {
			provider.RefreshToken = token.RefreshToken
		}
		if !token.Expiry.IsZero() {
			expiresAt := token.Expiry
			provider.TokenExpiresAt = &expiresAt
		}
		provider.UpdatedAt = time.Now()
		_, err = r.db.ExecContext(ctx,
			`UPDATE tbl_auth_provider SET user_id = $1, access_token = $2, refresh_token = $3, token_expires_at = $4, updated_at = $5 WHERE id = $6`,
			provider.UserID, provider.AccessToken, provider.RefreshToken, provider.TokenExpiresAt, provider.UpdatedAt, provider.ID)
		return true, err
	}
	if err != sql.ErrNoRows {
		return false, err
	}
	return false, nil
}
