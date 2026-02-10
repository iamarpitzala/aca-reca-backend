package middleware

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// normalizeOrigin trims trailing slashes for consistent comparison.
func normalizeOrigin(origin string) string {
	return strings.TrimSuffix(strings.TrimSpace(origin), "/")
}

func CorsMiddleware() gin.HandlerFunc {
	allowedOrigins := []string{
		"https://zenithive.lovable.app",
		"https://preview--zenithive.lovable.app",
		"http://localhost:5173",
		"https://lovable.dev/projects/d4f867cd-6d95-4580-8932-efc09c741d1e",
		"https://acareca.netlify.app",
	}

	if env := os.Getenv("CORS_ORIGINS"); env != "" {
		seen := make(map[string]bool)
		for _, o := range allowedOrigins {
			seen[normalizeOrigin(o)] = true
		}
		for _, o := range strings.Split(env, ",") {
			o = normalizeOrigin(o)
			if o != "" && !seen[o] {
				seen[o] = true
				allowedOrigins = append(allowedOrigins, o)
			}
		}
	}

	isAllowed := func(origin string) bool {
		origin = normalizeOrigin(origin)
		for _, o := range allowedOrigins {
			if origin == normalizeOrigin(o) {
				return true
			}
		}
		return false
	}

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")

		if origin != "" && isAllowed(origin) {
			c.Header("Access-Control-Allow-Origin", origin)
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			c.Header(
				"Access-Control-Allow-Headers",
				"Authorization, Content-Type, Accept, Origin, X-Requested-With",
			)
			c.Header("Access-Control-Max-Age", "86400")
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
