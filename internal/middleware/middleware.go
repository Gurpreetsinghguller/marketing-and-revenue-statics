package middleware

import (
	"context"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
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
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "", "", ErrMissingSecret
	}

	claims := &CustomClaims{}
	parsed, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return []byte(secret), nil
	})
	if err != nil || !parsed.Valid {
		return "", "", ErrInvalidToken
	}
	if claims.Subject == "" {
		return "", "", ErrInvalidToken
	}

	return claims.Subject, claims.Role, nil
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
				if userRole == role {
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
	// Simple rate limiter using in-memory map (in production, use Redis or similar)
	type RateLimitEntry struct {
		count     int
		resetTime time.Time
	}

	limitMap := make(map[string]*RateLimitEntry)
	const maxRequests = 100
	const windowDuration = 1 * time.Minute

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get identifier (user ID or IP)
		identifier := r.Header.Get("X-User-ID")
		if identifier == "" {
			identifier = getClientIP(r)
		}

		// Check current limit
		now := time.Now()
		entry, exists := limitMap[identifier]

		if !exists || now.After(entry.resetTime) {
			// New window or expired window
			limitMap[identifier] = &RateLimitEntry{
				count:     1,
				resetTime: now.Add(windowDuration),
			}
		} else if entry.count >= maxRequests {
			// Rate limit exceeded
			w.Header().Set("Retry-After", "60")
			http.Error(w, `{"error": "Rate limit exceeded"}`, http.StatusTooManyRequests)
			return
		} else {
			// Increment counter
			entry.count++
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

		// Log request details (in production, use proper logging library)
		_ = r.Method + " " + r.RequestURI

		next.ServeHTTP(w, r)

		// Log response time
		duration := time.Since(start)
		_ = duration // Use duration in actual logging
	})
}

// Helper functions

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
