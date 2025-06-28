package auth

import (
	"github.com/alxand/nalo-workspace/internal/domain/models"
	"github.com/alxand/nalo-workspace/internal/pkg/errors"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type AuthHandler struct {
	authService *Service
	logger      *zap.Logger
	validate    *validator.Validate
}

func NewAuthHandler(authService *Service, logger *zap.Logger) *AuthHandler {
	return &AuthHandler{
		authService: authService,
		logger:      logger,
		validate:    validator.New(),
	}
}

// Register godoc
// @Summary Register a new user
// @Tags auth
// @Accept json
// @Produce json
// @Param user body RegisterRequest true "User registration data"
// @Success 201 {object} models.User
// @Failure 400 {object} map[string]string
// @Failure 409 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *fiber.Ctx) error {
	var req RegisterRequest

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		h.logger.Error("Failed to parse register request", zap.Error(err))
		return errors.BadRequest("Invalid request body", err)
	}

	// Validate request
	if err := h.validate.Struct(req); err != nil {
		h.logger.Error("Register validation failed", zap.Error(err))
		return errors.ValidationError("Validation failed", err)
	}

	// Register user
	user, err := h.authService.Register(&req)
	if err != nil {
		h.logger.Error("Failed to register user", zap.String("email", req.Email), zap.Error(err))
		if err.Error() == "email already exists" || err.Error() == "username already exists" {
			return errors.BadRequest(err.Error(), err)
		}
		return errors.InternalServerError("Failed to register user", err)
	}

	h.logger.Info("User registered successfully", zap.String("email", req.Email), zap.Int64("user_id", user.ID))
	return c.Status(fiber.StatusCreated).JSON(user)
}

// Login godoc
// @Summary Login user
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body LoginRequest true "Login credentials"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} map[string]string
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *fiber.Ctx) error {
	var req LoginRequest

	// Parse request body
	if err := c.BodyParser(&req); err != nil {
		h.logger.Error("Failed to parse login request", zap.Error(err))
		return errors.BadRequest("Invalid request body", err)
	}

	// Validate request
	if err := h.validate.Struct(req); err != nil {
		h.logger.Error("Login validation failed", zap.Error(err))
		return errors.ValidationError("Validation failed", err)
	}

	// Login user
	response, err := h.authService.Login(&req)
	if err != nil {
		h.logger.Error("Failed to login user", zap.String("email", req.Email), zap.Error(err))
		if err.Error() == "invalid credentials" || err.Error() == "account is deactivated" {
			return errors.Unauthorized(err.Error(), err)
		}
		return errors.InternalServerError("Failed to login user", err)
	}

	h.logger.Info("User logged in successfully", zap.String("email", req.Email), zap.Int64("user_id", response.User.ID))
	return c.JSON(response)
}

// Profile godoc
// @Summary Get current user profile
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.User
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/profile [get]
func (h *AuthHandler) Profile(c *fiber.Ctx) error {
	// Get user from context (set by JWT middleware)
	user := c.Locals("user").(*models.User)
	if user == nil {
		return errors.Unauthorized("User not found in context", nil)
	}

	return c.JSON(user)
}

// RefreshToken godoc
// @Summary Refresh JWT token
// @Tags auth
// @Produce json
// @Security BearerAuth
// @Success 200 {object} LoginResponse
// @Failure 401 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /auth/refresh [post]
func (h *AuthHandler) RefreshToken(c *fiber.Ctx) error {
	// Get user from context (set by JWT middleware)
	user := c.Locals("user").(*models.User)
	if user == nil {
		return errors.Unauthorized("User not found in context", nil)
	}

	// Generate new token
	token, expiresAt, err := h.authService.GenerateJWT(user)
	if err != nil {
		h.logger.Error("Failed to generate refresh token", zap.Int64("user_id", user.ID), zap.Error(err))
		return errors.InternalServerError("Failed to refresh token", err)
	}

	response := &LoginResponse{
		Token:     token,
		Type:      "Bearer",
		User:      user,
		ExpiresAt: expiresAt,
	}

	h.logger.Info("Token refreshed successfully", zap.Int64("user_id", user.ID))
	return c.JSON(response)
}
