package main

import (
	"fmt"
	"log"
	"os"

	"github.com/alxand/nalo-workspace/internal/api"

	repository "github.com/alxand/nalo-workspace/internal/repository/postgres"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/joho/godotenv"

	_ "github.com/alxand/nalo-workspace/docs" // swagger docs

	swagger "github.com/gofiber/swagger"
)

func main() {
	// Load env vars, for demo just hardcode DSN and JWT_SECRET here or use os.Getenv
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	dsn := os.Getenv("DSN")
	if dsn == "" {
		log.Fatal("DSN is not set")
	}

	os.Setenv("JWT_SECRET", os.Getenv("JWT_SECRET")) // already required by GenerateJWT

	db := repository.InitDB(dsn)
	repo := repository.NewGormLogRepository(db)

	app := fiber.New()
	app.Use(logger.New())

	// Swagger route
	app.Get("/swagger/*", swagger.HandlerDefault)

	// Public route for getting a JWT token (demo only)
	app.Post("/login", func(c *fiber.Ctx) error {
		token, err := api.GenerateJWT()
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
		}
		return c.JSON(fiber.Map{"token": token})
	})

	// Protected log routes
	api.RegisterLogRoutes(app, repo)

	port := 3000
	fmt.Printf("Starting server on :%d\n", port)
	log.Fatal(app.Listen(fmt.Sprintf(":%d", port)))
}
