package main

import (
	"github.com/gofiber/fiber/v2"
)

func setupRoutes(app *fiber.App) {
	api := app.Group("/api/v1")
	api.Post("/signup", signupHandler)
	api.Post("/login", loginHandler)

	api.Post("/wallet/topup", authRequired(topupHandler)) // requires authentication
	api.Get("/wallet", authRequired(getWalletHandler))

	api.Post("/orders", authRequired(placeOrderHandler))
	api.Get("/orders/:id", authRequired(getOrderHandler))
}
