// internal/api/routes.go
package api

import (
	"os"
	"strconv"

	"github.com/alxand/nalo-workspace/internal/domain/models"
	repository "github.com/alxand/nalo-workspace/internal/repository/postgres"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func RegisterLogRoutes(app *fiber.App, repo *repository.GormLogRepository) {
	// Load .env file
	_ = godotenv.Load()

	handler := &LogHandler{Repo: repo}

	// Load base API path from environment variable, fallback to /logs
	basePath := os.Getenv("LOG_API_PREFIX")
	if basePath == "" {
		basePath = "/logs"
	}

	logs := app.Group(basePath)
	logs.Post("/", handler.CreateLog)
	logs.Get(":date", handler.GetLogsByDate)
	logs.Put(":id", handler.UpdateLog)
	logs.Delete(":id", handler.DeleteLog)
}

type LogHandler struct {
	Repo repository.GormLogRepository
}

func (h *LogHandler) CreateLog(c *fiber.Ctx) error {
	var log models.DailyTask
	if err := c.BodyParser(&log); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if err := h.Repo.Create(&log); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(log)
}

func (h *LogHandler) GetLogsByDate(c *fiber.Ctx) error {
	date := c.Params("date")
	logs, err := h.Repo.GetByDate(date)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(logs)
}

func (h *LogHandler) UpdateLog(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}
	var log models.DailyTask
	if err := c.BodyParser(&log); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	log.ID = id
	if err := h.Repo.Update(&log); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(log)
}

func (h *LogHandler) DeleteLog(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid id"})
	}
	if err := h.Repo.Delete(id); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.SendStatus(fiber.StatusNoContent)
}
