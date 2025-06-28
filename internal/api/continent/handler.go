package continent

import (
	"net/url"
	"strconv"
	"strings"

	"github.com/alxand/nalo-workspace/internal/domain/interfaces"
	"github.com/alxand/nalo-workspace/internal/domain/models"
	"github.com/alxand/nalo-workspace/internal/pkg/errors"
	"github.com/alxand/nalo-workspace/internal/pkg/validation"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type ContinentHandler struct {
	Repo   interfaces.ContinentInterface
	Logger *zap.Logger
}

func NewContinentHandler(repo interfaces.ContinentInterface, logger *zap.Logger) *ContinentHandler {
	return &ContinentHandler{
		Repo:   repo,
		Logger: logger,
	}
}

// CreateContinent godoc
// @Summary Create a new continent
// @Tags continents
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param continent body models.Continent true "Continent data"
// @Success 201 {object} models.Continent
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /continents [post]
func (h *ContinentHandler) CreateContinent(c *fiber.Ctx) error {
	var continent models.Continent

	if err := c.BodyParser(&continent); err != nil {
		h.Logger.Error("Failed to parse request body", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	if err := validation.ValidateContinent(&continent); err != nil {
		h.Logger.Error("Validation failed", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Validation failed",
			"details": err.Error(),
		})
	}

	if err := h.Repo.Create(&continent); err != nil {
		h.Logger.Error("Failed to create continent", zap.Error(err))
		return errors.DatabaseError("Failed to create continent", err)
	}

	h.Logger.Info("Continent created successfully", zap.Int64("continent_id", continent.ID))
	return c.Status(fiber.StatusCreated).JSON(continent)
}

// GetContinent godoc
// @Summary Get a continent by ID
// @Tags continents
// @Produce json
// @Security BearerAuth
// @Param id path int true "Continent ID"
// @Success 200 {object} models.Continent
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /continents/{id} [get]
func (h *ContinentHandler) GetContinent(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid continent ID"})
	}

	continent, err := h.Repo.GetByID(id)
	if err != nil {
		h.Logger.Error("Failed to get continent", zap.Int64("continent_id", id), zap.Error(err))
		return errors.NotFound("Continent not found", err)
	}

	return c.JSON(continent)
}

// GetContinentByCode godoc
// @Summary Get a continent by code
// @Tags continents
// @Produce json
// @Security BearerAuth
// @Param code path string true "Continent code"
// @Success 200 {object} models.Continent
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /continents/code/{code} [get]
func (h *ContinentHandler) GetContinentByCode(c *fiber.Ctx) error {
	codeParam := c.Params("code")
	// URL decode the parameter first, then trim whitespace
	decodedCode, err := url.QueryUnescape(codeParam)
	if err != nil {
		decodedCode = codeParam // fallback to original if decoding fails
	}
	code := strings.TrimSpace(decodedCode)
	if code == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Continent code is required"})
	}

	continent, err := h.Repo.GetByCode(code)
	if err != nil {
		h.Logger.Error("Failed to get continent by code", zap.String("code", code), zap.Error(err))
		return errors.NotFound("Continent not found", err)
	}

	return c.JSON(continent)
}

// GetAllContinents godoc
// @Summary Get all continents
// @Tags continents
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.Continent
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /continents [get]
func (h *ContinentHandler) GetAllContinents(c *fiber.Ctx) error {
	continents, err := h.Repo.GetAll()
	if err != nil {
		h.Logger.Error("Failed to get continents", zap.Error(err))
		return errors.DatabaseError("Failed to get continents", err)
	}

	h.Logger.Info("Continents retrieved successfully", zap.Int("count", len(continents)))
	return c.JSON(continents)
}

// UpdateContinent godoc
// @Summary Update a continent
// @Tags continents
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Continent ID"
// @Param continent body models.Continent true "Updated continent data"
// @Success 200 {object} models.Continent
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /continents/{id} [put]
func (h *ContinentHandler) UpdateContinent(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid continent ID"})
	}

	var continent models.Continent
	if err := c.BodyParser(&continent); err != nil {
		h.Logger.Error("Failed to parse request body", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	continent.ID = id

	if err := validation.ValidateContinent(&continent); err != nil {
		h.Logger.Error("Validation failed", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Validation failed", "details": err.Error()})
	}

	updatedContinent, err := h.Repo.Update(&continent)
	if err != nil {
		h.Logger.Error("Failed to update continent", zap.Int64("continent_id", id), zap.Error(err))
		return errors.DatabaseError("Failed to update continent", err)
	}

	h.Logger.Info("Continent updated successfully", zap.Int64("continent_id", id))
	return c.JSON(updatedContinent)
}

// DeleteContinent godoc
// @Summary Delete a continent
// @Tags continents
// @Security BearerAuth
// @Param id path int true "Continent ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /continents/{id} [delete]
func (h *ContinentHandler) DeleteContinent(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid continent ID"})
	}

	if err := h.Repo.Delete(id); err != nil {
		h.Logger.Error("Failed to delete continent", zap.Int64("continent_id", id), zap.Error(err))
		return errors.DatabaseError("Failed to delete continent", err)
	}

	h.Logger.Info("Continent deleted successfully", zap.Int64("continent_id", id))
	return c.SendStatus(fiber.StatusNoContent)
}
