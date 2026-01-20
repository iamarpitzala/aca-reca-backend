package repository

import (
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
	"github.com/jmoiron/sqlx"
	"golang.org/x/oauth2"
)

// Check if provider link already exists
func UpdateOrCreateOAuthProvider(ctx context.Context, db *sqlx.DB, oauthProvider *domain.OAuthProvider, provider string, providerUserID string, userID uuid.UUID, token *oauth2.Token) error {
	err := db.GetContext(ctx, oauthProvider,
		"SELECT id, user_id, provider, provider_user_id, provider_email, access_token, refresh_token, token_expires_at, created_at, updated_at FROM tbl_oauth_provider WHERE provider = $1 AND provider_user_id = $2 AND deleted_at IS NULL",
		provider, providerUserID)

	if err == nil {
		// Update existing link
		oauthProvider.UserID = userID
		oauthProvider.AccessToken = token.AccessToken
		if token.RefreshToken != "" {
			oauthProvider.RefreshToken = token.RefreshToken
		}
		if !token.Expiry.IsZero() {
			expiresAt := token.Expiry
			oauthProvider.TokenExpiresAt = &expiresAt
		}
		oauthProvider.UpdatedAt = time.Now()

		_, err = db.ExecContext(ctx,
			"UPDATE tbl_oauth_provider SET user_id = $1, access_token = $2, refresh_token = $3, token_expires_at = $4, updated_at = $5 WHERE id = $6",
			oauthProvider.UserID, oauthProvider.AccessToken, oauthProvider.RefreshToken, oauthProvider.TokenExpiresAt, oauthProvider.UpdatedAt, oauthProvider.ID)
		return err
	}

	if err != sql.ErrNoRows {
		return err
	}
	return err
}

func CreateOAuthProvider(ctx context.Context, db *sqlx.DB, oauthProvider *domain.OAuthProvider) error {
	query := `INSERT INTO tbl_oauth_provider (id, user_id, provider, provider_user_id, access_token, refresh_token, token_expires_at, created_at, updated_at)
		VALUES (:id, :user_id, :provider, :provider_user_id, :access_token, :refresh_token, :token_expires_at, :created_at, :updated_at)`

	_, err := db.NamedExecContext(ctx, query, oauthProvider)
	return err
}

func GetOAuthProvider(ctx context.Context, db *sqlx.DB, provider string, providerUserId string) (*domain.OAuthProvider, error) {
	var providerProvider domain.OAuthProvider
	err := db.GetContext(ctx, &providerProvider,
		"SELECT id, user_id, provider, provider_user_id, provider_email, access_token, refresh_token, token_expires_at, created_at, updated_at FROM tbl_oauth_provider WHERE provider = $1 AND provider_user_id = $2 AND deleted_at IS NULL",
		provider, providerUserId)
	if err != nil {
		return nil, err
	}
	return &providerProvider, nil
}
