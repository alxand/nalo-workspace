package company

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alxand/nalo-workspace/internal/domain/models"
	"github.com/alxand/nalo-workspace/internal/pkg/logger"
	"github.com/alxand/nalo-workspace/internal/pkg/middleware"
	"github.com/alxand/nalo-workspace/internal/pkg/validation"
	"github.com/alxand/nalo-workspace/internal/test_helpers"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// MockCompanyRepository is a mock implementation of CompanyInterface
type MockCompanyRepository struct {
	mock.Mock
}

func (m *MockCompanyRepository) Create(company *models.Company) error {
	args := m.Called(company)
	return args.Error(0)
}

func (m *MockCompanyRepository) GetByID(id int64) (*models.Company, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Company), args.Error(1)
}

func (m *MockCompanyRepository) GetByCode(code string) (*models.Company, error) {
	args := m.Called(code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Company), args.Error(1)
}

func (m *MockCompanyRepository) GetByCountry(countryID int64) ([]models.Company, error) {
	args := m.Called(countryID)
	return args.Get(0).([]models.Company), args.Error(1)
}

func (m *MockCompanyRepository) GetByIndustry(industry string) ([]models.Company, error) {
	args := m.Called(industry)
	return args.Get(0).([]models.Company), args.Error(1)
}

func (m *MockCompanyRepository) GetAll() ([]models.Company, error) {
	args := m.Called()
	return args.Get(0).([]models.Company), args.Error(1)
}

func (m *MockCompanyRepository) Update(company *models.Company) (*models.Company, error) {
	args := m.Called(company)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Company), args.Error(1)
}

func (m *MockCompanyRepository) Delete(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func setupTestApp() *fiber.App {
	logger := zap.NewNop()
	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler(logger),
	})
	return app
}

func setupTestHandler() (*CompanyHandler, *MockCompanyRepository) {
	// Initialize validation for all tests
	validation.Init()

	mockRepo := new(MockCompanyRepository)
	logger := zap.NewNop()
	handler := NewCompanyHandler(mockRepo, logger)
	return handler, mockRepo
}

func TestCreateCompany_Success(t *testing.T) {
	// Setup test database
	testDB, err := test_helpers.SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}
	defer test_helpers.CleanupTestDB(testDB)

	// Create test continent and country first
	continent, err := test_helpers.CreateTestContinent(testDB, "Europe", "EU")
	if err != nil {
		t.Fatalf("Failed to create test continent: %v", err)
	}

	country, err := test_helpers.CreateTestCountry(testDB, "Germany", "DEU", continent.ID)
	if err != nil {
		t.Fatalf("Failed to create test country: %v", err)
	}

	// Setup app and handler
	app := test_helpers.SetupTestApp()
	handler := NewCompanyHandler(testDB.CompanyRepo, logger.Get())

	company := models.Company{
		Name:        "TechCorp",
		Code:        "TECH001",
		CountryID:   country.ID,
		Description: "Technology company",
		Website:     "https://techcorp.com",
		Industry:    "Technology",
		Size:        "large",
		Founded:     2010,
	}

	app.Post("/companies", handler.CreateCompany)

	body, _ := json.Marshal(company)
	req := httptest.NewRequest("POST", "/companies", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

	var responseCompany models.Company
	json.NewDecoder(resp.Body).Decode(&responseCompany)
	assert.Equal(t, company.Name, responseCompany.Name)
	assert.Equal(t, company.Code, responseCompany.Code)
	assert.Equal(t, company.CountryID, responseCompany.CountryID)
	assert.NotZero(t, responseCompany.ID) // Should have been assigned by database
}

func TestCreateCompany_InvalidJSON(t *testing.T) {
	app := setupTestApp()
	handler, mockRepo := setupTestHandler()

	app.Post("/companies", handler.CreateCompany)

	req := httptest.NewRequest("POST", "/companies", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Contains(t, response["error"], "Invalid request body")

	mockRepo.AssertNotCalled(t, "Create")
}

func TestCreateCompany_ValidationError(t *testing.T) {
	app := setupTestApp()
	handler, mockRepo := setupTestHandler()

	company := models.Company{
		Name:      "", // Invalid: name is required
		Code:      "TECH001",
		CountryID: 1,
		Size:      "large",
	}

	app.Post("/companies", handler.CreateCompany)

	body, _ := json.Marshal(company)
	req := httptest.NewRequest("POST", "/companies", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Contains(t, response["error"], "Validation failed")

	mockRepo.AssertNotCalled(t, "Create")
}

func TestCreateCompany_InvalidSize(t *testing.T) {
	app := setupTestApp()
	handler, mockRepo := setupTestHandler()

	company := models.Company{
		Name:        "TechCorp",
		Code:        "TECH001",
		CountryID:   1,
		Description: "Technology company",
		Size:        "invalid_size", // Invalid size
	}

	app.Post("/companies", handler.CreateCompany)

	body, _ := json.Marshal(company)
	req := httptest.NewRequest("POST", "/companies", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Contains(t, response["error"], "Validation failed")

	mockRepo.AssertNotCalled(t, "Create")
}

func TestCreateCompany_DatabaseError(t *testing.T) {
	app := setupTestApp()
	handler, mockRepo := setupTestHandler()

	company := models.Company{
		Name:        "TechCorp",
		Code:        "TECH001",
		CountryID:   1,
		Description: "Technology company",
		Website:     "https://techcorp.com",
		Industry:    "Technology",
		Size:        "large",
		Founded:     2010,
	}

	mockRepo.On("Create", mock.AnythingOfType("*models.Company")).Return(errors.New("database error")).Once()

	app.Post("/companies", handler.CreateCompany)

	body, _ := json.Marshal(company)
	req := httptest.NewRequest("POST", "/companies", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	mockRepo.AssertExpectations(t)
}

func TestGetCompany_Success(t *testing.T) {
	app := setupTestApp()
	handler, mockRepo := setupTestHandler()

	expectedCompany := &models.Company{
		ID:          1,
		Name:        "TechCorp",
		Code:        "TECH001",
		CountryID:   1,
		Description: "Technology company",
		Website:     "https://techcorp.com",
		Industry:    "Technology",
		Size:        "large",
		Founded:     2010,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	mockRepo.On("GetByID", int64(1)).Return(expectedCompany, nil).Once()

	app.Get("/companies/:id", handler.GetCompany)

	req := httptest.NewRequest("GET", "/companies/1", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var responseCompany models.Company
	json.NewDecoder(resp.Body).Decode(&responseCompany)
	assert.Equal(t, expectedCompany.ID, responseCompany.ID)
	assert.Equal(t, expectedCompany.Name, responseCompany.Name)
	assert.Equal(t, expectedCompany.Code, responseCompany.Code)

	mockRepo.AssertExpectations(t)
}

func TestGetCompany_InvalidID(t *testing.T) {
	app := setupTestApp()
	handler, mockRepo := setupTestHandler()

	app.Get("/companies/:id", handler.GetCompany)

	req := httptest.NewRequest("GET", "/companies/invalid", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Contains(t, response["error"], "Invalid company ID")

	mockRepo.AssertNotCalled(t, "GetByID")
}

func TestGetCompany_NotFound(t *testing.T) {
	app := setupTestApp()
	handler, mockRepo := setupTestHandler()

	mockRepo.On("GetByID", int64(999)).Return(nil, errors.New("not found")).Once()

	app.Get("/companies/:id", handler.GetCompany)

	req := httptest.NewRequest("GET", "/companies/999", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

	mockRepo.AssertExpectations(t)
}

func TestGetCompanyByCode_Success(t *testing.T) {
	app := setupTestApp()
	handler, mockRepo := setupTestHandler()

	expectedCompany := &models.Company{
		ID:          1,
		Name:        "TechCorp",
		Code:        "TECH001",
		CountryID:   1,
		Description: "Technology company",
		Website:     "https://techcorp.com",
		Industry:    "Technology",
		Size:        "large",
		Founded:     2010,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	mockRepo.On("GetByCode", "TECH001").Return(expectedCompany, nil).Once()

	app.Get("/companies/code/:code", handler.GetCompanyByCode)

	req := httptest.NewRequest("GET", "/companies/code/TECH001", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var responseCompany models.Company
	json.NewDecoder(resp.Body).Decode(&responseCompany)
	assert.Equal(t, expectedCompany.Code, responseCompany.Code)

	mockRepo.AssertExpectations(t)
}

// TestGetCompanyByCode_EmptyString tests with an actual empty string parameter using SQLite
func TestGetCompanyByCode_EmptyString(t *testing.T) {
	// Setup test database
	testDB, err := test_helpers.SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}
	defer test_helpers.CleanupTestDB(testDB)

	// Create test continent and country first
	continent, err := test_helpers.CreateTestContinent(testDB, "Europe", "EU")
	if err != nil {
		t.Fatalf("Failed to create test continent: %v", err)
	}

	country, err := test_helpers.CreateTestCountry(testDB, "Germany", "DEU", continent.ID)
	if err != nil {
		t.Fatalf("Failed to create test country: %v", err)
	}

	// Create a test company (needed for database setup)
	_, err = test_helpers.CreateTestCompany(testDB, "TechCorp", "TECH001", country.ID)
	if err != nil {
		t.Fatalf("Failed to create test company: %v", err)
	}

	// Setup app and handler
	app := test_helpers.SetupTestApp()
	handler := NewCompanyHandler(testDB.CompanyRepo, logger.Get())

	app.Get("/companies/code/:code", handler.GetCompanyByCode)

	// Test with a space character which should be treated as empty after trimming
	req := httptest.NewRequest("GET", "/companies/code/%20", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Contains(t, response["error"], "Company code is required")
}

func TestGetCompanyByCode_NotFound(t *testing.T) {
	app := setupTestApp()
	handler, mockRepo := setupTestHandler()

	mockRepo.On("GetByCode", "INVALID").Return(nil, errors.New("not found")).Once()

	app.Get("/companies/code/:code", handler.GetCompanyByCode)

	req := httptest.NewRequest("GET", "/companies/code/INVALID", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

	mockRepo.AssertExpectations(t)
}

func TestGetCompaniesByCountry_Success(t *testing.T) {
	app := setupTestApp()
	handler, mockRepo := setupTestHandler()

	expectedCompanies := []models.Company{
		{
			ID:          1,
			Name:        "TechCorp",
			Code:        "TECH001",
			CountryID:   1,
			Description: "Technology company",
			Website:     "https://techcorp.com",
			Industry:    "Technology",
			Size:        "large",
			Founded:     2010,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          2,
			Name:        "DataCorp",
			Code:        "DATA001",
			CountryID:   1,
			Description: "Data analytics company",
			Website:     "https://datacorp.com",
			Industry:    "Technology",
			Size:        "medium",
			Founded:     2015,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	mockRepo.On("GetByCountry", int64(1)).Return(expectedCompanies, nil).Once()

	app.Get("/companies/country/:countryId", handler.GetCompaniesByCountry)

	req := httptest.NewRequest("GET", "/companies/country/1", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var responseCompanies []models.Company
	json.NewDecoder(resp.Body).Decode(&responseCompanies)
	assert.Len(t, responseCompanies, 2)
	assert.Equal(t, expectedCompanies[0].Name, responseCompanies[0].Name)
	assert.Equal(t, expectedCompanies[1].Name, responseCompanies[1].Name)

	mockRepo.AssertExpectations(t)
}

func TestGetCompaniesByCountry_InvalidCountryID(t *testing.T) {
	app := setupTestApp()
	handler, mockRepo := setupTestHandler()

	app.Get("/companies/country/:countryId", handler.GetCompaniesByCountry)

	req := httptest.NewRequest("GET", "/companies/country/invalid", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Contains(t, response["error"], "Invalid country ID")

	mockRepo.AssertNotCalled(t, "GetByCountry")
}

func TestGetCompaniesByCountry_DatabaseError(t *testing.T) {
	app := setupTestApp()
	handler, mockRepo := setupTestHandler()

	mockRepo.On("GetByCountry", int64(1)).Return([]models.Company{}, errors.New("database error")).Once()

	app.Get("/companies/country/:countryId", handler.GetCompaniesByCountry)

	req := httptest.NewRequest("GET", "/companies/country/1", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	mockRepo.AssertExpectations(t)
}

func TestGetCompaniesByIndustry_Success(t *testing.T) {
	app := setupTestApp()
	handler, mockRepo := setupTestHandler()

	expectedCompanies := []models.Company{
		{
			ID:          1,
			Name:        "TechCorp",
			Code:        "TECH001",
			CountryID:   1,
			Description: "Technology company",
			Website:     "https://techcorp.com",
			Industry:    "Technology",
			Size:        "large",
			Founded:     2010,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          2,
			Name:        "DataCorp",
			Code:        "DATA001",
			CountryID:   1,
			Description: "Data analytics company",
			Website:     "https://datacorp.com",
			Industry:    "Technology",
			Size:        "medium",
			Founded:     2015,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	mockRepo.On("GetByIndustry", "Technology").Return(expectedCompanies, nil).Once()

	app.Get("/companies/industry/:industry", handler.GetCompaniesByIndustry)

	req := httptest.NewRequest("GET", "/companies/industry/Technology", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var responseCompanies []models.Company
	json.NewDecoder(resp.Body).Decode(&responseCompanies)
	assert.Len(t, responseCompanies, 2)
	assert.Equal(t, expectedCompanies[0].Name, responseCompanies[0].Name)
	assert.Equal(t, expectedCompanies[1].Name, responseCompanies[1].Name)

	mockRepo.AssertExpectations(t)
}

// TestGetCompaniesByIndustry_EmptyString tests with an actual empty string parameter using SQLite
func TestGetCompaniesByIndustry_EmptyString(t *testing.T) {
	// Setup test database
	testDB, err := test_helpers.SetupTestDB()
	if err != nil {
		t.Fatalf("Failed to setup test database: %v", err)
	}
	defer test_helpers.CleanupTestDB(testDB)

	// Create test continent and country first
	continent, err := test_helpers.CreateTestContinent(testDB, "Europe", "EU")
	if err != nil {
		t.Fatalf("Failed to create test continent: %v", err)
	}

	country, err := test_helpers.CreateTestCountry(testDB, "Germany", "DEU", continent.ID)
	if err != nil {
		t.Fatalf("Failed to create test country: %v", err)
	}

	// Create a test company (needed for database setup)
	_, err = test_helpers.CreateTestCompany(testDB, "TechCorp", "TECH001", country.ID)
	if err != nil {
		t.Fatalf("Failed to create test company: %v", err)
	}

	// Setup app and handler
	app := test_helpers.SetupTestApp()
	handler := NewCompanyHandler(testDB.CompanyRepo, logger.Get())

	app.Get("/companies/industry/:industry", handler.GetCompaniesByIndustry)

	// Test with a space character which should be treated as empty after trimming
	req := httptest.NewRequest("GET", "/companies/industry/%20", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Contains(t, response["error"], "Industry is required")
}

func TestGetCompaniesByIndustry_DatabaseError(t *testing.T) {
	app := setupTestApp()
	handler, mockRepo := setupTestHandler()

	mockRepo.On("GetByIndustry", "Technology").Return([]models.Company{}, errors.New("database error")).Once()

	app.Get("/companies/industry/:industry", handler.GetCompaniesByIndustry)

	req := httptest.NewRequest("GET", "/companies/industry/Technology", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	mockRepo.AssertExpectations(t)
}

func TestGetAllCompanies_Success(t *testing.T) {
	app := setupTestApp()
	handler, mockRepo := setupTestHandler()

	expectedCompanies := []models.Company{
		{
			ID:          1,
			Name:        "TechCorp",
			Code:        "TECH001",
			CountryID:   1,
			Description: "Technology company",
			Website:     "https://techcorp.com",
			Industry:    "Technology",
			Size:        "large",
			Founded:     2010,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          2,
			Name:        "DataCorp",
			Code:        "DATA001",
			CountryID:   1,
			Description: "Data analytics company",
			Website:     "https://datacorp.com",
			Industry:    "Technology",
			Size:        "medium",
			Founded:     2015,
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	mockRepo.On("GetAll").Return(expectedCompanies, nil).Once()

	app.Get("/companies", handler.GetAllCompanies)

	req := httptest.NewRequest("GET", "/companies", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var responseCompanies []models.Company
	json.NewDecoder(resp.Body).Decode(&responseCompanies)
	assert.Len(t, responseCompanies, 2)
	assert.Equal(t, expectedCompanies[0].Name, responseCompanies[0].Name)
	assert.Equal(t, expectedCompanies[1].Name, responseCompanies[1].Name)

	mockRepo.AssertExpectations(t)
}

func TestGetAllCompanies_DatabaseError(t *testing.T) {
	app := setupTestApp()
	handler, mockRepo := setupTestHandler()

	mockRepo.On("GetAll").Return([]models.Company{}, errors.New("database error")).Once()

	app.Get("/companies", handler.GetAllCompanies)

	req := httptest.NewRequest("GET", "/companies", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	mockRepo.AssertExpectations(t)
}

func TestUpdateCompany_Success(t *testing.T) {
	app := setupTestApp()
	handler, mockRepo := setupTestHandler()

	updateData := models.Company{
		Name:        "Updated TechCorp",
		Code:        "TECH001",
		CountryID:   1,
		Description: "Updated Technology company",
		Website:     "https://updated-techcorp.com",
		Industry:    "Technology",
		Size:        "large",
		Founded:     2010,
	}

	expectedCompany := updateData
	expectedCompany.ID = 1
	expectedCompany.CreatedAt = time.Now()
	expectedCompany.UpdatedAt = time.Now()

	mockRepo.On("Update", mock.AnythingOfType("*models.Company")).Return(&expectedCompany, nil).Once()

	app.Put("/companies/:id", handler.UpdateCompany)

	body, _ := json.Marshal(updateData)
	req := httptest.NewRequest("PUT", "/companies/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var responseCompany models.Company
	json.NewDecoder(resp.Body).Decode(&responseCompany)
	assert.Equal(t, expectedCompany.Name, responseCompany.Name)
	assert.Equal(t, expectedCompany.Description, responseCompany.Description)
	assert.Equal(t, expectedCompany.Website, responseCompany.Website)

	mockRepo.AssertExpectations(t)
}

func TestUpdateCompany_InvalidID(t *testing.T) {
	app := setupTestApp()
	handler, mockRepo := setupTestHandler()

	app.Put("/companies/:id", handler.UpdateCompany)

	updateData := models.Company{
		Name:        "Updated TechCorp",
		Code:        "TECH001",
		CountryID:   1,
		Description: "Updated Technology company",
		Size:        "large",
	}

	body, _ := json.Marshal(updateData)
	req := httptest.NewRequest("PUT", "/companies/invalid", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Contains(t, response["error"], "Invalid company ID")

	mockRepo.AssertNotCalled(t, "Update")
}

func TestUpdateCompany_InvalidJSON(t *testing.T) {
	app := setupTestApp()
	handler, mockRepo := setupTestHandler()

	app.Put("/companies/:id", handler.UpdateCompany)

	req := httptest.NewRequest("PUT", "/companies/1", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Contains(t, response["error"], "Invalid request body")

	mockRepo.AssertNotCalled(t, "Update")
}

func TestUpdateCompany_ValidationError(t *testing.T) {
	app := setupTestApp()
	handler, mockRepo := setupTestHandler()

	updateData := models.Company{
		Name:      "", // Invalid: name is required
		Code:      "TECH001",
		CountryID: 1,
		Size:      "large",
	}

	app.Put("/companies/:id", handler.UpdateCompany)

	body, _ := json.Marshal(updateData)
	req := httptest.NewRequest("PUT", "/companies/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Contains(t, response["error"], "Validation failed")

	mockRepo.AssertNotCalled(t, "Update")
}

func TestUpdateCompany_DatabaseError(t *testing.T) {
	app := setupTestApp()
	handler, mockRepo := setupTestHandler()

	updateData := models.Company{
		Name:        "Updated TechCorp",
		Code:        "TECH001",
		CountryID:   1,
		Description: "Updated Technology company",
		Website:     "https://updated-techcorp.com",
		Industry:    "Technology",
		Size:        "large",
		Founded:     2010,
	}

	mockRepo.On("Update", mock.AnythingOfType("*models.Company")).Return(nil, errors.New("database error")).Once()

	app.Put("/companies/:id", handler.UpdateCompany)

	body, _ := json.Marshal(updateData)
	req := httptest.NewRequest("PUT", "/companies/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	mockRepo.AssertExpectations(t)
}

func TestDeleteCompany_Success(t *testing.T) {
	app := setupTestApp()
	handler, mockRepo := setupTestHandler()

	mockRepo.On("Delete", int64(1)).Return(nil).Once()

	app.Delete("/companies/:id", handler.DeleteCompany)

	req := httptest.NewRequest("DELETE", "/companies/1", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusNoContent, resp.StatusCode)

	mockRepo.AssertExpectations(t)
}

func TestDeleteCompany_InvalidID(t *testing.T) {
	app := setupTestApp()
	handler, mockRepo := setupTestHandler()

	app.Delete("/companies/:id", handler.DeleteCompany)

	req := httptest.NewRequest("DELETE", "/companies/invalid", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Contains(t, response["error"], "Invalid company ID")

	mockRepo.AssertNotCalled(t, "Delete")
}

func TestDeleteCompany_DatabaseError(t *testing.T) {
	app := setupTestApp()
	handler, mockRepo := setupTestHandler()

	mockRepo.On("Delete", int64(1)).Return(errors.New("database error")).Once()

	app.Delete("/companies/:id", handler.DeleteCompany)

	req := httptest.NewRequest("DELETE", "/companies/1", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	mockRepo.AssertExpectations(t)
}
