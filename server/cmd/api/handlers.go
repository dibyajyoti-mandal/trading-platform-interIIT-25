package main

import (
	"log"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
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
	r.Email = strings.TrimSpace(strings.ToLower(r.Email))
	if r.Email == "" || r.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "email & password required"})
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(r.Password), bcrypt.DefaultCost)
	if err != nil {
		log.Printf("signup: bcrypt error: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "server error"})
	}

	_, err = db.Exec(c.Context(), "INSERT INTO users (email, password_hash) VALUES ($1,$2)", r.Email, string(hash))
	if err != nil {
		log.Printf("signup: db insert error: %v\n", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "could not create user", "details": err.Error()})
	}

	token, err := createToken(r.Email)
	if err != nil {
		log.Printf("signup: token error: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "token error"})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"token": token})
}

func loginHandler(c *fiber.Ctx) error {
	var r authRequest
	if err := c.BodyParser(&r); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid json"})
	}
	r.Email = strings.TrimSpace(strings.ToLower(r.Email))
	if r.Email == "" || r.Password == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "email & password required"})
	}

	var stored string
	err := db.QueryRow(c.Context(), "SELECT password_hash FROM users WHERE email=$1", r.Email).Scan(&stored)
	if err != nil {
		if err == pgx.ErrNoRows {
			log.Printf("login failed: no user found for email=%s\n", r.Email)
			// keep response generic for security
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
		}
		// some other DB error
		log.Printf("login: db query error: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "server error"})
	}

	if err := bcrypt.CompareHashAndPassword([]byte(stored), []byte(r.Password)); err != nil {
		log.Printf("login failed: bcrypt mismatch for email=%s (err=%v)\n", r.Email, err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid credentials"})
	}

	token, err := createToken(r.Email)
	if err != nil {
		log.Printf("login: token error: %v\n", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "token error"})
	}
	return c.JSON(fiber.Map{"token": token})
}

func topupHandler(c *fiber.Ctx) error {
	return c.SendString("topupHandler placeholder")
}

func getWalletHandler(c *fiber.Ctx) error {
	return c.SendString("getWalletHandler placeholder")
}

func placeOrderHandler(c *fiber.Ctx) error {
	return c.SendString("placeOrderHandler placeholder")
}

func getOrderHandler(c *fiber.Ctx) error {
	return c.SendString("getOrderHandler placeholder")
}
