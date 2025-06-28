package middleware

import (
	"strings"

	"github.com/alxand/nalo-workspace/internal/api/auth"
	"github.com/alxand/nalo-workspace/internal/domain/models"
	"github.com/alxand/nalo-workspace/internal/pkg/errors"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// ErrorHandler handles application errors
func ErrorHandler(logger *zap.Logger) fiber.ErrorHandler {
	return func(c *fiber.Ctx, err error) error {
		// Check if it's our custom AppError
		if appErr, ok := err.(*errors.AppError); ok {
			logger.Error("Application error",
				zap.String("path", c.Path()),
				zap.String("method", c.Method()),
				zap.Int("status", appErr.Code),
				zap.String("message", appErr.Message),
				zap.Error(appErr.Err),
			)
			return c.Status(appErr.Code).JSON(fiber.Map{
				"error":   appErr.Message,
				"code":    appErr.Code,
				"details": appErr.Err.Error(),
			})
		}

		// Handle other errors
		logger.Error("Unhandled error",
			zap.String("path", c.Path()),
			zap.String("method", c.Method()),
			zap.Error(err),
		)

		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Internal server error",
		})
	}
}

// JWT middleware using the auth service
func JWT(authService *auth.Service) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return errors.Unauthorized("Missing authorization header", nil)
		}

		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenStr == authHeader {
			return errors.Unauthorized("Invalid authorization header format", nil)
		}

		// Get user from token
		user, err := authService.GetUserFromToken(tokenStr)
		if err != nil {
			return errors.Unauthorized("Invalid or expired token", err)
		}

		// Add user to context
		c.Locals("user", user)
		return c.Next()
	}
}

// RoleMiddleware checks if user has required role
func RoleMiddleware(requiredRoles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		user := c.Locals("user")
		if user == nil {
			return errors.Unauthorized("User not found in context", nil)
		}

		// Type assertion
		userObj, ok := user.(*models.User)
		if !ok {
			return errors.Unauthorized("Invalid user object in context", nil)
		}

		// Check if user has any of the required roles
		hasRole := false
		for _, role := range requiredRoles {
			if userObj.HasRole(models.UserRole(role)) {
				hasRole = true
				break
			}
		}

		if !hasRole {
			return errors.Forbidden("Insufficient permissions", nil)
		}

		return c.Next()
	}
}

// AdminMiddleware checks if user is admin
func AdminMiddleware() fiber.Handler {
	return RoleMiddleware("admin")
}

// RequestLogger logs incoming requests
func RequestLogger(logger *zap.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := c.Context().Time()

		err := c.Next()

		duration := c.Context().Time().Sub(start)

		logger.Info("Request processed",
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.Int("status", c.Response().StatusCode()),
			zap.Duration("duration", duration),
			zap.String("ip", c.IP()),
			zap.String("user_agent", c.Get("User-Agent")),
		)

		return err
	}
}

// CORS middleware
func CORS() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set("Access-Control-Allow-Origin", "*")
		c.Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS")
		c.Set("Access-Control-Allow-Headers", "Content-Type,Authorization")

		if c.Method() == "OPTIONS" {
			return c.SendStatus(fiber.StatusNoContent)
		}

		return c.Next()
	}
}
