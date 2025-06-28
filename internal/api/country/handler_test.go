package country

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alxand/nalo-workspace/internal/domain/models"
	"github.com/alxand/nalo-workspace/internal/pkg/middleware"
	"github.com/alxand/nalo-workspace/internal/pkg/validation"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// MockCountryRepository is a mock implementation of CountryInterface
type MockCountryRepository struct {
	mock.Mock
}

func (m *MockCountryRepository) Create(country *models.Country) error {
	args := m.Called(country)
	return args.Error(0)
}

func (m *MockCountryRepository) GetByID(id int64) (*models.Country, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Country), args.Error(1)
}

func (m *MockCountryRepository) GetByCode(code string) (*models.Country, error) {
	args := m.Called(code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Country), args.Error(1)
}

func (m *MockCountryRepository) GetByContinent(continentID int64) ([]models.Country, error) {
	args := m.Called(continentID)
	return args.Get(0).([]models.Country), args.Error(1)
}

func (m *MockCountryRepository) GetAll() ([]models.Country, error) {
	args := m.Called()
	return args.Get(0).([]models.Country), args.Error(1)
}

func (m *MockCountryRepository) Update(country *models.Country) (*models.Country, error) {
	args := m.Called(country)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Country), args.Error(1)
}

func (m *MockCountryRepository) Delete(id int64) error {
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

func setupTestHandler() (*CountryHandler, *MockCountryRepository) {
	// Initialize validation for all tests
	validation.Init()

	mockRepo := new(MockCountryRepository)
	logger := zap.NewNop()
	handler := NewCountryHandler(mockRepo, logger)
	return handler, mockRepo
}

func TestCreateCountry_Success(t *testing.T) {
	// Initialize validation
	validation.Init()

	app := setupTestApp()
	handler, mockRepo := setupTestHandler()

	country := models.Country{
		Name:        "Germany",
		Code:        "DEU",
		ContinentID: 1,
		Description: "Federal Republic of Germany",
	}

	expectedCountry := country
	expectedCountry.ID = 1
	expectedCountry.CreatedAt = time.Now()
	expectedCountry.UpdatedAt = time.Now()

	mockRepo.On("Create", mock.AnythingOfType("*models.Country")).Return(nil).Once()

	app.Post("/countries", handler.CreateCountry)

	body, _ := json.Marshal(country)
	req := httptest.NewRequest("POST", "/countries", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

	var responseCountry models.Country
	json.NewDecoder(resp.Body).Decode(&responseCountry)
	assert.Equal(t, country.Name, responseCountry.Name)
	assert.Equal(t, country.Code, responseCountry.Code)
	assert.Equal(t, country.ContinentID, responseCountry.ContinentID)

	mockRepo.AssertExpectations(t)
}

func TestCreateCountry_InvalidJSON(t *testing.T) {
	app := setupTestApp()
	handler, mockRepo := setupTestHandler()

	app.Post("/countries", handler.CreateCountry)

	req := httptest.NewRequest("POST", "/countries", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Contains(t, response["error"], "Invalid request body")

	mockRepo.AssertNotCalled(t, "Create")
}

func TestCreateCountry_ValidationError(t *testing.T) {
	// Initialize validation
	validation.Init()

	app := setupTestApp()
	handler, mockRepo := setupTestHandler()

	country := models.Country{
		Name:        "", // Invalid: name is required
		Code:        "DEU",
		ContinentID: 1,
	}

	app.Post("/countries", handler.CreateCountry)

	body, _ := json.Marshal(country)
	req := httptest.NewRequest("POST", "/countries", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Contains(t, response["error"], "Validation failed")

	mockRepo.AssertNotCalled(t, "Create")
}

func TestCreateCountry_DatabaseError(t *testing.T) {
	// Initialize validation
	validation.Init()

	app := setupTestApp()
	handler, mockRepo := setupTestHandler()

	country := models.Country{
		Name:        "Germany",
		Code:        "DEU",
		ContinentID: 1,
		Description: "Federal Republic of Germany",
	}

	mockRepo.On("Create", mock.AnythingOfType("*models.Country")).Return(errors.New("database error")).Once()

	app.Post("/countries", handler.CreateCountry)

	body, _ := json.Marshal(country)
	req := httptest.NewRequest("POST", "/countries", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	mockRepo.AssertExpectations(t)
}

func TestGetCountry_Success(t *testing.T) {
	app := setupTestApp()
	handler, mockRepo := setupTestHandler()

	expectedCountry := &models.Country{
		ID:          1,
		Name:        "Germany",
		Code:        "DEU",
		ContinentID: 1,
		Description: "Federal Republic of Germany",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	mockRepo.On("GetByID", int64(1)).Return(expectedCountry, nil).Once()

	app.Get("/countries/:id", handler.GetCountry)

	req := httptest.NewRequest("GET", "/countries/1", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var responseCountry models.Country
	json.NewDecoder(resp.Body).Decode(&responseCountry)
	assert.Equal(t, expectedCountry.ID, responseCountry.ID)
	assert.Equal(t, expectedCountry.Name, responseCountry.Name)
	assert.Equal(t, expectedCountry.Code, responseCountry.Code)

	mockRepo.AssertExpectations(t)
}

func TestGetCountry_InvalidID(t *testing.T) {
	app := setupTestApp()
	handler, mockRepo := setupTestHandler()

	app.Get("/countries/:id", handler.GetCountry)

	req := httptest.NewRequest("GET", "/countries/invalid", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Contains(t, response["error"], "Invalid country ID")

	mockRepo.AssertNotCalled(t, "GetByID")
}

func TestGetCountry_NotFound(t *testing.T) {
	app := setupTestApp()
	handler, mockRepo := setupTestHandler()

	mockRepo.On("GetByID", int64(999)).Return(nil, errors.New("not found")).Once()

	app.Get("/countries/:id", handler.GetCountry)

	req := httptest.NewRequest("GET", "/countries/999", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

	mockRepo.AssertExpectations(t)
}

func TestGetCountryByCode_Success(t *testing.T) {
	app := setupTestApp()
	handler, mockRepo := setupTestHandler()

	expectedCountry := &models.Country{
		ID:          1,
		Name:        "Germany",
		Code:        "DEU",
		ContinentID: 1,
		Description: "Federal Republic of Germany",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	mockRepo.On("GetByCode", "DEU").Return(expectedCountry, nil).Once()

	app.Get("/countries/code/:code", handler.GetCountryByCode)

	req := httptest.NewRequest("GET", "/countries/code/DEU", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var responseCountry models.Country
	json.NewDecoder(resp.Body).Decode(&responseCountry)
	assert.Equal(t, expectedCountry.Code, responseCountry.Code)

	mockRepo.AssertExpectations(t)
}

func TestGetCountryByCode_NotFound(t *testing.T) {
	app := setupTestApp()
	handler, mockRepo := setupTestHandler()

	mockRepo.On("GetByCode", "XXX").Return(nil, errors.New("not found")).Once()

	app.Get("/countries/code/:code", handler.GetCountryByCode)

	req := httptest.NewRequest("GET", "/countries/code/XXX", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

	mockRepo.AssertExpectations(t)
}

func TestGetCountriesByContinent_Success(t *testing.T) {
	app := setupTestApp()
	handler, mockRepo := setupTestHandler()

	expectedCountries := []models.Country{
		{
			ID:          1,
			Name:        "Germany",
			Code:        "DEU",
			ContinentID: 1,
			Description: "Federal Republic of Germany",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          2,
			Name:        "France",
			Code:        "FRA",
			ContinentID: 1,
			Description: "French Republic",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	mockRepo.On("GetByContinent", int64(1)).Return(expectedCountries, nil).Once()

	app.Get("/countries/continent/:continentId", handler.GetCountriesByContinent)

	req := httptest.NewRequest("GET", "/countries/continent/1", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var responseCountries []models.Country
	json.NewDecoder(resp.Body).Decode(&responseCountries)
	assert.Len(t, responseCountries, 2)
	assert.Equal(t, expectedCountries[0].Name, responseCountries[0].Name)
	assert.Equal(t, expectedCountries[1].Name, responseCountries[1].Name)

	mockRepo.AssertExpectations(t)
}

func TestGetCountriesByContinent_InvalidContinentID(t *testing.T) {
	app := setupTestApp()
	handler, mockRepo := setupTestHandler()

	app.Get("/countries/continent/:continentId", handler.GetCountriesByContinent)

	req := httptest.NewRequest("GET", "/countries/continent/invalid", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Contains(t, response["error"], "Invalid continent ID")

	mockRepo.AssertNotCalled(t, "GetByContinent")
}

func TestGetCountriesByContinent_DatabaseError(t *testing.T) {
	app := setupTestApp()
	handler, mockRepo := setupTestHandler()

	mockRepo.On("GetByContinent", int64(1)).Return([]models.Country{}, errors.New("database error")).Once()

	app.Get("/countries/continent/:continentId", handler.GetCountriesByContinent)

	req := httptest.NewRequest("GET", "/countries/continent/1", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	mockRepo.AssertExpectations(t)
}

func TestGetAllCountries_Success(t *testing.T) {
	app := setupTestApp()
	handler, mockRepo := setupTestHandler()

	expectedCountries := []models.Country{
		{
			ID:          1,
			Name:        "Germany",
			Code:        "DEU",
			ContinentID: 1,
			Description: "Federal Republic of Germany",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          2,
			Name:        "France",
			Code:        "FRA",
			ContinentID: 1,
			Description: "French Republic",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	mockRepo.On("GetAll").Return(expectedCountries, nil).Once()

	app.Get("/countries", handler.GetAllCountries)

	req := httptest.NewRequest("GET", "/countries", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var responseCountries []models.Country
	json.NewDecoder(resp.Body).Decode(&responseCountries)
	assert.Len(t, responseCountries, 2)
	assert.Equal(t, expectedCountries[0].Name, responseCountries[0].Name)
	assert.Equal(t, expectedCountries[1].Name, responseCountries[1].Name)

	mockRepo.AssertExpectations(t)
}

func TestGetAllCountries_DatabaseError(t *testing.T) {
	app := setupTestApp()
	handler, mockRepo := setupTestHandler()

	mockRepo.On("GetAll").Return([]models.Country{}, errors.New("database error")).Once()

	app.Get("/countries", handler.GetAllCountries)

	req := httptest.NewRequest("GET", "/countries", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	mockRepo.AssertExpectations(t)
}

func TestUpdateCountry_Success(t *testing.T) {
	// Initialize validation
	validation.Init()

	app := setupTestApp()
	handler, mockRepo := setupTestHandler()

	updateData := models.Country{
		Name:        "Updated Germany",
		Code:        "DEU",
		ContinentID: 1,
		Description: "Updated Federal Republic of Germany",
	}

	expectedCountry := updateData
	expectedCountry.ID = 1
	expectedCountry.CreatedAt = time.Now()
	expectedCountry.UpdatedAt = time.Now()

	mockRepo.On("Update", mock.AnythingOfType("*models.Country")).Return(&expectedCountry, nil).Once()

	app.Put("/countries/:id", handler.UpdateCountry)

	body, _ := json.Marshal(updateData)
	req := httptest.NewRequest("PUT", "/countries/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var responseCountry models.Country
	json.NewDecoder(resp.Body).Decode(&responseCountry)
	assert.Equal(t, expectedCountry.Name, responseCountry.Name)
	assert.Equal(t, expectedCountry.Description, responseCountry.Description)

	mockRepo.AssertExpectations(t)
}

func TestUpdateCountry_InvalidID(t *testing.T) {
	app := setupTestApp()
	handler, mockRepo := setupTestHandler()

	app.Put("/countries/:id", handler.UpdateCountry)

	updateData := models.Country{
		Name:        "Updated Germany",
		Code:        "DEU",
		ContinentID: 1,
		Description: "Updated Federal Republic of Germany",
	}

	body, _ := json.Marshal(updateData)
	req := httptest.NewRequest("PUT", "/countries/invalid", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Contains(t, response["error"], "Invalid country ID")

	mockRepo.AssertNotCalled(t, "Update")
}

func TestUpdateCountry_InvalidJSON(t *testing.T) {
	app := setupTestApp()
	handler, mockRepo := setupTestHandler()

	app.Put("/countries/:id", handler.UpdateCountry)

	req := httptest.NewRequest("PUT", "/countries/1", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Contains(t, response["error"], "Invalid request body")

	mockRepo.AssertNotCalled(t, "Update")
}

func TestUpdateCountry_ValidationError(t *testing.T) {
	// Initialize validation
	validation.Init()

	app := setupTestApp()
	handler, mockRepo := setupTestHandler()

	updateData := models.Country{
		Name:        "", // Invalid: name is required
		Code:        "DEU",
		ContinentID: 1,
	}

	app.Put("/countries/:id", handler.UpdateCountry)

	body, _ := json.Marshal(updateData)
	req := httptest.NewRequest("PUT", "/countries/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Contains(t, response["error"], "Validation failed")

	mockRepo.AssertNotCalled(t, "Update")
}

func TestUpdateCountry_DatabaseError(t *testing.T) {
	// Initialize validation
	validation.Init()

	app := setupTestApp()
	handler, mockRepo := setupTestHandler()

	updateData := models.Country{
		Name:        "Updated Germany",
		Code:        "DEU",
		ContinentID: 1,
		Description: "Updated Federal Republic of Germany",
	}

	mockRepo.On("Update", mock.AnythingOfType("*models.Country")).Return(nil, errors.New("database error")).Once()

	app.Put("/countries/:id", handler.UpdateCountry)

	body, _ := json.Marshal(updateData)
	req := httptest.NewRequest("PUT", "/countries/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	mockRepo.AssertExpectations(t)
}

func TestDeleteCountry_Success(t *testing.T) {
	app := setupTestApp()
	handler, mockRepo := setupTestHandler()

	mockRepo.On("Delete", int64(1)).Return(nil).Once()

	app.Delete("/countries/:id", handler.DeleteCountry)

	req := httptest.NewRequest("DELETE", "/countries/1", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusNoContent, resp.StatusCode)

	mockRepo.AssertExpectations(t)
}

func TestDeleteCountry_InvalidID(t *testing.T) {
	app := setupTestApp()
	handler, mockRepo := setupTestHandler()

	app.Delete("/countries/:id", handler.DeleteCountry)

	req := httptest.NewRequest("DELETE", "/countries/invalid", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Contains(t, response["error"], "Invalid country ID")

	mockRepo.AssertNotCalled(t, "Delete")
}

func TestDeleteCountry_DatabaseError(t *testing.T) {
	app := setupTestApp()
	handler, mockRepo := setupTestHandler()

	mockRepo.On("Delete", int64(1)).Return(errors.New("database error")).Once()

	app.Delete("/countries/:id", handler.DeleteCountry)

	req := httptest.NewRequest("DELETE", "/countries/1", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	mockRepo.AssertExpectations(t)
}
