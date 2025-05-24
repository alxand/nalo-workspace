package handler

import (
	"github.com/alxand/nalo-workspace/internal/domain/models"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber"
)

var validate = validator.New()

// CreateLog godoc
// @Summary Create a new daily log
// @Tags logs
// @Accept json
// @Produce json
// @Param log body models.DailyTask true "Daily log data"
// @Success 201 {object} models.DailyTask
// @Failure 400 {object} map[string]string
// @Router /logs [post]
func (h *LogHandler) CreateLog(c *fiber.Ctx) error {
	var log models.DailyTask
	if err := c.BodyParser(&log); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}
	if err := validate.Struct(log); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"validation_error": err.Error()})
	}
	if err := h.Repo.Create(&log); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": err.Error()})
	}
	return c.Status(fiber.StatusCreated).JSON(log)
}
