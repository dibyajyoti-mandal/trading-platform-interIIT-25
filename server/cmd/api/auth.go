package main

import (
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret string

func init() {
	jwtSecret = os.Getenv("JWT_SECRET")
	// No need to log.Fatal here, main will handle the check
}

func createToken(subject string) (string, error) {
	claims := jwt.MapClaims{
		"sub": subject,
		"exp": time.Now().Add(24 * time.Hour).Unix(),
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return t.SignedString([]byte(jwtSecret))
}
