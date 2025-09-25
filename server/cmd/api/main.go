package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/gofiber/fiber/v2"
	fiberlogger "github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"

	"github.com/jackc/pgx/v5/pgxpool"
)

var db *pgxpool.Pool

func main() {
	// load .env if present
	_ = godotenv.Load()

	_ = os.Getenv("JWT_SECRET") // jwtSecret is now local to main.go, but not used directly here

	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL must be set")
	}

	ctx := context.Background()
	var err error
	db, err = pgxpool.New(ctx, databaseURL)
	if err != nil {
		log.Fatalf("failed to connect to db: %v", err)
	}
	defer db.Close()

	if err := ensureSchema(ctx); err != nil {
		log.Fatalf("failed to ensure schema: %v", err)
	}

	app := fiber.New()
	app.Use(fiberlogger.New())

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.SendString("ok")
	})

	setupRoutes(app)

	// start server
	go func() {
		if err := app.Listen(":8080"); err != nil {
			log.Fatalf("server failed: %v", err)
		}
	}()

	// graceful shutdown on Ctrl+C
	sigctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()
	<-sigctx.Done()

	shutdownCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	_ = app.Shutdown()
	<-shutdownCtx.Done()
}

func ensureSchema(ctx context.Context) error {
	schema := `
CREATE EXTENSION IF NOT EXISTS "pgcrypto";
CREATE TABLE IF NOT EXISTS users (
  id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
  email TEXT UNIQUE NOT NULL,
  password_hash TEXT NOT NULL,
  created_at TIMESTAMPTZ DEFAULT now()
);
`
	_, err := db.Exec(ctx, schema)
	return err
}
