package dailytask

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alxand/nalo-workspace/internal/domain/models"
	"github.com/alxand/nalo-workspace/internal/pkg/errors"
	"github.com/alxand/nalo-workspace/internal/pkg/validation"
	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"go.uber.org/zap"
)

// MockRepository is a mock implementation of the DailyTaskInterface
type MockRepository struct {
	mock.Mock
}

func (m *MockRepository) Create(task *models.DailyTask) error {
	args := m.Called(task)
	return args.Error(0)
}

func (m *MockRepository) GetByID(id int64) (*models.DailyTask, error) {
	args := m.Called(id)
	return args.Get(0).(*models.DailyTask), args.Error(1)
}

func (m *MockRepository) GetByDate(date string) ([]models.DailyTask, error) {
	args := m.Called(date)
	return args.Get(0).([]models.DailyTask), args.Error(1)
}

func (m *MockRepository) GetByDateAndUser(date string, userID int64) ([]models.DailyTask, error) {
	args := m.Called(date, userID)
	return args.Get(0).([]models.DailyTask), args.Error(1)
}

func (m *MockRepository) Update(task *models.DailyTask) (*models.DailyTask, error) {
	args := m.Called(task)
	return args.Get(0).(*models.DailyTask), args.Error(1)
}

func (m *MockRepository) Delete(id int64) error {
	args := m.Called(id)
	return args.Error(0)
}

func (m *MockRepository) List(limit, offset int) ([]models.DailyTask, error) {
	args := m.Called(limit, offset)
	return args.Get(0).([]models.DailyTask), args.Error(1)
}

// TestHelper provides common test utilities
type TestHelper struct {
	app    *fiber.App
	repo   *MockRepository
	logger *zap.Logger
}

func setupTest() *TestHelper {
	validation.Init() // Register custom validation functions
	app := fiber.New(fiber.Config{
		ErrorHandler: func(c *fiber.Ctx, err error) error {
			// Return error as JSON with status code
			code := fiber.StatusInternalServerError
			if e, ok := err.(*fiber.Error); ok {
				code = e.Code
			}
			return c.Status(code).JSON(fiber.Map{"error": err.Error()})
		},
	})
	repo := new(MockRepository)
	logger, _ := zap.NewDevelopment()

	handler := NewTDailyTaskHandler(repo, logger)

	// Setup routes with mock user middleware
	app.Post("/tasks", func(c *fiber.Ctx) error {
		// Add mock user to context
		mockUser := &models.User{
			ID:       1,
			Email:    "test@example.com",
			Username: "testuser",
			Role:     models.RoleUser,
		}
		c.Locals("user", mockUser)
		return handler.CreateDailyTask(c)
	})

	app.Get("/tasks/:date", func(c *fiber.Ctx) error {
		// Add mock user to context
		mockUser := &models.User{
			ID:       1,
			Email:    "test@example.com",
			Username: "testuser",
			Role:     models.RoleUser,
		}
		c.Locals("user", mockUser)
		return handler.GetTasksByDate(c)
	})

	app.Put("/tasks/:id", func(c *fiber.Ctx) error {
		// Add mock user to context
		mockUser := &models.User{
			ID:       1,
			Email:    "test@example.com",
			Username: "testuser",
			Role:     models.RoleUser,
		}
		c.Locals("user", mockUser)
		return handler.UpdateTask(c)
	})

	app.Delete("/tasks/:id", func(c *fiber.Ctx) error {
		// Add mock user to context
		mockUser := &models.User{
			ID:       1,
			Email:    "test@example.com",
			Username: "testuser",
			Role:     models.RoleUser,
		}
		c.Locals("user", mockUser)
		return handler.DeleteTask(c)
	})

	return &TestHelper{
		app:    app,
		repo:   repo,
		logger: logger,
	}
}

func TestCreateDailyTask_Success(t *testing.T) {
	helper := setupTest()

	task := models.DailyTask{
		Day:               "Monday",
		Date:              time.Now(),
		StartTime:         time.Now(),
		EndTime:           time.Now().Add(time.Hour),
		Status:            "pending",
		Score:             8,
		ProductivityScore: 7,
	}

	helper.repo.On("Create", mock.AnythingOfType("*models.DailyTask")).Return(nil)

	body, _ := json.Marshal(task)
	req := httptest.NewRequest("POST", "/tasks", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := helper.app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusCreated, resp.StatusCode)

	helper.repo.AssertExpectations(t)
}

func TestCreateDailyTask_InvalidBody(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("panic occurred: %v", r)
		}
	}()

	helper := setupTest()

	req := httptest.NewRequest("POST", "/tasks", bytes.NewReader([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")

	resp, err := helper.app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, 400, resp.StatusCode)

	// Assert that Create was NOT called
	helper.repo.AssertNotCalled(t, "Create", mock.Anything)
}

func TestCreateDailyTask_ValidationError(t *testing.T) {
	helper := setupTest()

	// Create task with invalid data (empty required fields)
	task := models.DailyTask{
		// Missing required fields
	}

	body, _ := json.Marshal(task)
	req := httptest.NewRequest("POST", "/tasks", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := helper.app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestCreateDailyTask_DatabaseError(t *testing.T) {
	helper := setupTest()

	task := models.DailyTask{
		Day:               "Monday",
		Date:              time.Now(),
		StartTime:         time.Now(),
		EndTime:           time.Now().Add(time.Hour),
		Status:            "pending",
		Score:             8,
		ProductivityScore: 7,
	}

	helper.repo.On("Create", mock.AnythingOfType("*models.DailyTask")).Return(errors.DatabaseError("connection failed", nil))

	body, _ := json.Marshal(task)
	req := httptest.NewRequest("POST", "/tasks", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := helper.app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	helper.repo.AssertExpectations(t)
}

func TestGetTasksByDate_Success(t *testing.T) {
	helper := setupTest()

	expectedTasks := []models.DailyTask{
		{
			ID:     1,
			Day:    "Monday",
			Date:   time.Now(),
			Status: "completed",
		},
	}

	helper.repo.On("GetByDateAndUser", "2024-01-15", int64(1)).Return(expectedTasks, nil)

	req := httptest.NewRequest("GET", "/tasks/2024-01-15", nil)
	resp, err := helper.app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	helper.repo.AssertExpectations(t)
}

func TestGetTasksByDate_MissingDate(t *testing.T) {
	helper := setupTest()

	req := httptest.NewRequest("GET", "/tasks/", nil)
	resp, err := helper.app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusMethodNotAllowed, resp.StatusCode) // Fiber returns 405 when route doesn't match
}

func TestGetTasksByDate_DatabaseError(t *testing.T) {
	helper := setupTest()

	helper.repo.On("GetByDateAndUser", "2024-01-15", int64(1)).Return([]models.DailyTask{}, errors.DatabaseError("connection failed", nil))

	req := httptest.NewRequest("GET", "/tasks/2024-01-15", nil)
	resp, err := helper.app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	helper.repo.AssertExpectations(t)
}

func TestUpdateTask_Success(t *testing.T) {
	helper := setupTest()

	task := models.DailyTask{
		Day:               "Monday",
		Date:              time.Now(),
		StartTime:         time.Now(),
		EndTime:           time.Now().Add(time.Hour),
		Status:            "completed",
		Score:             9,
		ProductivityScore: 8,
	}

	existingTask := models.DailyTask{
		ID:     1,
		UserID: 1, // Same user ID as mock user
		Day:    "Monday",
		Date:   time.Now(),
		Status: "pending",
	}

	updatedTask := task
	updatedTask.ID = 1
	updatedTask.UserID = 1

	helper.repo.On("GetByID", int64(1)).Return(&existingTask, nil)
	helper.repo.On("Update", mock.AnythingOfType("*models.DailyTask")).Return(&updatedTask, nil)

	body, _ := json.Marshal(task)
	req := httptest.NewRequest("PUT", "/tasks/1", bytes.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := helper.app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	helper.repo.AssertExpectations(t)
}

func TestUpdateTask_InvalidID(t *testing.T) {
	helper := setupTest()

	req := httptest.NewRequest("PUT", "/tasks/invalid", nil)
	resp, err := helper.app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestDeleteTask_Success(t *testing.T) {
	helper := setupTest()

	existingTask := models.DailyTask{
		ID:     1,
		UserID: 1, // Same user ID as mock user
		Day:    "Monday",
		Date:   time.Now(),
		Status: "completed",
	}

	helper.repo.On("GetByID", int64(1)).Return(&existingTask, nil)
	helper.repo.On("Delete", int64(1)).Return(nil)

	req := httptest.NewRequest("DELETE", "/tasks/1", nil)
	resp, err := helper.app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusNoContent, resp.StatusCode)

	helper.repo.AssertExpectations(t)
}

func TestDeleteTask_InvalidID(t *testing.T) {
	helper := setupTest()

	req := httptest.NewRequest("DELETE", "/tasks/invalid", nil)
	resp, err := helper.app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusBadRequest, resp.StatusCode)
}

func TestDeleteTask_DatabaseError(t *testing.T) {
	helper := setupTest()

	existingTask := models.DailyTask{
		ID:     1,
		UserID: 1, // Same user ID as mock user
		Day:    "Monday",
		Date:   time.Now(),
		Status: "completed",
	}

	helper.repo.On("GetByID", int64(1)).Return(&existingTask, nil)
	helper.repo.On("Delete", int64(1)).Return(errors.DatabaseError("connection failed", nil))

	req := httptest.NewRequest("DELETE", "/tasks/1", nil)
	resp, err := helper.app.Test(req)
	assert.NoError(t, err)
	assert.Equal(t, http.StatusInternalServerError, resp.StatusCode)

	helper.repo.AssertExpectations(t)
}
