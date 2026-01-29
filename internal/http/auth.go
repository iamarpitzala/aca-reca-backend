package http

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
	"github.com/iamarpitzala/aca-reca-backend/internal/service"
	utils "github.com/iamarpitzala/aca-reca-backend/util"
)

// contains checks if a string contains a substring (case-insensitive)
func contains(s, substr string) bool {
	return strings.Contains(strings.ToLower(s), strings.ToLower(substr))
}

type AuthHandler struct {
	authService  *service.AuthService
	oauthService *service.OAuthService
	frontendURL  string
}

func NewAuthHandler(authService *service.AuthService, oauthService *service.OAuthService, frontendURL string) *AuthHandler {
	return &AuthHandler{
		authService:  authService,
		oauthService: oauthService,
		frontendURL:  frontendURL,
	}
}

// Register handles user registration
// POST /api/v1/auth/register
// @Summary Register a new user
// @Description Register a new user with email and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param register_request body domain.RegisterRequest true "Register request"
// @Success 201 {object} domain.AuthResponse
// @Failure 400 {object} domain.H
// @Failure 500 {object} domain.H
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req domain.RegisterRequest
	if err := utils.BindAndValidate(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.authService.Register(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, response)
}

// Login handles user login
// POST /api/v1/auth/login
// @Summary Login a user
// @Description Login a user with email and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param login_request body domain.LoginRequest true "Login request"
// @Success 200 {object} domain.AuthResponse
// @Failure 400 {object} domain.H
// @Failure 500 {object} domain.H
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req domain.LoginRequest

	if err := utils.BindAndValidate(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	response, err := h.authService.Login(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}

// RefreshToken handles token refresh
// POST /api/v1/auth/refresh
// @Summary Refresh a token
// @Description Refresh a token with a refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Param refreshToken body string true "Refresh token"
// @Success 200 {object} domain.AuthResponse
// @Failure 400 {object} domain.H
// @Failure 500 {object} domain.H
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refreshToken" validate:"required"`
	}

	if err := utils.BindAndValidate(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	tokenPair, err := h.authService.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, tokenPair)
}

// Logout handles user logout
// POST /api/v1/auth/logout
// @Summary Logout a user
// @Description Logout a user with a session ID
// @Tags Auth
// @Accept json
// @Produce json
// @Param sessionId path string true "Session ID"
// @Success 200 {object} domain.H
// @Failure 400 {object} domain.H
// @Failure 500 {object} domain.H
// @Router /auth/logout/{sessionId} [post]
func (h *AuthHandler) Logout(c *gin.Context) {
	sessionID := c.Param("sessionId")
	sessionUUID, err := uuid.Parse(sessionID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid session ID"})
		return
	}

	if err := h.authService.Logout(c.Request.Context(), sessionUUID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "logged out successfully"})
}

// InitiateOAuth initiates OAuth flow
// GET /api/v1/auth/oauth/:provider
// @Summary Initiate OAuth flow
// @Description Initiate OAuth flow with a provider
// @Tags Auth
// @Accept json
// @Produce json
// @Param provider query string true "Provider"
// @Param state query string true "State"
// @Success 307 {string} string "Redirect to OAuth provider"
// @Failure 400 {object} domain.H
// @Failure 500 {object} domain.H
// @Router /auth/oauth [get]
func (h *AuthHandler) InitiateOAuth(c *gin.Context) {
	provider := c.Param("provider")

	// Generate state token for CSRF protection
	state := uuid.New().String()

	// Store state in cookie; SameSite=Lax so cookie is sent when Google redirects back
	secure := c.Request.URL.Scheme == "https" || c.Request.TLS != nil
	cookie := &http.Cookie{
		Name:     "oauth_state",
		Value:    state,
		Path:     "/",
		MaxAge:   600,
		HttpOnly: true,
		Secure:   secure,
		SameSite: http.SameSiteLaxMode,
	}
	http.SetCookie(c.Writer, cookie)

	authURL, err := h.oauthService.GetAuthURL(provider, state)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
			"hint":  "Ensure OAuth provider is configured with CLIENT_ID and CLIENT_SECRET",
		})
		return
	}

	c.Redirect(http.StatusTemporaryRedirect, authURL)
}

// OAuthCallback handles OAuth callback
// GET /api/v1/auth/oauth/:provider/callback
// @Summary OAuth callback
// @Description OAuth callback with a provider
// @Tags Auth
// @Accept json
// @Produce json
// @Param provider path string true "Provider"
// @Success 200 {object} domain.AuthResponse
// @Failure 400 {object} domain.H
// @Failure 500 {object} domain.H
// @Router /auth/oauth/{provider}/callback [get]
func (h *AuthHandler) OAuthCallback(c *gin.Context) {
	provider := c.Param("provider")

	// Verify state
	stateCookie, err := c.Cookie("oauth_state")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing state cookie"})
		return
	}

	state := c.Query("state")
	if state != stateCookie {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid state parameter"})
		return
	}

	code := c.Query("code")
	if code == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "missing authorization code"})
		return
	}

	// Exchange code for token
	token, err := h.oauthService.ExchangeCode(c.Request.Context(), provider, code)
	if err != nil {
		// Provide helpful error message for redirect_uri_mismatch
		errorMsg := err.Error()
		if contains(errorMsg, "redirect_uri_mismatch") {
			redirectURI, _ := h.oauthService.GetRedirectURI(provider)
			c.JSON(http.StatusBadRequest, gin.H{
				"error":                   "OAuth redirect URI mismatch",
				"details":                 errorMsg,
				"configured_redirect_uri": redirectURI,
				"hint": "Ensure this exact redirect URI is added in your OAuth provider's console. " +
					"For Google: https://console.cloud.google.com/apis/credentials",
			})
			return
		}
		c.JSON(http.StatusBadRequest, gin.H{"error": errorMsg})
		return
	}

	// Get user info from provider
	userInfo, err := h.oauthService.GetUserInfo(c.Request.Context(), provider, token)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to get user info from provider", "details": err.Error()})
		return
	}

	// Find or create user
	existingUser, err := h.oauthService.FindUserByProvider(c.Request.Context(), provider, userInfo.ID)
	if err != nil {
		// User doesn't exist, create new user
		newUser, err := h.oauthService.CreateUserFromOAuth(c.Request.Context(), userInfo)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user", "details": err.Error()})
			return
		}

		// Link OAuth provider
		if err := h.oauthService.LinkProvider(c.Request.Context(), newUser.ID, provider, userInfo.ID, userInfo.Email, token); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to link provider", "details": err.Error()})
			return
		}

		// Generate tokens and create session
		response, err := h.authService.OAuthLogin(c.Request.Context(), newUser.ID)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "details": err.Error()})
			return
		}

		h.redirectToFrontendWithTokens(c, response)
		return
	}

	// User exists, link provider if not already linked
	if err := h.oauthService.LinkProvider(c.Request.Context(), existingUser.ID, provider, userInfo.ID, userInfo.Email, token); err != nil {
		// Provider might already be linked, log but continue
		// In production, you might want to log this error

		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to link provider", "details": err.Error()})
		return
	}

	// Generate tokens and create session
	response, err := h.authService.OAuthLogin(c.Request.Context(), existingUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "details": err.Error()})
		return
	}

	h.redirectToFrontendWithTokens(c, response)
}

// redirectToFrontendWithTokens redirects the browser to the frontend callback with tokens in the URL hash (not logged).
func (h *AuthHandler) redirectToFrontendWithTokens(c *gin.Context, response *domain.AuthResponse) {
	if h.frontendURL == "" {
		c.JSON(http.StatusOK, response)
		return
	}
	redirectURL, err := url.Parse(strings.TrimSuffix(h.frontendURL, "/") + "/auth/callback")
	if err != nil {
		c.JSON(http.StatusOK, response)
		return
	}
	q := redirectURL.Query()
	q.Set("access_token", response.AccessToken)
	q.Set("refresh_token", response.RefreshToken)
	if response.User != nil {
		userJSON, _ := json.Marshal(response.User)
		q.Set("user", base64.URLEncoding.EncodeToString(userJSON))
	}
	redirectURL.RawQuery = q.Encode()
	redirectURL.Fragment = redirectURL.RawQuery
	redirectURL.RawQuery = ""
	c.Redirect(http.StatusTemporaryRedirect, redirectURL.String())
}
