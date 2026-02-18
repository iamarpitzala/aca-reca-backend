package service

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"time"

	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/config"
	"github.com/iamarpitzala/aca-reca-backend/internal/application/port"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
	"golang.org/x/oauth2"
)

type OAuthService struct {
	config       config.OAuthConfig
	providerRepo port.OAuthProviderRepository
	userRepo     port.UserRepository
	providers    map[string]*oauth2.Config
}

func NewOAuthService(cfg config.OAuthConfig, providerRepo port.OAuthProviderRepository, userRepo port.UserRepository) *OAuthService {
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
		config:       cfg,
		providerRepo: providerRepo,
		userRepo:     userRepo,
		providers:    providers,
	}
}

func (os *OAuthService) GetAuthURL(provider string, state string) (string, error) {
	oauthConfig, ok := os.providers[provider]
	if !ok {
		return "", fmt.Errorf("oauth provider %s not configured", provider)
	}
	return oauthConfig.AuthCodeURL(state, oauth2.AccessTypeOffline), nil
}

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
			EmailVerified: true,
		}
	default:
		return nil, fmt.Errorf("unsupported provider: %s", provider)
	}
	return &userInfo, nil
}

func (os *OAuthService) LinkProvider(ctx context.Context, userID uuid.UUID, provider string, providerUserID string, providerEmail string, token *oauth2.Token) error {
	var oauthProvider domain.OAuthProvider
	existed, err := os.providerRepo.UpdateOrCreate(ctx, &oauthProvider, provider, providerUserID, userID, token)
	if err != nil {
		return fmt.Errorf("failed to update OAuth provider: %w", err)
	}
	if existed {
		if providerEmail != "" && oauthProvider.ProviderEmail != providerEmail {
			oauthProvider.ProviderEmail = providerEmail
			oauthProvider.UpdatedAt = time.Now()
			if err := os.providerRepo.Update(ctx, &oauthProvider); err != nil {
				return fmt.Errorf("failed to update provider email: %w", err)
			}
		}
		return nil
	}
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
	if err := os.providerRepo.Create(ctx, &oauthProvider); err != nil {
		return fmt.Errorf("failed to create OAuth provider: %w", err)
	}
	return nil
}

func (os *OAuthService) FindUserByProvider(ctx context.Context, provider string, providerUserID string) (*domain.User, error) {
	oauthProvider, err := os.providerRepo.GetByProviderAndProviderUserID(ctx, provider, providerUserID)
	if err != nil {
		return nil, err
	}
	if oauthProvider == nil {
		return nil, nil
	}
	return os.userRepo.GetByID(ctx, oauthProvider.UserID)
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
	if err := os.userRepo.Create(ctx, &user); err != nil {
		return nil, err
	}
	return &user, nil
}
