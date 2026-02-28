package main

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

func main() {
	userID := flag.String("user", "", "User ID to embed in the token")
	role := flag.String("role", "", "User role to embed in the token")
	ttl := flag.String("ttl", "24h", "Token TTL (e.g. 24h, 7d, 15m)")
	secret := flag.String("secret", "", "JWT secret (or set JWT_SECRET env var)")
	flag.Parse()

	if *userID == "" {
		fmt.Fprintln(os.Stderr, "Usage: tokengen -user <user_id>")
		os.Exit(1)
	}

	jwtSecret := *secret
	if jwtSecret == "" {
		jwtSecret = os.Getenv("JWT_SECRET")
	}
	if jwtSecret == "" {
		jwtSecret = readSecretFile("shared/secret")
	}
	if jwtSecret == "" {
		fmt.Fprintln(os.Stderr, "Missing JWT secret. Use -secret, set JWT_SECRET, or add shared/secret.")
		os.Exit(1)
	}

	duration, err := time.ParseDuration(*ttl)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Invalid ttl: %v\n", err)
		os.Exit(1)
	}

	now := time.Now().UTC()
	claims := jwt.MapClaims{
		"sub":  *userID,
		"role": *role,
		"iat":  now.Unix(),
		"exp":  now.Add(duration).Unix(),
	}

	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := jwtToken.SignedString([]byte(jwtSecret))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to sign token: %v\n", err)
		os.Exit(1)
	}

	fmt.Println(signed)
}

func readSecretFile(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(data))
}
