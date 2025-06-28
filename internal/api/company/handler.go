package company

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

type CompanyHandler struct {
	Repo   interfaces.CompanyInterface
	Logger *zap.Logger
}

func NewCompanyHandler(repo interfaces.CompanyInterface, logger *zap.Logger) *CompanyHandler {
	return &CompanyHandler{
		Repo:   repo,
		Logger: logger,
	}
}

// CreateCompany godoc
// @Summary Create a new company
// @Tags companies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param company body models.Company true "Company data"
// @Success 201 {object} models.Company
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /companies [post]
func (h *CompanyHandler) CreateCompany(c *fiber.Ctx) error {
	var company models.Company

	if err := c.BodyParser(&company); err != nil {
		h.Logger.Error("Failed to parse request body", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
	}

	if err := validation.ValidateCompany(&company); err != nil {
		h.Logger.Error("Validation failed", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error":   "Validation failed",
			"details": err.Error(),
		})
	}

	if err := h.Repo.Create(&company); err != nil {
		h.Logger.Error("Failed to create company", zap.Error(err))
		return errors.DatabaseError("Failed to create company", err)
	}

	h.Logger.Info("Company created successfully", zap.Int64("company_id", company.ID))
	return c.Status(fiber.StatusCreated).JSON(company)
}

// GetCompany godoc
// @Summary Get a company by ID
// @Tags companies
// @Produce json
// @Security BearerAuth
// @Param id path int true "Company ID"
// @Success 200 {object} models.Company
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /companies/{id} [get]
func (h *CompanyHandler) GetCompany(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid company ID"})
	}

	company, err := h.Repo.GetByID(id)
	if err != nil {
		h.Logger.Error("Failed to get company", zap.Int64("company_id", id), zap.Error(err))
		return errors.NotFound("Company not found", err)
	}

	return c.JSON(company)
}

// GetCompanyByCode godoc
// @Summary Get a company by code
// @Tags companies
// @Produce json
// @Security BearerAuth
// @Param code path string true "Company code"
// @Success 200 {object} models.Company
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /companies/code/{code} [get]
func (h *CompanyHandler) GetCompanyByCode(c *fiber.Ctx) error {
	codeParam := c.Params("code")
	// URL decode the parameter first, then trim whitespace
	decodedCode, err := url.QueryUnescape(codeParam)
	if err != nil {
		decodedCode = codeParam // fallback to original if decoding fails
	}
	code := strings.TrimSpace(decodedCode)
	if code == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Company code is required"})
	}

	company, err := h.Repo.GetByCode(code)
	if err != nil {
		h.Logger.Error("Failed to get company by code", zap.String("code", code), zap.Error(err))
		return errors.NotFound("Company not found", err)
	}

	return c.JSON(company)
}

// GetCompaniesByCountry godoc
// @Summary Get companies by country ID
// @Tags companies
// @Produce json
// @Security BearerAuth
// @Param countryId path int true "Country ID"
// @Success 200 {array} models.Company
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /companies/country/{countryId} [get]
func (h *CompanyHandler) GetCompaniesByCountry(c *fiber.Ctx) error {
	countryIDParam := c.Params("countryId")
	countryID, err := strconv.ParseInt(countryIDParam, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid country ID"})
	}

	companies, err := h.Repo.GetByCountry(countryID)
	if err != nil {
		h.Logger.Error("Failed to get companies by country", zap.Int64("country_id", countryID), zap.Error(err))
		return errors.DatabaseError("Failed to get companies", err)
	}

	h.Logger.Info("Companies retrieved by country", zap.Int64("country_id", countryID), zap.Int("count", len(companies)))
	return c.JSON(companies)
}

// GetCompaniesByIndustry godoc
// @Summary Get companies by industry
// @Tags companies
// @Produce json
// @Security BearerAuth
// @Param industry path string true "Industry name"
// @Success 200 {array} models.Company
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /companies/industry/{industry} [get]
func (h *CompanyHandler) GetCompaniesByIndustry(c *fiber.Ctx) error {
	industryParam := c.Params("industry")
	// URL decode the parameter first, then trim whitespace
	decodedIndustry, err := url.QueryUnescape(industryParam)
	if err != nil {
		decodedIndustry = industryParam // fallback to original if decoding fails
	}
	industry := strings.TrimSpace(decodedIndustry)
	if industry == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Industry is required"})
	}

	companies, err := h.Repo.GetByIndustry(industry)
	if err != nil {
		h.Logger.Error("Failed to get companies by industry", zap.String("industry", industry), zap.Error(err))
		return errors.DatabaseError("Failed to get companies", err)
	}

	h.Logger.Info("Companies retrieved by industry", zap.String("industry", industry), zap.Int("count", len(companies)))
	return c.JSON(companies)
}

// GetAllCompanies godoc
// @Summary Get all companies
// @Tags companies
// @Produce json
// @Security BearerAuth
// @Success 200 {array} models.Company
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /companies [get]
func (h *CompanyHandler) GetAllCompanies(c *fiber.Ctx) error {
	companies, err := h.Repo.GetAll()
	if err != nil {
		h.Logger.Error("Failed to get companies", zap.Error(err))
		return errors.DatabaseError("Failed to get companies", err)
	}

	h.Logger.Info("Companies retrieved successfully", zap.Int("count", len(companies)))
	return c.JSON(companies)
}

// UpdateCompany godoc
// @Summary Update a company
// @Tags companies
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Company ID"
// @Param company body models.Company true "Updated company data"
// @Success 200 {object} models.Company
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /companies/{id} [put]
func (h *CompanyHandler) UpdateCompany(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid company ID"})
	}

	var company models.Company
	if err := c.BodyParser(&company); err != nil {
		h.Logger.Error("Failed to parse request body", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	company.ID = id

	if err := validation.ValidateCompany(&company); err != nil {
		h.Logger.Error("Validation failed", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Validation failed", "details": err.Error()})
	}

	updatedCompany, err := h.Repo.Update(&company)
	if err != nil {
		h.Logger.Error("Failed to update company", zap.Int64("company_id", id), zap.Error(err))
		return errors.DatabaseError("Failed to update company", err)
	}

	h.Logger.Info("Company updated successfully", zap.Int64("company_id", id))
	return c.JSON(updatedCompany)
}

// DeleteCompany godoc
// @Summary Delete a company
// @Tags companies
// @Security BearerAuth
// @Param id path int true "Company ID"
// @Success 204 "No Content"
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /companies/{id} [delete]
func (h *CompanyHandler) DeleteCompany(c *fiber.Ctx) error {
	idParam := c.Params("id")
	id, err := strconv.ParseInt(idParam, 10, 64)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid company ID"})
	}

	if err := h.Repo.Delete(id); err != nil {
		h.Logger.Error("Failed to delete company", zap.Int64("company_id", id), zap.Error(err))
		return errors.DatabaseError("Failed to delete company", err)
	}

	h.Logger.Info("Company deleted successfully", zap.Int64("company_id", id))
	return c.SendStatus(fiber.StatusNoContent)
}
