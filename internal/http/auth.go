package http

import (
	"net/http"
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
}

func NewAuthHandler(authService *service.AuthService, oauthService *service.OAuthService) *AuthHandler {
	return &AuthHandler{
		authService:  authService,
		oauthService: oauthService,
	}
}

// Register handles user registration
// POST /api/v1/auth/register
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
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req struct {
		RefreshToken string `json:"refresh_token" binding:"required"`
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
func (h *AuthHandler) Logout(c *gin.Context) {
	sessionID := c.Param("session_id")
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
func (h *AuthHandler) InitiateOAuth(c *gin.Context) {
	provider := c.Param("provider")

	// Generate state token for CSRF protection
	state := uuid.New().String()

	// Store state in session/cookie (in production, use secure cookie)
	c.SetCookie("oauth_state", state, 600, "/", "", false, true)

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

		c.JSON(http.StatusOK, response)
		return
	}

	// User exists, link provider if not already linked
	if err := h.oauthService.LinkProvider(c.Request.Context(), existingUser.ID, provider, userInfo.ID, userInfo.Email, token); err != nil {
		// Provider might already be linked, log but continue
		// In production, you might want to log this error
	}

	// Generate tokens and create session
	response, err := h.authService.OAuthLogin(c.Request.Context(), existingUser.ID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error(), "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, response)
}
