package main

import (
	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

type authRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func signupHandler(c *fiber.Ctx) error {
	var r authRequest
	if err := c.BodyParser(&r); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid json"})
	}
	if r.Email == "" || r.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "email & password required"})
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(r.Password), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "server error"})
	}

	_, err = db.Exec(c.Context(), "INSERT INTO users (email, password_hash) VALUES ($1,$2)", r.Email, string(hash))
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "could not create user", "details": err.Error()})
	}

	token, err := createToken(r.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "token error"})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"token": token})
}

func loginHandler(c *fiber.Ctx) error {
	var r authRequest
	if err := c.BodyParser(&r); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid json"})
	}
	var stored string
	err := db.QueryRow(c.Context(), "SELECT password_hash FROM users WHERE email=$1", r.Email).Scan(&stored)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
	}
	if err := bcrypt.CompareHashAndPassword([]byte(stored), []byte(r.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
	}
	token, err := createToken(r.Email)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "token error"})
	}
	return c.JSON(fiber.Map{"token": token})
}
