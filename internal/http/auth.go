package http

import (
	"crypto/rand"
	"encoding/hex"
	"net/http"
	"net/url"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/iamarpitzala/aca-reca-backend/internal/application/usecase"
	"github.com/iamarpitzala/aca-reca-backend/internal/domain"
	"github.com/iamarpitzala/aca-reca-backend/internal/service"
	utils "github.com/iamarpitzala/aca-reca-backend/util"
)

// AuthHandler handles auth HTTP endpoints (register, login, refresh, logout, OAuth).
type AuthHandler struct {
	authUC       *usecase.AuthService
	oauthService *service.OAuthService
	frontendURL  string
}

// NewAuthHandler returns a new AuthHandler.
func NewAuthHandler(authUC *usecase.AuthService, oauthService *service.OAuthService, frontendURL string) *AuthHandler {
	return &AuthHandler{
		authUC:       authUC,
		oauthService: oauthService,
		frontendURL:  frontendURL,
	}
}

// refreshTokenRequest is the body for POST /auth/refresh.
type refreshTokenRequest struct {
	RefreshToken string `json:"refreshToken" binding:"required"`
}

// Register handles POST /auth/register.
func (h *AuthHandler) Register(c *gin.Context) {
	var req domain.RegisterRequest
	if err := utils.BindAndValidate(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := h.authUC.Register(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// Login handles POST /auth/login.
func (h *AuthHandler) Login(c *gin.Context) {
	var req domain.LoginRequest
	if err := utils.BindAndValidate(c, &req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	resp, err := h.authUC.Login(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// RefreshToken handles POST /auth/refresh.
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	var req refreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.ErrRefreshTokenRequired})
		return
	}
	resp, err := h.authUC.RefreshToken(c.Request.Context(), req.RefreshToken)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, resp)
}

// Logout handles POST /auth/logout/:sessionId.
func (h *AuthHandler) Logout(c *gin.Context) {
	sessionIDStr := c.Param("sessionId")
	sessionID, err := uuid.Parse(sessionIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.ErrInvalidSessionID})
		return
	}
	if err := h.authUC.Logout(c.Request.Context(), sessionID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": utils.MsgLoggedOutSuccessfully})
}

// InitiateOAuth handles GET /auth/oauth/:provider - redirects to provider's authorization URL.
func (h *AuthHandler) InitiateOAuth(c *gin.Context) {
	provider := c.Param("provider")
	if provider == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.ErrProviderRequired})
		return
	}
	stateBytes := make([]byte, 16)
	if _, err := rand.Read(stateBytes); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": utils.ErrFailedToGenerateState})
		return
	}
	state := hex.EncodeToString(stateBytes)
	authURL, err := h.oauthService.GetAuthURL(provider, state)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	c.Redirect(http.StatusFound, authURL)
}

// OAuthCallback handles GET /auth/oauth/:provider/callback - exchanges code for tokens and redirects to frontend.
func (h *AuthHandler) OAuthCallback(c *gin.Context) {
	provider := c.Param("provider")
	if provider == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": utils.ErrProviderRequired})
		return
	}
	code := c.Query(utils.OAuthParamCode)
	if code == "" {
		c.Redirect(http.StatusTemporaryRedirect, h.frontendURL+"?"+utils.OAuthParamError+"="+utils.OAuthErrorMissingCode)
		return
	}
	ctx := c.Request.Context()
	token, err := h.oauthService.ExchangeCode(ctx, provider, code)
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, h.frontendURL+"?"+utils.OAuthParamError+"="+utils.OAuthErrorExchangeFailed)
		return
	}
	userInfo, err := h.oauthService.GetUserInfo(ctx, provider, token)
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, h.frontendURL+"?"+utils.OAuthParamError+"="+utils.OAuthErrorUserInfoFailed)
		return
	}
	user, err := h.oauthService.FindUserByProvider(ctx, provider, userInfo.ID)
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, h.frontendURL+"?"+utils.OAuthParamError+"="+utils.OAuthErrorLookupFailed)
		return
	}
	if user == nil {
		user, err = h.oauthService.CreateUserFromOAuth(ctx, userInfo)
		if err != nil {
			c.Redirect(http.StatusTemporaryRedirect, h.frontendURL+"?"+utils.OAuthParamError+"="+utils.OAuthErrorCreateUserFailed)
			return
		}
		if err := h.authUC.EnsureDefaultClinicForUser(ctx, user.ID, user.FirstName); err != nil {
			// non-fatal: user exists
		}
	} else {
		if err := h.oauthService.LinkProvider(ctx, user.ID, provider, userInfo.ID, userInfo.Email, token); err != nil {
			// optional: link for future logins
		}
	}
	authResp, err := h.authUC.OAuthLogin(ctx, user.ID)
	if err != nil {
		c.Redirect(http.StatusTemporaryRedirect, h.frontendURL+"?"+utils.OAuthParamError+"="+utils.OAuthErrorLoginFailed)
		return
	}
	q := url.Values{}
	q.Set(utils.OAuthParamAccessToken, authResp.AccessToken)
	q.Set(utils.OAuthParamRefreshToken, authResp.RefreshToken)
	q.Set(utils.OAuthParamTokenType, authResp.TokenType)
	redirectURL := h.frontendURL + "?" + q.Encode()
	c.Redirect(http.StatusFound, redirectURL)
}

// GetAuthUserID extracts the authenticated user ID from the request context (set by auth middleware).
// If the user is not authenticated or the context value is invalid, it responds with 401 and returns (uuid.Nil, false).
func GetAuthUserID(c *gin.Context) (uuid.UUID, bool) {
	v, ok := c.Get("user_id")
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": utils.ErrUnauthorized})
		return uuid.Nil, false
	}
	userID, ok := v.(uuid.UUID)
	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": utils.ErrInvalidUserContext})
		return uuid.Nil, false
	}
	return userID, true
}

// RequireClinicAccess verifies the authenticated user has access to the given clinic.
// If the user is not authenticated or does not have access, it responds with 401/403 and returns false.
// Callers should return immediately when RequireClinicAccess returns false (response already sent).
func RequireClinicAccess(c *gin.Context, userClinicUC *usecase.UserClinicService, clinicID uuid.UUID) bool {
	userID, ok := GetAuthUserID(c)
	if !ok {
		return false
	}
	hasAccess, err := userClinicUC.UserHasAccessToClinic(c.Request.Context(), userID, clinicID)
	if err != nil || !hasAccess {
		c.JSON(http.StatusForbidden, gin.H{"error": utils.ErrAccessDenied})
		return false
	}
	return true
}
