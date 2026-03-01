package middleware

import (
	"context"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/common/config"
	"github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/common/errors"
	"github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/common/logger"
	"github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/common/util"
	"github.com/golang-jwt/jwt/v5"
)

var (
	middlewareLog  = logger.Get()
	configLoadOnce sync.Once
	appConfig      *config.Config
)

// AuthMiddleware validates JWT token and extracts user info
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Extract JWT token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, `{"error": "Missing authorization header"}`, http.StatusUnauthorized)
			return
		}

		// Expected format: "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			http.Error(w, `{"error": "Invalid authorization header format"}`, http.StatusUnauthorized)
			return
		}

		token := parts[1]

		// Validate token
		userID, userRole, err := validateToken(token)
		if err != nil {
			middlewareLog.WithError(err).Warn("invalid or expired auth token")
			http.Error(w, `{"error": "Invalid or expired token"}`, http.StatusUnauthorized)
			return
		}

		// Add user info to request context
		ctx := context.WithValue(r.Context(), "user_id", userID)
		// Also add to headers for easier access in handlers
		r.Header.Set("X-User-ID", userID)
		if userRole != "" {
			r.Header.Set("X-User-Role", userRole)
		}

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// validateToken validates a JWT token and extracts user info.
func validateToken(token string) (string, string, error) {
	secret := loadJWTSecret()
	if secret == "" {
		return "", "", errors.ErrMissingSecret
	}

	claims := &CustomClaims{}
	parsed, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, errors.ErrInvalidToken
		}
		return []byte(secret), nil
	})
	if err != nil || !parsed.Valid {
		if err != nil {
			middlewareLog.WithError(err).Warn("token parse failed")
		}
		return "", "", errors.ErrInvalidToken
	}
	if claims.Subject == "" {
		return "", "", errors.ErrInvalidToken
	}

	return claims.Subject, claims.Role, nil
}

func loadJWTSecret() string {
	if secret := strings.TrimSpace(os.Getenv("JWT_SECRET")); secret != "" {
		return secret
	}

	cfg := getAppConfig()
	secretPath := strings.TrimSpace(cfg.Auth.SecretFile)
	if secretPath == "" {
		secretPath = "shared/secret"
	}

	data, err := os.ReadFile(secretPath)
	if err != nil {
		middlewareLog.WithError(err).WithField("path", secretPath).Warn("failed to read jwt secret file")
		return ""
	}

	return strings.TrimSpace(string(data))
}

// RoleMiddleware checks user role authorization
func RoleMiddleware(allowedRoles ...string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract user role from request context or header
			userRole := r.Header.Get("X-User-Role")
			if userRole == "" {
				http.Error(w, `{"error": "User role not found"}`, http.StatusForbidden)
				return
			}

			// Check if role is in allowedRoles
			isAllowed := false
			for _, role := range allowedRoles {
				if strings.EqualFold(strings.ToLower(userRole), strings.ToLower(role)) {
					isAllowed = true
					break
				}
			}

			if !isAllowed {
				http.Error(w, `{"error": "User role not authorized for this resource"}`, http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// RateLimitMiddleware applies rate limiting
func RateLimitMiddleware(next http.Handler) http.Handler {
	cfg := getAppConfig()
	limiter := util.NewFixedWindowRateLimiter(cfg.RateLimit.MaxRequests, time.Duration(cfg.RateLimit.WindowSeconds)*time.Second)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get identifier (user ID or IP)
		identifier := r.Header.Get("X-User-ID")
		if identifier == "" {
			identifier = getClientIP(r)
		}

		allowed, retryAfter := limiter.Allow(identifier, time.Now())
		if !allowed {
			w.Header().Set("Retry-After", strconv.Itoa(retryAfter))
			middlewareLog.WithFields(map[string]interface{}{
				"identifier":  identifier,
				"retry_after": retryAfter,
				"path":        r.URL.Path,
				"method":      r.Method,
			}).Warn("rate limit exceeded")
			http.Error(w, `{"error": "Rate limit exceeded"}`, http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// CORSMiddleware handles CORS headers
func CORSMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization, X-User-ID, X-User-Role")
		w.Header().Set("Access-Control-Max-Age", "3600")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// LoggingMiddleware logs HTTP requests
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		writer := &statusResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}

		next.ServeHTTP(writer, r)

		middlewareLog.WithFields(map[string]interface{}{
			"method":      r.Method,
			"path":        r.URL.Path,
			"status":      writer.statusCode,
			"duration_ms": time.Since(start).Milliseconds(),
			"client_ip":   getClientIP(r),
			"user_id":     r.Header.Get("X-User-ID"),
		}).Info("http request")
	})
}

// CustomClaims defines JWT claims used by the API.
type CustomClaims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

// getClientIP extracts client IP from request
func getClientIP(r *http.Request) string {
	// Check X-Forwarded-For header first (for proxied requests)
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		return strings.Split(forwarded, ",")[0]
	}

	// Falls back to RemoteAddr
	return strings.Split(r.RemoteAddr, ":")[0]
}

func getAppConfig() *config.Config {
	configLoadOnce.Do(func() {
		cfg, err := config.Load(config.DefaultConfigPath)
		if err != nil {
			middlewareLog.WithError(err).Warn("failed to load middleware config; using defaults")
			cfg = config.Default()
		}
		appConfig = cfg
	})

	return appConfig
}

type statusResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func (w *statusResponseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}
