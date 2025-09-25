package main

import (
	"github.com/gofiber/fiber/v2"
)

func setupRoutes(app *fiber.App) {
	api := app.Group("/api/v1")
	api.Post("/signup", signupHandler)
	api.Post("/login", loginHandler)
}
