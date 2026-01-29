package config

import (
	"os"
	"strconv"
	"time"
)

type RedisConfig struct {
	Host       string
	Port       string
	Password   string
	DB         int
	TLSEnabled bool
}

// Config struct
type Config struct {
	Server ServerConfig
	DB     DBConfig
	JWT    JWTConfig
	OAuth  OAuthConfig
}

type ServerConfig struct {
	Port string
	Env  string
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

type JWTConfig struct {
	SecretKey       string
	AccessTokenTTL  time.Duration
	RefreshTokenTTL time.Duration
	Issuer          string
}

type SessionConfig struct {
	TTL time.Duration
}

type OAuthConfig struct {
	RedirectURL  string
	FrontendURL  string // Where to redirect after OAuth success (tokens in hash)
	Providers    map[string]OAuthProviderConfig
}

type OAuthProviderConfig struct {
	ClientID     string
	ClientSecret string
	AuthURL      string
	TokenURL     string
	UserInfoURL  string
	Scopes       []string
}

func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
			Env:  getEnv("ENV", "development"),
		},
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			DBName:   getEnv("DB_NAME", "acareca"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		JWT: JWTConfig{
			SecretKey:       getEnv("JWT_SECRET_KEY", "your-secret-key-change-in-production"),
			AccessTokenTTL:  getEnvAsDuration("JWT_ACCESS_TTL", 15*time.Minute),
			RefreshTokenTTL: getEnvAsDuration("JWT_REFRESH_TTL", 7*24*time.Hour),
			Issuer:          getEnv("JWT_ISSUER", "motocabz-sso"),
		},
		OAuth: OAuthConfig{
			// OAUTH_REDIRECT_URL: base for OAuth callbacks (e.g. http://localhost:8080/api/v1/auth/oauth -> .../google/callback)
			RedirectURL: getEnv("OAUTH_REDIRECT_URL", "http://localhost:8080/api/v1/auth/oauth"),
			// FRONTEND_URL: where to redirect after OAuth success (tokens in hash)
			FrontendURL: getEnv("FRONTEND_URL", "http://localhost:5173"),
			Providers: map[string]OAuthProviderConfig{
				"google": {
					ClientID:     getEnv("GOOGLE_CLIENT_ID", ""),
					ClientSecret: getEnv("GOOGLE_CLIENT_SECRET", ""),
					AuthURL:      "https://accounts.google.com/o/oauth2/v2/auth",
					TokenURL:     "https://oauth2.googleapis.com/token",
					UserInfoURL:  "https://www.googleapis.com/oauth2/v2/userinfo",
					Scopes:       []string{"openid", "profile", "email"},
				},
				"microsoft": {
					ClientID:     getEnv("MICROSOFT_CLIENT_ID", ""),
					ClientSecret: getEnv("MICROSOFT_CLIENT_SECRET", ""),
					AuthURL:      "https://login.microsoftonline.com/common/oauth2/v2.0/authorize",
					TokenURL:     "https://login.microsoftonline.com/common/oauth2/v2.0/token",
					UserInfoURL:  "https://graph.microsoft.com/v1.0/me",
					Scopes:       []string{"openid", "profile", "email"},
				},
			},
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := getEnv(key, "")
	if value, err := time.ParseDuration(valueStr); err == nil {
		return value
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := getEnv(key, "")
	if valueStr == "" {
		return defaultValue
	}
	value, err := strconv.ParseBool(valueStr)
	if err != nil {
		return defaultValue
	}
	return value
}
