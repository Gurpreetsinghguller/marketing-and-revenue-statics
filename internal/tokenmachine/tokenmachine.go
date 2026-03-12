package tokenmachine

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/common/config"
	cmnerr "github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/common/errors"
	"github.com/Gurpreetsinghguller/marketing-and-revenue-statics/internal/common/logger"
	"github.com/golang-jwt/jwt/v5"
	"github.com/sirupsen/logrus"
)

type TokenMachine interface {
	GenerateToken(userID string, role string) (string, error)
	ValidateToken(token string) (string, string, error)
}

type JWTTokenMachine struct {
	cfg *config.Config
	log *logrus.Logger
}

var DefaultTTL = 24 * time.Hour

type tokenClaims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

// CustomClaims defines JWT claims used by the API.
type CustomClaims struct {
	Role string `json:"role"`
	jwt.RegisteredClaims
}

func NewJWTTokenMachine(cfg *config.Config) TokenMachine {
	return &JWTTokenMachine{
		cfg: cfg,
		log: logger.Get(),
	}
}

func (tm *JWTTokenMachine) GenerateToken(userID string, role string) (string, error) {
	secret := tm.loadJWTSecret()
	if secret == "" {
		return "", errors.New("jwt secret is missing")
	}

	now := time.Now().UTC()
	tokenTTL, err := time.ParseDuration(tm.cfg.Auth.TokenTTL)
	if err != nil {
		tm.log.WithError(err).Warn("failed to parse token TTL")
		tokenTTL = DefaultTTL
	}
	claims := tokenClaims{
		Role: string(role),
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(now.Add(tokenTTL)),
		},
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return jwtToken.SignedString([]byte(secret))
}

func (tm *JWTTokenMachine) ValidateToken(token string) (string, string, error) {
	secret := tm.loadJWTSecret()
	if secret == "" {
		return "", "", cmnerr.ErrMissingSecret
	}
	fmt.Println("secret used for JWT validation:", secret) // Debug log - remove in production
	claims := &CustomClaims{}
	parsed, err := jwt.ParseWithClaims(token, claims, func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, cmnerr.ErrInvalidToken
		}
		return []byte(secret), nil
	})
	if err != nil || !parsed.Valid {
		if err != nil {
			tm.log.WithError(err).Warn("token parse failed")
		}
		return "", "", cmnerr.ErrInvalidToken
	}
	if claims.Subject == "" {
		return "", "", cmnerr.ErrInvalidToken
	}

	return claims.Subject, claims.Role, nil
}

func (tm *JWTTokenMachine) loadJWTSecret() string {
	if secret := strings.TrimSpace(os.Getenv("JWT_SECRET")); secret != "" {
		return secret
	}

	secretPath := strings.TrimSpace(tm.cfg.Auth.SecretFile)

	data, err := os.ReadFile(secretPath)
	if err != nil {
		tm.log.WithError(err).WithField("path", secretPath).Warn("failed to read jwt secret file")
		return ""
	}

	return strings.TrimSpace(string(data))
}
