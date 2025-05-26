package api

import (
	"strconv"

	"github.com/alxand/nalo-workspace/internal/domain/models"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

// CreateDailyTask godoc
// @Summary Create a new daily log
// @Tags logs
// @Accept json
// @Produce json
// @Param log body models.DailyTask true "Daily log data"
// @Success 201 {object} models.DailyTask
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /logs [post]
func (h *TaskHandler) CreateDailyTask(c *fiber.Ctx) error {
	var log models.DailyTask

	// Parse request body
	if err := c.BodyParser(&log); err != nil {
		return badRequest(c, "Invalid request body", err)
	}

	// Validate request
	if err := validate.Struct(log); err != nil {
		return badRequest(c, "Validation failed", err)
	}

	// Save log
	if err := h.Repo.Create(&log); err != nil {
		return internalServerError(c, "Failed to create log", err)
	}

	return c.Status(fiber.StatusCreated).JSON(log)
}

func (h *TaskHandler) GettasksByDate(c *fiber.Ctx) error {
	date := c.Params("date")
	tasks, err := h.Repo.GetByDate(date)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.JSON(tasks)
}

func (h *TaskHandler) UpdateLog(c *fiber.Ctx) error {
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

func (h *TaskHandler) DeleteLog(c *fiber.Ctx) error {
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

// Helper to respond with a 400 Bad Request
func badRequest(c *fiber.Ctx, message string, err error) error {
	return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
		"error":   message,
		"details": err.Error(),
	})
}

// Helper to respond with a 500 Internal Server Error
func internalServerError(c *fiber.Ctx, message string, err error) error {
	return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
		"error":   message,
		"details": err.Error(),
	})
}
