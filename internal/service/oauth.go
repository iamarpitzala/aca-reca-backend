package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/config"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
	"github.com/iamarpitzala/aca-reca-backend/internal/repository"
	"github.com/jmoiron/sqlx"
	"golang.org/x/oauth2"
)

type OAuthService struct {
	config    config.OAuthConfig
	db        *sqlx.DB
	providers map[string]*oauth2.Config
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

func (os *OAuthService) GetUserInfo(ctx context.Context, provider string, token *oauth2.Token) (*domain.OAuthUserInfo, error) {
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

	var userInfo domain.OAuthUserInfo
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
		userInfo = domain.OAuthUserInfo{
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
		userInfo = domain.OAuthUserInfo{
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

func (os *OAuthService) LinkProvider(ctx context.Context, userID uuid.UUID, provider string, providerUserID string, providerEmail string, token *oauth2.Token) error {
	var oauthProvider domain.OAuthProvider

	// Check if provider link already exists and update it
	exists, err := repository.UpdateOrCreateOAuthProvider(ctx, os.db, &oauthProvider, provider, providerUserID, userID, token)
	if err != nil {
		return fmt.Errorf("failed to update OAuth provider: %w", err)
	}

	// If provider already exists and was updated, update email if needed
	if exists {
		if providerEmail != "" && oauthProvider.ProviderEmail != providerEmail {
			oauthProvider.ProviderEmail = providerEmail
			_, err = os.db.ExecContext(ctx,
				"UPDATE tbl_auth_provider SET provider_email = $1, updated_at = $2 WHERE id = $3",
				providerEmail, time.Now(), oauthProvider.ID)
			if err != nil {
				return fmt.Errorf("failed to update provider email: %w", err)
			}
		}
		return nil
	}

	// Create new link
	expiresAt := token.Expiry
	oauthProvider = domain.OAuthProvider{
		ID:             uuid.New(),
		UserID:         userID,
		Provider:       provider,
		ProviderUserID: providerUserID,
		ProviderEmail:  providerEmail,
		AccessToken:    token.AccessToken,
		RefreshToken:   token.RefreshToken,
		TokenExpiresAt: &expiresAt,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}

	// Create new OAuth provider link
	err = repository.CreateOAuthProvider(ctx, os.db, &oauthProvider)
	if err != nil {
		return fmt.Errorf("failed to create OAuth provider: %w", err)
	}
	return nil
}

func (os *OAuthService) FindUserByProvider(ctx context.Context, provider string, providerUserID string) (*domain.User, error) {
	oauthProvider, err := repository.GetOAuthProvider(ctx, os.db, provider, providerUserID)
	if err != nil {
		return nil, err
	}

	user, err := repository.GetUserByID(ctx, os.db, oauthProvider.UserID)
	if err != nil {
		return nil, err
	}

	return user, nil
}

func (os *OAuthService) CreateUserFromOAuth(ctx context.Context, userInfo *domain.OAuthUserInfo) (*domain.User, error) {
	user := domain.User{
		ID:              uuid.New(),
		Email:           userInfo.Email,
		FirstName:       userInfo.FirstName,
		LastName:        userInfo.LastName,
		AvatarURL:       userInfo.AvatarURL,
		IsEmailVerified: userInfo.EmailVerified,
		CreatedAt:       time.Now(),
		UpdatedAt:       time.Now(),
	}

	err := repository.CreateUser(ctx, os.db, &user)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
