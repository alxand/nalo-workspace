package country

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

type CountryHandler struct {
	Repo   interfaces.CountryInterface
	Logger *zap.Logger
}

func NewCountryHandler(repo interfaces.CountryInterface, logger *zap.Logger) *CountryHandler {
	return &CountryHandler{
		Repo:   repo,
		Logger: logger,
	}
}

// CreateCountry godoc
// @Summary Create a new country
// @Tags countries
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param country body models.Country true "Country data"
// @Success 201 {object} models.Country
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /countries [post]
func (h *CountryHandler) CreateCountry(c *fiber.Ctx) error {
	var country models.Country

	if err := c.BodyParser(&country); err != nil {
		h.Logger.Error("Failed to parse request body", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	if err := validation.ValidateCountry(&country); err != nil {
		h.Logger.Error("Validation failed", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Validation failed",
			"details": err.Error(),
		})
	}

	if err := h.Repo.Create(&country); err != nil {
		h.Logger.Error("Failed to create country", zap.Error(err))
		return errors.DatabaseError("Failed to create country", err)
	}

	h.Logger.Info("Country created successfully", zap.Int64("country_id", country.ID))
	return c.Status(fiber.StatusCreated).JSON(country)
}

// GetCountry godoc
// @Summary Get a country by ID
// @Tags countries
// @Produce json
// @Security BearerAuth
// @Param id path int true "Country ID"
// @Success 200 {object} models.Country
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /countries/{id} [get]
func (h *CountryHandler) GetCountry(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid country ID"})
	}

	country, err := h.Repo.GetByID(id)
	if err != nil {
		h.Logger.Error("Failed to get country", zap.Int64("country_id", id), zap.Error(err))
		return errors.NotFound("Country not found", err)
	}

	return c.JSON(country)
}

// GetCountryByCode godoc
// @Summary Get a country by code
// @Tags countries
// @Produce json
// @Security BearerAuth
// @Param code path string true "Country code"
// @Success 200 {object} models.Country
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /countries/code/{code} [get]
func (h *CountryHandler) GetCountryByCode(c *fiber.Ctx) error {
	codeParam := c.Params("code")
	// URL decode the parameter first, then trim whitespace
	decodedCode, err := url.QueryUnescape(codeParam)
	if err != nil {
		decodedCode = codeParam // fallback to original if decoding fails
	}
	code := strings.TrimSpace(decodedCode)
	if code == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Country code is required"})
	}

	country, err := h.Repo.GetByCode(code)
	if err != nil {
		h.Logger.Error("Failed to get country by code", zap.String("code", code), zap.Error(err))
		return errors.NotFound("Country not found", err)
	}

	return c.JSON(country)
}

// GetCountriesByContinent godoc
// @Summary Get countries by continent ID
// @Tags countries
// @Produce json
// @Security BearerAuth
// @Param continentId path int true "Continent ID"
// @Success 200 {array} models.Country
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /countries/continent/{continentId} [get]
func (h *CountryHandler) GetCountriesByContinent(c *fiber.Ctx) error {
	continentIDParam := c.Params("continentId")
	continentID, err := strconv.ParseInt(continentIDParam, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid continent ID"})
	}

	countries, err := h.Repo.GetByContinent(continentID)
	if err != nil {
		h.Logger.Error("Failed to get countries by continent", zap.Int64("continent_id", continentID), zap.Error(err))
		return errors.DatabaseError("Failed to get countries", err)
	}

	h.Logger.Info("Countries retrieved by continent", zap.Int64("continent_id", continentID), zap.Int("count", len(countries)))
	return c.JSON(countries)
}

// GetAllCountries godoc
// @Summary Get all countries
// @Tags countries
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.Country
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /countries [get]
func (h *CountryHandler) GetAllCountries(c *fiber.Ctx) error {
	countries, err := h.Repo.GetAll()
	if err != nil {
		h.Logger.Error("Failed to get countries", zap.Error(err))
		return errors.DatabaseError("Failed to get countries", err)
	}

	h.Logger.Info("Countries retrieved successfully", zap.Int("count", len(countries)))
	return c.JSON(countries)
}

// UpdateCountry godoc
// @Summary Update a country
// @Tags countries
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Country ID"
// @Param country body models.Country true "Updated country data"
// @Success 200 {object} models.Country
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /countries/{id} [put]
func (h *CountryHandler) UpdateCountry(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid country ID"})
	}

	var country models.Country
	if err := c.BodyParser(&country); err != nil {
		h.Logger.Error("Failed to parse request body", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	country.ID = id

	if err := validation.ValidateCountry(&country); err != nil {
		h.Logger.Error("Validation failed", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Validation failed", "details": err.Error()})
	}

	updatedCountry, err := h.Repo.Update(&country)
	if err != nil {
		h.Logger.Error("Failed to update country", zap.Int64("country_id", id), zap.Error(err))
		return errors.DatabaseError("Failed to update country", err)
	}

	h.Logger.Info("Country updated successfully", zap.Int64("country_id", id))
	return c.JSON(updatedCountry)
}

// DeleteCountry godoc
// @Summary Delete a country
// @Tags countries
// @Security BearerAuth
// @Param id path int true "Country ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /countries/{id} [delete]
func (h *CountryHandler) DeleteCountry(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid country ID"})
	}

	if err := h.Repo.Delete(id); err != nil {
		h.Logger.Error("Failed to delete country", zap.Int64("country_id", id), zap.Error(err))
		return errors.DatabaseError("Failed to delete country", err)
	}

	h.Logger.Info("Country deleted successfully", zap.Int64("country_id", id))
	return c.SendStatus(fiber.StatusNoContent)
}
