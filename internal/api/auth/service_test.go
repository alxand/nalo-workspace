package auth

import (
	"testing"
	"time"

	"github.com/alxand/nalo-workspace/internal/config"
	"github.com/alxand/nalo-workspace/internal/domain/models"
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
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetByEmail(email string) (*models.User, error) {
	args := m.Called(email)
	return args.Get(0).(*models.User), args.Error(1)
}

func (m *MockUserRepository) GetByUsername(username string) (*models.User, error) {
	args := m.Called(username)
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

func TestAuthService_Register(t *testing.T) {
	// Setup
	logger, _ := zap.NewDevelopment()
	jwtConfig := config.JWTConfig{
		Secret:     "test-secret",
		Expiration: 24 * time.Hour,
	}
	mockRepo := new(MockUserRepository)
	service := NewService(jwtConfig, mockRepo, logger)

	req := &RegisterRequest{
		Email:     "test@example.com",
		Username:  "testuser",
		Password:  "password123",
		FirstName: "John",
		LastName:  "Doe",
		Role:      models.RoleUser,
	}

	// Mock expectations
	mockRepo.On("ExistsByEmail", "test@example.com").Return(false, nil)
	mockRepo.On("ExistsByUsername", "testuser").Return(false, nil)
	mockRepo.On("Create", mock.AnythingOfType("*models.User")).Return(nil)

	// Test
	user, err := service.Register(req)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, user)
	assert.Equal(t, "test@example.com", user.Email)
	assert.Equal(t, "testuser", user.Username)
	assert.Equal(t, "John", user.FirstName)
	assert.Equal(t, "Doe", user.LastName)
	assert.Equal(t, models.RoleUser, user.Role)
	assert.True(t, user.IsActive)

	mockRepo.AssertExpectations(t)
}

func TestAuthService_Login(t *testing.T) {
	// Setup
	logger, _ := zap.NewDevelopment()
	jwtConfig := config.JWTConfig{
		Secret:     "test-secret",
		Expiration: 24 * time.Hour,
	}
	mockRepo := new(MockUserRepository)
	service := NewService(jwtConfig, mockRepo, logger)

	// Create a test user with hashed password
	testUser := &models.User{
		ID:        1,
		Email:     "test@example.com",
		Username:  "testuser",
		Password:  "$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi", // "password"
		FirstName: "John",
		LastName:  "Doe",
		Role:      models.RoleUser,
		IsActive:  true,
	}

	req := &LoginRequest{
		Email:    "test@example.com",
		Password: "password",
	}

	// Mock expectations
	mockRepo.On("GetByEmail", "test@example.com").Return(testUser, nil)
	mockRepo.On("UpdateLastLogin", int64(1)).Return(nil)

	// Test
	response, err := service.Login(req)

	// Assertions
	assert.NoError(t, err)
	assert.NotNil(t, response)
	assert.NotEmpty(t, response.Token)
	assert.Equal(t, "Bearer", response.Type)
	assert.Equal(t, testUser, response.User)
	assert.True(t, response.ExpiresAt.After(time.Now()))

	mockRepo.AssertExpectations(t)
}

func TestAuthService_GenerateJWT(t *testing.T) {
	// Setup
	logger, _ := zap.NewDevelopment()
	jwtConfig := config.JWTConfig{
		Secret:     "test-secret",
		Expiration: 24 * time.Hour,
	}
	mockRepo := new(MockUserRepository)
	service := NewService(jwtConfig, mockRepo, logger)

	user := &models.User{
		ID:       1,
		Email:    "test@example.com",
		Username: "testuser",
		Role:     models.RoleAdmin,
	}

	// Test
	token, expiresAt, err := service.GenerateJWT(user)

	// Assertions
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.True(t, expiresAt.After(time.Now()))

	// Validate the token
	claims, err := service.ValidateJWT(token)
	assert.NoError(t, err)
	assert.NotNil(t, claims)

	// Check claims
	userID, err := service.ExtractUserID(claims)
	assert.NoError(t, err)
	assert.Equal(t, int64(1), userID)

	role, err := service.ExtractUserRole(claims)
	assert.NoError(t, err)
	assert.Equal(t, models.RoleAdmin, role)
}

func TestAuthService_ValidateJWT(t *testing.T) {
	// Setup
	logger, _ := zap.NewDevelopment()
	jwtConfig := config.JWTConfig{
		Secret:     "test-secret",
		Expiration: 24 * time.Hour,
	}
	mockRepo := new(MockUserRepository)
	service := NewService(jwtConfig, mockRepo, logger)

	user := &models.User{
		ID:       1,
		Email:    "test@example.com",
		Username: "testuser",
		Role:     models.RoleUser,
	}

	// Generate a valid token
	token, _, err := service.GenerateJWT(user)
	assert.NoError(t, err)

	// Test valid token
	claims, err := service.ValidateJWT(token)
	assert.NoError(t, err)
	assert.NotNil(t, claims)

	// Test invalid token
	_, err = service.ValidateJWT("invalid-token")
	assert.Error(t, err)
}
