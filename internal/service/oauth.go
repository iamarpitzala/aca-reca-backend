package service

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/config"
	"github.com/iamarpitzala/aca-reca-backend/internal/model"
	"github.com/jmoiron/sqlx"
	"golang.org/x/oauth2"
)

type OAuthService struct {
	config    config.OAuthConfig
	db        *sqlx.DB
	providers map[string]*oauth2.Config
}

type OAuthUserInfo struct {
	ID            string
	Email         string
	FirstName     string
	LastName      string
	AvatarURL     string
	EmailVerified bool
}

func NewOAuthService(cfg config.OAuthConfig, db *sqlx.DB) *OAuthService {
	providers := make(map[string]*oauth2.Config)

	for name, providerCfg := range cfg.Providers {
		if providerCfg.ClientID != "" && providerCfg.ClientSecret != "" {
			providers[name] = &oauth2.Config{
				ClientID:     providerCfg.ClientID,
				ClientSecret: providerCfg.ClientSecret,
				RedirectURL:  fmt.Sprintf("%s/%s/callback", cfg.RedirectURL, name),
				Scopes:       providerCfg.Scopes,
				Endpoint: oauth2.Endpoint{
					AuthURL:  providerCfg.AuthURL,
					TokenURL: providerCfg.TokenURL,
				},
			}
		}
	}

	return &OAuthService{
		config:    cfg,
		db:        db,
		providers: providers,
	}
}

func (os *OAuthService) GetAuthURL(provider string, state string) (string, error) {
	oauthConfig, ok := os.providers[provider]
	if !ok {
		return "", fmt.Errorf("oauth provider %s not configured", provider)
	}

	return oauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline), nil
}

// GetRedirectURI returns the configured redirect URI for a provider
func (os *OAuthService) GetRedirectURI(provider string) (string, error) {
	oauthConfig, ok := os.providers[provider]
	if !ok {
		return "", fmt.Errorf("oauth provider %s not configured", provider)
	}

	return oauthConfig.RedirectURL, nil
}

func (os *OAuthService) ExchangeCode(ctx context.Context, provider string, code string) (*oauth2.Token, error) {
	oauthConfig, ok := os.providers[provider]
	if !ok {
		return nil, fmt.Errorf("oauth provider %s not configured", provider)
	}

	token, err := oauthConfig.Exchange(ctx, code)
	if err != nil {
		return nil, fmt.Errorf("failed to exchange code for token: %w", err)
	}

	return token, nil
}

func (os *OAuthService) GetUserInfo(ctx context.Context, provider string, token *oauth2.Token) (*OAuthUserInfo, error) {
	providerCfg, ok := os.config.Providers[provider]
	if !ok {
		return nil, fmt.Errorf("oauth provider %s not configured", provider)
	}

	client := oauth2.NewClient(ctx, oauth2.StaticTokenSource(token))
	resp, err := client.Get(providerCfg.UserInfoURL)
	if err != nil {
		return nil, fmt.Errorf("failed to get user info: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read user info response: %w", err)
	}

	var userInfo OAuthUserInfo
	switch provider {
	case "google":
		var googleUser struct {
			ID            string `json:"id"`
			Email         string `json:"email"`
			VerifiedEmail bool   `json:"verified_email"`
			Name          string `json:"name"`
			GivenName     string `json:"given_name"`
			FamilyName    string `json:"family_name"`
			Picture       string `json:"picture"`
		}
		if err := json.Unmarshal(body, &googleUser); err != nil {
			return nil, fmt.Errorf("failed to parse Google user info: %w", err)
		}
		userInfo = OAuthUserInfo{
			ID:            googleUser.ID,
			Email:         googleUser.Email,
			FirstName:     googleUser.GivenName,
			LastName:      googleUser.FamilyName,
			AvatarURL:     googleUser.Picture,
			EmailVerified: googleUser.VerifiedEmail,
		}
	case "microsoft":
		var msUser struct {
			ID                string `json:"id"`
			Mail              string `json:"mail"`
			UserPrincipalName string `json:"userPrincipalName"`
			GivenName         string `json:"givenName"`
			Surname           string `json:"surname"`
		}
		if err := json.Unmarshal(body, &msUser); err != nil {
			return nil, fmt.Errorf("failed to parse Microsoft user info: %w", err)
		}
		email := msUser.Mail
		if email == "" {
			email = msUser.UserPrincipalName
		}
		userInfo = OAuthUserInfo{
			ID:            msUser.ID,
			Email:         email,
			FirstName:     msUser.GivenName,
			LastName:      msUser.Surname,
			EmailVerified: true, // Microsoft accounts are typically verified
		}
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}

	return &userInfo, nil
}

func (os *OAuthService) LinkProvider(ctx context.Context, userID uuid.UUID, provider string, providerUserID string, token *oauth2.Token) error {
	var oauthProvider model.OAuthProvider

	// Check if provider link already exists
	err := os.db.GetContext(ctx, &oauthProvider,
		"SELECT id, user_id, provider, provider_user_id, provider_email, access_token, refresh_token, token_expires_at, created_at, updated_at FROM oauth_providers WHERE provider = $1 AND provider_user_id = $2 AND deleted_at IS NULL",
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

		_, err = os.db.ExecContext(ctx,
			"UPDATE oauth_providers SET user_id = $1, access_token = $2, refresh_token = $3, token_expires_at = $4, updated_at = $5 WHERE id = $6",
			oauthProvider.UserID, oauthProvider.AccessToken, oauthProvider.RefreshToken, oauthProvider.TokenExpiresAt, oauthProvider.UpdatedAt, oauthProvider.ID)
		return err
	}

	if err != sql.ErrNoRows {
		return err
	}

	// Create new link
	expiresAt := token.Expiry
	oauthProvider = model.OAuthProvider{
		ID:             uuid.New(),
		UserID:         userID,
		Provider:       provider,
		ProviderUserID: providerUserID,
		AccessToken:    token.AccessToken,
		RefreshToken:   token.RefreshToken,
		TokenExpiresAt: &expiresAt,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	query := `INSERT INTO oauth_providers (id, user_id, provider, provider_user_id, access_token, refresh_token, token_expires_at, created_at, updated_at)
		VALUES (:id, :user_id, :provider, :provider_user_id, :access_token, :refresh_token, :token_expires_at, :created_at, :updated_at)`

	_, err = os.db.NamedExecContext(ctx, query, oauthProvider)
	return err
}

func (os *OAuthService) FindUserByProvider(ctx context.Context, provider string, providerUserID string) (*model.User, error) {
	var oauthProvider model.OAuthProvider
	err := os.db.GetContext(ctx, &oauthProvider,
		"SELECT id, user_id, provider, provider_user_id, provider_email, access_token, refresh_token, token_expires_at, created_at, updated_at FROM oauth_providers WHERE provider = $1 AND provider_user_id = $2 AND deleted_at IS NULL",
		provider, providerUserID)
	if err != nil {
		return nil, err
	}

	var user model.User
	err = os.db.GetContext(ctx, &user,
		"SELECT id, email, password, first_name, last_name, phone, avatar_url, is_active, is_email_verified, created_at, updated_at FROM users WHERE id = $1 AND deleted_at IS NULL",
		oauthProvider.UserID)
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (os *OAuthService) CreateUserFromOAuth(ctx context.Context, userInfo *OAuthUserInfo) (*model.User, error) {
	user := model.User{
		ID:              uuid.New(),
		Email:           userInfo.Email,
		FirstName:       userInfo.FirstName,
		LastName:        userInfo.LastName,
		AvatarURL:       userInfo.AvatarURL,
		IsActive:        true,
		IsEmailVerified: userInfo.EmailVerified,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	query := `INSERT INTO users (id, email, first_name, last_name, avatar_url, is_active, is_email_verified, created_at, updated_at)
		VALUES (:id, :email, :first_name, :last_name, :avatar_url, :is_active, :is_email_verified, :created_at, :updated_at)`

	_, err := os.db.NamedExecContext(ctx, query, user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
