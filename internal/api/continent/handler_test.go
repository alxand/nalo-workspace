package continent

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

// MockContinentRepository is a mock implementation of ContinentInterface
type MockContinentRepository struct {
	mock.Mock
}

func (m *MockContinentRepository) Create(continent *models.Continent) error {
	args := m.Called(continent)
	return args.Error(0)
}

func (m *MockContinentRepository) GetByID(id int64) (*models.Continent, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Continent), args.Error(1)
}

func (m *MockContinentRepository) GetByCode(code string) (*models.Continent, error) {
	args := m.Called(code)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Continent), args.Error(1)
}

func (m *MockContinentRepository) GetAll() ([]models.Continent, error) {
	args := m.Called()
	return args.Get(0).([]models.Continent), args.Error(1)
}

func (m *MockContinentRepository) Update(continent *models.Continent) (*models.Continent, error) {
	args := m.Called(continent)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.Continent), args.Error(1)
}

func (m *MockContinentRepository) Delete(id int64) error {
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

func setupTestHandler() (*ContinentHandler, *MockContinentRepository) {
	// Initialize validation for all tests
	validation.Init()

	mockRepo := new(MockContinentRepository)
	logger := zap.NewNop()
	handler := NewContinentHandler(mockRepo, logger)
	return handler, mockRepo
}

func TestCreateContinent_Success(t *testing.T) {
	// Initialize validation
	validation.Init()

	app := setupTestApp()
	handler, mockRepo := setupTestHandler()

	continent := models.Continent{
		Name:        "Europe",
		Code:        "EU",
		Description: "European continent",
	}

	expectedContinent := continent
	expectedContinent.ID = 1
	expectedContinent.CreatedAt = time.Now()
	expectedContinent.UpdatedAt = time.Now()

	mockRepo.On("Create", mock.AnythingOfType("*models.Continent")).Return(nil).Once()

	app.Post("/continents", handler.CreateContinent)

	body, _ := json.Marshal(continent)
	req := httptest.NewRequest("POST", "/continents", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusCreated, resp.StatusCode)

	var responseContinent models.Continent
	json.NewDecoder(resp.Body).Decode(&responseContinent)
	assert.Equal(t, continent.Name, responseContinent.Name)
	assert.Equal(t, continent.Code, responseContinent.Code)
	assert.Equal(t, continent.Description, responseContinent.Description)

	mockRepo.AssertExpectations(t)
}

func TestCreateContinent_InvalidJSON(t *testing.T) {
	app := setupTestApp()
	handler, mockRepo := setupTestHandler()

	app.Post("/continents", handler.CreateContinent)

	req := httptest.NewRequest("POST", "/continents", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Contains(t, response["error"], "Invalid request body")

	mockRepo.AssertNotCalled(t, "Create")
}

func TestCreateContinent_ValidationError(t *testing.T) {
	// Initialize validation
	validation.Init()

	app := setupTestApp()
	handler, mockRepo := setupTestHandler()

	continent := models.Continent{
		Name: "", // Invalid: name is required
		Code: "EU",
	}

	app.Post("/continents", handler.CreateContinent)

	body, _ := json.Marshal(continent)
	req := httptest.NewRequest("POST", "/continents", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Contains(t, response["error"], "Validation failed")

	mockRepo.AssertNotCalled(t, "Create")
}

func TestCreateContinent_DatabaseError(t *testing.T) {
	// Initialize validation
	validation.Init()

	app := setupTestApp()
	handler, mockRepo := setupTestHandler()

	continent := models.Continent{
		Name:        "Europe",
		Code:        "EU",
		Description: "European continent",
	}

	mockRepo.On("Create", mock.AnythingOfType("*models.Continent")).Return(errors.New("database error")).Once()

	app.Post("/continents", handler.CreateContinent)

	body, _ := json.Marshal(continent)
	req := httptest.NewRequest("POST", "/continents", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	mockRepo.AssertExpectations(t)
}

func TestGetContinent_Success(t *testing.T) {
	app := setupTestApp()
	handler, mockRepo := setupTestHandler()

	expectedContinent := &models.Continent{
		ID:          1,
		Name:        "Europe",
		Code:        "EU",
		Description: "European continent",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	mockRepo.On("GetByID", int64(1)).Return(expectedContinent, nil).Once()

	app.Get("/continents/:id", handler.GetContinent)

	req := httptest.NewRequest("GET", "/continents/1", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var responseContinent models.Continent
	json.NewDecoder(resp.Body).Decode(&responseContinent)
	assert.Equal(t, expectedContinent.ID, responseContinent.ID)
	assert.Equal(t, expectedContinent.Name, responseContinent.Name)
	assert.Equal(t, expectedContinent.Code, responseContinent.Code)

	mockRepo.AssertExpectations(t)
}

func TestGetContinent_InvalidID(t *testing.T) {
	app := setupTestApp()
	handler, mockRepo := setupTestHandler()

	app.Get("/continents/:id", handler.GetContinent)

	req := httptest.NewRequest("GET", "/continents/invalid", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Contains(t, response["error"], "Invalid continent ID")

	mockRepo.AssertNotCalled(t, "GetByID")
}

func TestGetContinent_NotFound(t *testing.T) {
	app := setupTestApp()
	handler, mockRepo := setupTestHandler()

	mockRepo.On("GetByID", int64(999)).Return(nil, errors.New("not found")).Once()

	app.Get("/continents/:id", handler.GetContinent)

	req := httptest.NewRequest("GET", "/continents/999", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

	mockRepo.AssertExpectations(t)
}

func TestGetContinentByCode_Success(t *testing.T) {
	app := setupTestApp()
	handler, mockRepo := setupTestHandler()

	expectedContinent := &models.Continent{
		ID:          1,
		Name:        "Europe",
		Code:        "EU",
		Description: "European continent",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	mockRepo.On("GetByCode", "EU").Return(expectedContinent, nil).Once()

	app.Get("/continents/code/:code", handler.GetContinentByCode)

	req := httptest.NewRequest("GET", "/continents/code/EU", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var responseContinent models.Continent
	json.NewDecoder(resp.Body).Decode(&responseContinent)
	assert.Equal(t, expectedContinent.Code, responseContinent.Code)

	mockRepo.AssertExpectations(t)
}

func TestGetContinentByCode_NotFound(t *testing.T) {
	app := setupTestApp()
	handler, mockRepo := setupTestHandler()

	mockRepo.On("GetByCode", "XX").Return(nil, errors.New("not found")).Once()

	app.Get("/continents/code/:code", handler.GetContinentByCode)

	req := httptest.NewRequest("GET", "/continents/code/XX", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

	mockRepo.AssertExpectations(t)
}

func TestGetAllContinents_Success(t *testing.T) {
	app := setupTestApp()
	handler, mockRepo := setupTestHandler()

	expectedContinents := []models.Continent{
		{
			ID:          1,
			Name:        "Europe",
			Code:        "EU",
			Description: "European continent",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
		{
			ID:          2,
			Name:        "Asia",
			Code:        "AS",
			Description: "Asian continent",
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
		},
	}

	mockRepo.On("GetAll").Return(expectedContinents, nil).Once()

	app.Get("/continents", handler.GetAllContinents)

	req := httptest.NewRequest("GET", "/continents", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var responseContinents []models.Continent
	json.NewDecoder(resp.Body).Decode(&responseContinents)
	assert.Len(t, responseContinents, 2)
	assert.Equal(t, expectedContinents[0].Name, responseContinents[0].Name)
	assert.Equal(t, expectedContinents[1].Name, responseContinents[1].Name)

	mockRepo.AssertExpectations(t)
}

func TestGetAllContinents_DatabaseError(t *testing.T) {
	app := setupTestApp()
	handler, mockRepo := setupTestHandler()

	mockRepo.On("GetAll").Return([]models.Continent{}, errors.New("database error")).Once()

	app.Get("/continents", handler.GetAllContinents)

	req := httptest.NewRequest("GET", "/continents", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	mockRepo.AssertExpectations(t)
}

func TestUpdateContinent_Success(t *testing.T) {
	// Initialize validation
	validation.Init()

	app := setupTestApp()
	handler, mockRepo := setupTestHandler()

	updateData := models.Continent{
		Name:        "Updated Europe",
		Code:        "EU",
		Description: "Updated European continent",
	}

	expectedContinent := updateData
	expectedContinent.ID = 1
	expectedContinent.CreatedAt = time.Now()
	expectedContinent.UpdatedAt = time.Now()

	mockRepo.On("Update", mock.AnythingOfType("*models.Continent")).Return(&expectedContinent, nil).Once()

	app.Put("/continents/:id", handler.UpdateContinent)

	body, _ := json.Marshal(updateData)
	req := httptest.NewRequest("PUT", "/continents/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var responseContinent models.Continent
	json.NewDecoder(resp.Body).Decode(&responseContinent)
	assert.Equal(t, expectedContinent.Name, responseContinent.Name)
	assert.Equal(t, expectedContinent.Description, responseContinent.Description)

	mockRepo.AssertExpectations(t)
}

func TestUpdateContinent_InvalidID(t *testing.T) {
	app := setupTestApp()
	handler, mockRepo := setupTestHandler()

	app.Put("/continents/:id", handler.UpdateContinent)

	updateData := models.Continent{
		Name:        "Updated Europe",
		Code:        "EU",
		Description: "Updated European continent",
	}

	body, _ := json.Marshal(updateData)
	req := httptest.NewRequest("PUT", "/continents/invalid", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Contains(t, response["error"], "Invalid continent ID")

	mockRepo.AssertNotCalled(t, "Update")
}

func TestUpdateContinent_InvalidJSON(t *testing.T) {
	app := setupTestApp()
	handler, mockRepo := setupTestHandler()

	app.Put("/continents/:id", handler.UpdateContinent)

	req := httptest.NewRequest("PUT", "/continents/1", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Contains(t, response["error"], "Invalid request body")

	mockRepo.AssertNotCalled(t, "Update")
}

func TestUpdateContinent_ValidationError(t *testing.T) {
	// Initialize validation
	validation.Init()

	app := setupTestApp()
	handler, mockRepo := setupTestHandler()

	updateData := models.Continent{
		Name: "", // Invalid: name is required
		Code: "EU",
	}

	app.Put("/continents/:id", handler.UpdateContinent)

	body, _ := json.Marshal(updateData)
	req := httptest.NewRequest("PUT", "/continents/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Contains(t, response["error"], "Validation failed")

	mockRepo.AssertNotCalled(t, "Update")
}

func TestUpdateContinent_DatabaseError(t *testing.T) {
	// Initialize validation
	validation.Init()

	app := setupTestApp()
	handler, mockRepo := setupTestHandler()

	updateData := models.Continent{
		Name:        "Updated Europe",
		Code:        "EU",
		Description: "Updated European continent",
	}

	mockRepo.On("Update", mock.AnythingOfType("*models.Continent")).Return(nil, errors.New("database error")).Once()

	app.Put("/continents/:id", handler.UpdateContinent)

	body, _ := json.Marshal(updateData)
	req := httptest.NewRequest("PUT", "/continents/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	mockRepo.AssertExpectations(t)
}

func TestDeleteContinent_Success(t *testing.T) {
	app := setupTestApp()
	handler, mockRepo := setupTestHandler()

	mockRepo.On("Delete", int64(1)).Return(nil).Once()

	app.Delete("/continents/:id", handler.DeleteContinent)

	req := httptest.NewRequest("DELETE", "/continents/1", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusNoContent, resp.StatusCode)

	mockRepo.AssertExpectations(t)
}

func TestDeleteContinent_InvalidID(t *testing.T) {
	app := setupTestApp()
	handler, mockRepo := setupTestHandler()

	app.Delete("/continents/:id", handler.DeleteContinent)

	req := httptest.NewRequest("DELETE", "/continents/invalid", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Contains(t, response["error"], "Invalid continent ID")

	mockRepo.AssertNotCalled(t, "Delete")
}

func TestDeleteContinent_DatabaseError(t *testing.T) {
	app := setupTestApp()
	handler, mockRepo := setupTestHandler()

	mockRepo.On("Delete", int64(1)).Return(errors.New("database error")).Once()

	app.Delete("/continents/:id", handler.DeleteContinent)

	req := httptest.NewRequest("DELETE", "/continents/1", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	mockRepo.AssertExpectations(t)
}
