// internal/api/routes.go
package api

import (
	"fmt"
	"os"

	repository "github.com/alxand/nalo-workspace/internal/repository/postgres"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

type TaskHandler struct {
	Repo repository.DailyTaskRepository
}

func RegisterDailyTaskRoutes(app *fiber.App, repo repository.DailyTaskRepository) {
	// Load .env file
	_ = godotenv.Load()
	if err := godotenv.Load(); err != nil {
		fmt.Println("Warning: .env file not found or could not be loaded")
	}

	handler := &TaskHandler{Repo: repo}

	// Load base API path from environment variable, fallback to /tasks
	basePath := os.Getenv("LOG_API_PREFIX")
	if basePath == "" {
		basePath = "/tasks"
	}
	tasks := app.Group(basePath)
	tasks.Post("/", handler.CreateDailyTask)
	tasks.Get("/:date", handler.GettasksByDate)
	tasks.Put("/:id", handler.UpdateLog)
	tasks.Delete("/:id", handler.DeleteLog)
}
