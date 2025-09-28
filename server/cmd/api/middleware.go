package main

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func getUserIDByEmail(ctx context.Context, email string) (string, error) {
	var id string
	err := db.QueryRow(ctx, "SELECT id FROM users WHERE email=$1", email).Scan(&id)
	return id, err
}

// authRequired extracts Bearer token, validates it, and sets "user_id" in Locals
func authRequired(next fiber.Handler) fiber.Handler {
	return func(c *fiber.Ctx) error {
		auth := c.Get("Authorization")
		if auth == "" || !strings.HasPrefix(auth, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "missing bearer token"})
		}
		tokenStr := strings.TrimPrefix(auth, "Bearer ")
		token, err := jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
			// no alg check here for brevity; ensure in production
			return []byte(os.Getenv("JWT_SECRET")), nil
		})
		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token"})
		}
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token claims"})
		}
		sub, _ := claims["sub"].(string)
		sub = strings.TrimSpace(strings.ToLower(sub))
		userID, err := getUserIDByEmail(c.Context(), sub)
		if err != nil {
			log.Printf("authRequired: user lookup failed: %v\n", err)
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid token user"})
		}
		c.Locals("user_id", userID)
		return next(c)
	}
}
