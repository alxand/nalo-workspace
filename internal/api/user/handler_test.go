package user

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"

	"github.com/alxand/nalo-workspace/internal/domain/models"
	"github.com/alxand/nalo-workspace/internal/pkg/middleware"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// MockUserRepository is a mock implementation of UserInterface
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) Create(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) GetByID(id int64) (*models.User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetByUsername(username string) (*models.User, error) {
	args := m.Called(username)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetByCountry(countryID int64) ([]models.User, error) {
	args := m.Called(countryID)
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockUserRepository) GetByCompany(companyID int64) ([]models.User, error) {
	args := m.Called(companyID)
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockUserRepository) Update(user *models.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) Delete(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserRepository) List(limit, offset int) ([]models.User, error) {
	args := m.Called(limit, offset)
	return args.Get(0).([]models.User), args.Error(1)
}

func (m *MockUserRepository) UpdateLastLogin(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockUserRepository) ExistsByEmail(email string) (bool, error) {
	args := m.Called(email)
	return args.Bool(0), args.Error(1)
}

func (m *MockUserRepository) ExistsByUsername(username string) (bool, error) {
	args := m.Called(username)
	return args.Bool(0), args.Error(1)
}

func setupTestApp() *fiber.App {
	logger := zap.NewNop()
	app := fiber.New(fiber.Config{
		ErrorHandler: middleware.ErrorHandler(logger),
	})
	return app
}

func setupTestHandler() (*UserHandler, *MockUserRepository) {
	mockRepo := new(MockUserRepository)
	logger := zap.NewNop()
	handler := NewUserHandler(mockRepo, logger)
	return handler, mockRepo
}

func TestListUsers_Success(t *testing.T) {
	app := setupTestApp()
	handler, mockRepo := setupTestHandler()

	expectedUsers := []models.User{
		{ID: 1, Email: "user1@example.com", Username: "user1"},
		{ID: 2, Email: "user2@example.com", Username: "user2"},
	}
	mockRepo.On("List", 10, 0).Return(expectedUsers, nil).Once()

	app.Get("/admin/users", handler.ListUsers)

	req := httptest.NewRequest("GET", "/admin/users", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var responseUsers []models.User
	json.NewDecoder(resp.Body).Decode(&responseUsers)
	assert.Len(t, responseUsers, 2)
	assert.Equal(t, expectedUsers[0].Email, responseUsers[0].Email)
	assert.Equal(t, expectedUsers[1].Email, responseUsers[1].Email)

	mockRepo.AssertExpectations(t)
}

func TestListUsers_WithPagination(t *testing.T) {
	app := setupTestApp()
	handler, mockRepo := setupTestHandler()

	expectedUsers := []models.User{
		{ID: 1, Email: "user1@example.com", Username: "user1"},
	}

	mockRepo.On("List", 5, 10).Return(expectedUsers, nil).Once()

	app.Get("/admin/users", handler.ListUsers)

	req := httptest.NewRequest("GET", "/admin/users?limit=5&offset=10", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var responseUsers []models.User
	json.NewDecoder(resp.Body).Decode(&responseUsers)
	assert.Len(t, responseUsers, 1)

	mockRepo.AssertExpectations(t)
}

func TestListUsers_InvalidPagination(t *testing.T) {
	app := setupTestApp()
	handler, mockRepo := setupTestHandler()

	expectedUsers := []models.User{
		{ID: 1, Email: "user1@example.com", Username: "user1"},
	}

	// Should use default values for invalid pagination
	mockRepo.On("List", 10, 0).Return(expectedUsers, nil).Once()

	app.Get("/admin/users", handler.ListUsers)

	req := httptest.NewRequest("GET", "/admin/users?limit=invalid&offset=invalid", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	mockRepo.AssertExpectations(t)
}

func TestListUsers_DatabaseError(t *testing.T) {
	app := setupTestApp()
	handler, mockRepo := setupTestHandler()

	mockRepo.On("List", 10, 0).Return([]models.User{}, errors.New("database error")).Once()

	app.Get("/admin/users", handler.ListUsers)

	req := httptest.NewRequest("GET", "/admin/users", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Contains(t, response["error"], "Failed to list users")

	mockRepo.AssertExpectations(t)
}

func TestGetUser_Success(t *testing.T) {
	app := setupTestApp()
	handler, mockRepo := setupTestHandler()

	expectedUser := &models.User{
		ID:        1,
		Email:     "user@example.com",
		Username:  "testuser",
		FirstName: "John",
		LastName:  "Doe",
		Role:      models.RoleUser,
	}

	mockRepo.On("GetByID", int64(1)).Return(expectedUser, nil).Once()

	app.Get("/admin/users/:id", handler.GetUser)

	req := httptest.NewRequest("GET", "/admin/users/1", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var responseUser models.User
	json.NewDecoder(resp.Body).Decode(&responseUser)
	assert.Equal(t, expectedUser.Email, responseUser.Email)
	assert.Equal(t, expectedUser.Username, responseUser.Username)

	mockRepo.AssertExpectations(t)
}

func TestGetUser_InvalidID(t *testing.T) {
	app := setupTestApp()
	handler, mockRepo := setupTestHandler()

	app.Get("/admin/users/:id", handler.GetUser)

	req := httptest.NewRequest("GET", "/admin/users/invalid", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Contains(t, response["error"], "Invalid user ID")

	mockRepo.AssertNotCalled(t, "GetByID")
}

func TestGetUser_NotFound(t *testing.T) {
	app := setupTestApp()
	handler, mockRepo := setupTestHandler()

	mockRepo.On("GetByID", int64(999)).Return(nil, errors.New("not found")).Once()

	app.Get("/admin/users/:id", handler.GetUser)

	req := httptest.NewRequest("GET", "/admin/users/999", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Contains(t, response["error"], "User not found")

	mockRepo.AssertExpectations(t)
}

func TestUpdateUser_Success(t *testing.T) {
	app := setupTestApp()
	handler, mockRepo := setupTestHandler()

	existingUser := &models.User{
		ID:        1,
		Email:     "user@example.com",
		Username:  "testuser",
		FirstName: "John",
		LastName:  "Doe",
	}

	updatedUser := models.User{
		FirstName: "Jane",
		LastName:  "Smith",
		Email:     "jane@example.com",
	}

	mockRepo.On("GetByID", int64(1)).Return(existingUser, nil).Once()
	mockRepo.On("Update", mock.AnythingOfType("*models.User")).Return(nil).Once()

	app.Put("/admin/users/:id", handler.UpdateUser)

	body, _ := json.Marshal(updatedUser)
	req := httptest.NewRequest("PUT", "/admin/users/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusOK, resp.StatusCode)

	var responseUser models.User
	json.NewDecoder(resp.Body).Decode(&responseUser)
	assert.Equal(t, int64(1), responseUser.ID)
	assert.Equal(t, updatedUser.FirstName, responseUser.FirstName)
	assert.Equal(t, updatedUser.LastName, responseUser.LastName)

	mockRepo.AssertExpectations(t)
}

func TestUpdateUser_InvalidID(t *testing.T) {
	app := setupTestApp()
	handler, mockRepo := setupTestHandler()

	app.Put("/admin/users/:id", handler.UpdateUser)

	body, _ := json.Marshal(models.User{})
	req := httptest.NewRequest("PUT", "/admin/users/invalid", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Contains(t, response["error"], "Invalid user ID")

	mockRepo.AssertNotCalled(t, "GetByID")
}

func TestUpdateUser_InvalidJSON(t *testing.T) {
	app := setupTestApp()
	handler, mockRepo := setupTestHandler()

	// Add mock expectation for GetByID since the handler calls it first
	mockRepo.On("GetByID", int64(1)).Return(&models.User{ID: 1}, nil).Once()

	app.Put("/admin/users/:id", handler.UpdateUser)

	req := httptest.NewRequest("PUT", "/admin/users/1", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Contains(t, response["error"], "Invalid request body")

	mockRepo.AssertExpectations(t)
}

func TestUpdateUser_UserNotFound(t *testing.T) {
	app := setupTestApp()
	handler, mockRepo := setupTestHandler()

	mockRepo.On("GetByID", int64(999)).Return(nil, errors.New("not found")).Once()

	app.Put("/admin/users/:id", handler.UpdateUser)

	body, _ := json.Marshal(models.User{})
	req := httptest.NewRequest("PUT", "/admin/users/999", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Contains(t, response["error"], "User not found")

	mockRepo.AssertExpectations(t)
}

func TestUpdateUser_DatabaseError(t *testing.T) {
	app := setupTestApp()
	handler, mockRepo := setupTestHandler()

	existingUser := &models.User{
		ID:       1,
		Email:    "user@example.com",
		Username: "testuser",
	}

	updatedUser := models.User{
		FirstName: "Jane",
		LastName:  "Smith",
	}

	mockRepo.On("GetByID", int64(1)).Return(existingUser, nil).Once()
	mockRepo.On("Update", mock.AnythingOfType("*models.User")).Return(errors.New("database error")).Once()

	app.Put("/admin/users/:id", handler.UpdateUser)

	body, _ := json.Marshal(updatedUser)
	req := httptest.NewRequest("PUT", "/admin/users/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Contains(t, response["error"], "Failed to update user")

	mockRepo.AssertExpectations(t)
}

func TestDeleteUser_Success(t *testing.T) {
	app := setupTestApp()
	handler, mockRepo := setupTestHandler()

	existingUser := &models.User{
		ID:       1,
		Email:    "user@example.com",
		Username: "testuser",
	}

	mockRepo.On("GetByID", int64(1)).Return(existingUser, nil).Once()
	mockRepo.On("Delete", int64(1)).Return(nil).Once()

	app.Delete("/admin/users/:id", handler.DeleteUser)

	req := httptest.NewRequest("DELETE", "/admin/users/1", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusNoContent, resp.StatusCode)

	mockRepo.AssertExpectations(t)
}

func TestDeleteUser_InvalidID(t *testing.T) {
	app := setupTestApp()
	handler, mockRepo := setupTestHandler()

	app.Delete("/admin/users/:id", handler.DeleteUser)

	req := httptest.NewRequest("DELETE", "/admin/users/invalid", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusBadRequest, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Contains(t, response["error"], "Invalid user ID")

	mockRepo.AssertNotCalled(t, "Delete")
}

func TestDeleteUser_UserNotFound(t *testing.T) {
	app := setupTestApp()
	handler, mockRepo := setupTestHandler()

	mockRepo.On("GetByID", int64(999)).Return(nil, errors.New("not found")).Once()

	app.Delete("/admin/users/:id", handler.DeleteUser)

	req := httptest.NewRequest("DELETE", "/admin/users/999", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusNotFound, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Contains(t, response["error"], "User not found")

	mockRepo.AssertExpectations(t)
}

func TestDeleteUser_DatabaseError(t *testing.T) {
	app := setupTestApp()
	handler, mockRepo := setupTestHandler()

	existingUser := &models.User{
		ID:       1,
		Email:    "user@example.com",
		Username: "testuser",
	}

	mockRepo.On("GetByID", int64(1)).Return(existingUser, nil).Once()
	mockRepo.On("Delete", int64(1)).Return(errors.New("database error")).Once()

	app.Delete("/admin/users/:id", handler.DeleteUser)

	req := httptest.NewRequest("DELETE", "/admin/users/1", nil)
	resp, _ := app.Test(req)

	assert.Equal(t, fiber.StatusInternalServerError, resp.StatusCode)

	var response map[string]interface{}
	json.NewDecoder(resp.Body).Decode(&response)
	assert.Contains(t, response["error"], "Failed to delete user")

	mockRepo.AssertExpectations(t)
}
