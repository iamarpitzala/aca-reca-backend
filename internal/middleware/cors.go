package middleware

import (
	"net/url"
	"os"
	"regexp"
	"strings"

	"github.com/gin-gonic/gin"
)

const defaultAllowHeaders = "Authorization, Content-Type, Accept, Accept-Language, Origin, X-Requested-With, X-Request-ID"

// normalizeOrigin trims trailing slashes for consistent comparison.
func normalizeOrigin(origin string) string {
	return strings.TrimSuffix(strings.TrimSpace(origin), "/")
}

var safeHeaderRe = regexp.MustCompile(`^[a-zA-Z0-9\-_]+$`)

func CorsMiddleware() gin.HandlerFunc {
	allowedOrigins := []string{
		"https://zenithive.lovable.app",
		"https://preview--zenithive.lovable.app",
		"http://localhost:5173",
		"http://127.0.0.1:5173",
		"https://lovable.dev/projects/d4f867cd-6d95-4580-8932-efc09c741d1e",
		"https://acareca.netlify.app",
		"https://lovableproject.com",
		"https://d4f867cd-6d95-4580-8932-efc09c741d1e.lovableproject.com",
		"https://id-preview--d4f867cd-6d95-4580-8932-efc09c741d1e.lovable.app",
	}

	allowAll := false
	if env := os.Getenv("CORS_ORIGINS"); env != "" {
		for _, o := range strings.Split(env, ",") {
			o = strings.TrimSpace(o)
			if o == "*" {
				allowAll = true
				break
			}
			normalized := normalizeOrigin(o)
			if normalized == "" {
				continue
			}
			found := false
			for _, existing := range allowedOrigins {
				if normalizeOrigin(existing) == normalized {
					found = true
					break
				}
			}
			if !found {
				allowedOrigins = append(allowedOrigins, o)
			}
		}
	}

	isAllowed := func(origin string) bool {
		if allowAll {
			return true
		}
		origin = normalizeOrigin(origin)
		if origin == "" {
			return false
		}
		for _, o := range allowedOrigins {
			o = strings.TrimSpace(o)
			if strings.HasPrefix(o, "*.") {
				// Pattern: *.domain.com matches https://anything.domain.com
				pattern := strings.TrimPrefix(o, "*")
				u, err := url.Parse(origin)
				if err != nil {
					continue
				}
				host := strings.ToLower(u.Hostname())
				if host == strings.TrimPrefix(strings.ToLower(pattern), ".") ||
					strings.HasSuffix(host, strings.ToLower(pattern)) {
					return true
				}
			}
			if origin == normalizeOrigin(o) {
				return true
			}
		}
		return false
	}

	return func(c *gin.Context) {
		origin := c.Request.Header.Get("Origin")
		allowed := (allowAll || (origin != "" && isAllowed(origin)))

		if allowed {
			if allowAll {
				c.Header("Access-Control-Allow-Origin", "*")
			} else {
				c.Header("Access-Control-Allow-Origin", origin)
				c.Header("Access-Control-Allow-Credentials", "true")
			}
			c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE, OPTIONS")
			allowHeaders := defaultAllowHeaders
			if reqHeaders := c.Request.Header.Get("Access-Control-Request-Headers"); reqHeaders != "" {
				// Reflect requested headers (validated) so any header the frontend needs is allowed
				var validated []string
				seen := make(map[string]bool)
				for _, h := range strings.Split(reqHeaders, ",") {
					trimmed := strings.TrimSpace(h)
					key := strings.ToLower(trimmed)
					if trimmed != "" && safeHeaderRe.MatchString(trimmed) && !seen[key] {
						seen[key] = true
						validated = append(validated, trimmed)
					}
				}
				if len(validated) > 0 {
					allowHeaders = defaultAllowHeaders + ", " + strings.Join(validated, ", ")
				}
			}
			c.Header("Access-Control-Allow-Headers", allowHeaders)
			c.Header("Access-Control-Max-Age", "604800")
		}

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
